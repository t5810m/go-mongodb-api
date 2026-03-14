package main

import (
	"context"
	"errors"
	"fmt"
	"go-mongodb-api/config"
	"go-mongodb-api/handlers"
	authMW "go-mongodb-api/middleware"
	"go-mongodb-api/repositories"
	"go-mongodb-api/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize configuration
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize MongoDB connection
	_, err = config.InitMongo()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer func() {
		if err := config.DisconnectMongo(); err != nil {
			log.Printf("error disconnecting MongoDB: %v", err)
		}
	}()

	// Get database
	db, err := config.GetDatabase("job_board")
	if err != nil {
		log.Fatalf("Failed to get database: %v", err)
	}
	fmt.Printf("Connected to database: %s\n", db.Name())

	// Ensure indexes
	idxCtx, idxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer idxCancel()
	if err := config.EnsureIndexes(idxCtx, db); err != nil {
		log.Fatalf("Failed to ensure indexes: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	jobRepo := repositories.NewJobRepository(db)
	skillRepo := repositories.NewSkillRepository(db)
	applicationRepo := repositories.NewApplicationRepository(db)
	jobCategoryRepo := repositories.NewJobCategoryRepository(db)
	candidateSkillRepo := repositories.NewCandidateSkillRepository(db)
	jobSkillRepo := repositories.NewJobSkillRepository(db)
	articleRepo := repositories.NewArticleRepository(db)
	countryRepo := repositories.NewCountryRepository(db)
	educationLevelRepo := repositories.NewEducationLevelRepository(db)
	jobTypeRepo := repositories.NewJobTypeRepository(db)
	knowledgeLevelRepo := repositories.NewKnowledgeLevelRepository(db)
	locationAvailabilityRepo := repositories.NewLocationAvailabilityRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	jobService := services.NewJobService(jobRepo, userRepo, jobCategoryRepo)
	skillService := services.NewSkillService(skillRepo)
	applicationService := services.NewApplicationService(applicationRepo, jobRepo, userRepo)
	jobCategoryService := services.NewJobCategoryService(jobCategoryRepo)
	candidateSkillService := services.NewCandidateSkillService(candidateSkillRepo, userRepo, skillRepo)
	jobSkillService := services.NewJobSkillService(jobSkillRepo, jobRepo, skillRepo)
	authService := services.NewAuthService(userService, cfg.JWTSecret)
	articleService := services.NewArticleService(articleRepo)
	countryService := services.NewCountryService(countryRepo)
	educationLevelService := services.NewEducationLevelService(educationLevelRepo)
	jobTypeService := services.NewJobTypeService(jobTypeRepo)
	knowledgeLevelService := services.NewKnowledgeLevelService(knowledgeLevelRepo)
	locationAvailabilityService := services.NewLocationAvailabilityService(locationAvailabilityRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	jobHandler := handlers.NewJobHandler(jobService)
	skillHandler := handlers.NewSkillHandler(skillService)
	applicationHandler := handlers.NewApplicationHandler(applicationService)
	jobCategoryHandler := handlers.NewJobCategoryHandler(jobCategoryService)
	candidateSkillHandler := handlers.NewCandidateSkillHandler(candidateSkillService)
	jobSkillHandler := handlers.NewJobSkillHandler(jobSkillService)
	authHandler := handlers.NewAuthHandler(authService)
	articleHandler := handlers.NewArticleHandler(articleService)
	countryHandler := handlers.NewCountryHandler(countryService)
	educationLevelHandler := handlers.NewEducationLevelHandler(educationLevelService)
	jobTypeHandler := handlers.NewJobTypeHandler(jobTypeService)
	knowledgeLevelHandler := handlers.NewKnowledgeLevelHandler(knowledgeLevelService)
	locationAvailabilityHandler := handlers.NewLocationAvailabilityHandler(locationAvailabilityService)

	// Validate port
	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		log.Fatalf("Invalid port: %s (must be 1-65535)", cfg.Port)
	}

	// Setup router with middleware
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	// ── Public routes ────────────────────────────────────────────────────────
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, `{"status":"healthy","version":"1.0"}`); err != nil {
			log.Printf("error writing health response: %v", err)
		}
	})

	// Auth endpoints
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/register", authHandler.Register)

	// Public read-only
	r.Get("/jobs", jobHandler.GetAllJobs)
	r.Get("/jobs/{id}", jobHandler.GetJobByID)
	r.Get("/users/{userId}/jobs", jobHandler.GetJobsByUser)
	r.Get("/skills", skillHandler.GetAllSkills)
	r.Get("/skills/{id}", skillHandler.GetSkillByID)
	r.Get("/jobcategories", jobCategoryHandler.GetAllJobCategories)
	r.Get("/jobcategories/{id}", jobCategoryHandler.GetJobCategoryByID)
	r.Get("/jobs/{jobId}/skills", jobSkillHandler.GetJobSkillsByJobID)
	r.Get("/articles", articleHandler.GetAllArticles)
	r.Get("/articles/{id}", articleHandler.GetArticleByID)
	r.Get("/countries", countryHandler.GetAllCountries)
	r.Get("/countries/{id}", countryHandler.GetCountryByID)
	r.Get("/educationlevels", educationLevelHandler.GetAllEducationLevels)
	r.Get("/educationlevels/{id}", educationLevelHandler.GetEducationLevelByID)
	r.Get("/jobtypes", jobTypeHandler.GetAllJobTypes)
	r.Get("/jobtypes/{id}", jobTypeHandler.GetJobTypeByID)
	r.Get("/knowledgelevels", knowledgeLevelHandler.GetAllKnowledgeLevels)
	r.Get("/knowledgelevels/{id}", knowledgeLevelHandler.GetKnowledgeLevelByID)
	r.Get("/locationavailabilities", locationAvailabilityHandler.GetAllLocationAvailabilities)
	r.Get("/locationavailabilities/{id}", locationAvailabilityHandler.GetLocationAvailabilityByID)

	// Public user listing (supports ?role=candidate etc.)
	r.Get("/users", userHandler.GetAllUsers)
	r.Get("/users/{id}", userHandler.GetUserByID)

	// ── Authenticated routes ─────────────────────────────────────────────────
	r.Group(func(r chi.Router) {
		r.Use(authMW.Authenticate(cfg.JWTSecret))

		// admin only
		r.Group(func(r chi.Router) {
			r.Use(authMW.RequireRoles("admin"))
			r.Post("/users", userHandler.CreateUser)
			r.Put("/users/{id}", userHandler.UpdateUser)
			r.Delete("/users/{id}", userHandler.DeleteUser)
			r.Post("/skills", skillHandler.CreateSkill)
			r.Put("/skills/{id}", skillHandler.UpdateSkill)
			r.Delete("/skills/{id}", skillHandler.DeleteSkill)
			r.Post("/jobcategories", jobCategoryHandler.CreateJobCategory)
			r.Put("/jobcategories/{id}", jobCategoryHandler.UpdateJobCategory)
			r.Delete("/jobcategories/{id}", jobCategoryHandler.DeleteJobCategory)
			r.Post("/articles", articleHandler.CreateArticle)
			r.Put("/articles/{id}", articleHandler.UpdateArticle)
			r.Delete("/articles/{id}", articleHandler.DeleteArticle)
			r.Post("/countries", countryHandler.CreateCountry)
			r.Put("/countries/{id}", countryHandler.UpdateCountry)
			r.Delete("/countries/{id}", countryHandler.DeleteCountry)
			r.Post("/educationlevels", educationLevelHandler.CreateEducationLevel)
			r.Put("/educationlevels/{id}", educationLevelHandler.UpdateEducationLevel)
			r.Delete("/educationlevels/{id}", educationLevelHandler.DeleteEducationLevel)
			r.Post("/jobtypes", jobTypeHandler.CreateJobType)
			r.Put("/jobtypes/{id}", jobTypeHandler.UpdateJobType)
			r.Delete("/jobtypes/{id}", jobTypeHandler.DeleteJobType)
			r.Post("/knowledgelevels", knowledgeLevelHandler.CreateKnowledgeLevel)
			r.Put("/knowledgelevels/{id}", knowledgeLevelHandler.UpdateKnowledgeLevel)
			r.Delete("/knowledgelevels/{id}", knowledgeLevelHandler.DeleteKnowledgeLevel)
			r.Post("/locationavailabilities", locationAvailabilityHandler.CreateLocationAvailability)
			r.Put("/locationavailabilities/{id}", locationAvailabilityHandler.UpdateLocationAvailability)
			r.Delete("/locationavailabilities/{id}", locationAvailabilityHandler.DeleteLocationAvailability)
			r.Get("/candidateskills", candidateSkillHandler.GetAllCandidateSkills)
			r.Get("/jobskills", jobSkillHandler.GetAllJobSkills)
			r.Get("/applications", applicationHandler.GetAllApplications)
		})

		// admin + recruiter
		r.Group(func(r chi.Router) {
			r.Use(authMW.RequireRoles("admin", "recruiter"))
			r.Post("/jobs", jobHandler.CreateJob)
			r.Delete("/jobs/{id}", jobHandler.DeleteJob)
			r.Post("/jobskills", jobSkillHandler.CreateJobSkill)
			r.Get("/jobskills/{id}", jobSkillHandler.GetJobSkillByID)
			r.Put("/jobskills/{id}", jobSkillHandler.UpdateJobSkillProficiencyLevel)
			r.Delete("/jobskills/{id}", jobSkillHandler.DeleteJobSkill)
			r.Get("/jobs/{jobId}/applications", applicationHandler.GetApplicationsByJobID)
			r.Put("/applications/{id}", applicationHandler.UpdateApplicationStatus)
		})

		// admin + candidate
		r.Group(func(r chi.Router) {
			r.Use(authMW.RequireRoles("admin", "candidate"))
			r.Post("/applications", applicationHandler.CreateApplication)
			r.Post("/candidateskills", candidateSkillHandler.CreateCandidateSkill)
			r.Put("/candidateskills/{id}", candidateSkillHandler.UpdateCandidateSkillProficiencyLevel)
			r.Delete("/candidateskills/{id}", candidateSkillHandler.DeleteCandidateSkill)
			r.Get("/users/{userId}/applications", applicationHandler.GetApplicationsByUserID)
			r.Delete("/applications/{id}", applicationHandler.DeleteApplication)
			r.Get("/candidateskills/{id}", candidateSkillHandler.GetCandidateSkillByID)
		})

		// admin + candidate + recruiter
		r.Group(func(r chi.Router) {
			r.Use(authMW.RequireRoles("admin", "candidate", "recruiter"))
			r.Get("/applications/{id}", applicationHandler.GetApplicationByID)
			r.Get("/users/{userId}/skills", candidateSkillHandler.GetCandidateSkillsByUserID)
		})
	})

	// Start server
	fmt.Printf("Starting API server on port %s\n", cfg.Port)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	serverErrors := make(chan error, 1)

	go func() {
		fmt.Println("Server is running. Press Ctrl+C to gracefully shutdown...")
		serverErrors <- server.ListenAndServe()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		fmt.Println("\nShutdown signal received, gracefully shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during graceful shutdown: %v", err)
		} else {
			fmt.Println("Server shut down gracefully")
		}
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}
}
