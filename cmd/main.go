package main

import (
	"context"
	"errors"
	"fmt"
	"go-mongodb-api/config"
	"go-mongodb-api/handlers"
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
	db, err := config.GetDatabase("jobs_db")
	if err != nil {
		log.Fatalf("Failed to get database: %v", err)
	}
	fmt.Printf("Connected to database: %s\n", db.Name())

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	jobRepo := repositories.NewJobRepository(db)
	candidateRepo := repositories.NewCandidateRepository(db)
	recruiterRepo := repositories.NewRecruiterRepository(db)
	companyRepo := repositories.NewCompanyRepository(db)
	skillRepo := repositories.NewSkillRepository(db)
	applicationRepo := repositories.NewApplicationRepository(db)
	jobCategoryRepo := repositories.NewJobCategoryRepository(db)
	candidateSkillRepo := repositories.NewCandidateSkillRepository(db)
	jobSkillRepo := repositories.NewJobSkillRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	jobService := services.NewJobServiceWithDeps(jobRepo, recruiterRepo, companyRepo, jobCategoryRepo)
	candidateService := services.NewCandidateService(candidateRepo)
	recruiterService := services.NewRecruiterService(recruiterRepo)
	companyService := services.NewCompanyService(companyRepo)
	skillService := services.NewSkillService(skillRepo)
	applicationService := services.NewApplicationServiceWithDeps(applicationRepo, jobRepo, candidateRepo)
	jobCategoryService := services.NewJobCategoryService(jobCategoryRepo)
	candidateSkillService := services.NewCandidateSkillServiceWithDeps(candidateSkillRepo, candidateRepo, skillRepo)
	jobSkillService := services.NewJobSkillServiceWithDeps(jobSkillRepo, jobRepo, skillRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	jobHandler := handlers.NewJobHandler(jobService)
	candidateHandler := handlers.NewCandidateHandler(candidateService)
	recruiterHandler := handlers.NewRecruiterHandler(recruiterService)
	companyHandler := handlers.NewCompanyHandler(companyService)
	skillHandler := handlers.NewSkillHandler(skillService)
	applicationHandler := handlers.NewApplicationHandler(applicationService)
	jobCategoryHandler := handlers.NewJobCategoryHandler(jobCategoryService)
	candidateSkillHandler := handlers.NewCandidateSkillHandler(candidateSkillService)
	jobSkillHandler := handlers.NewJobSkillHandler(jobSkillService)

	// Validate port
	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		log.Fatalf("Invalid port: %s (must be 1-65535)", cfg.Port)
	}

	// Setup router with middleware
	r := chi.NewRouter()
	r.Use(middleware.Recoverer) // Panic recovery middleware

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, `{"status":"healthy","version":"1.0"}`); err != nil {
			log.Printf("error writing health response: %v", err)
		}
	})

	// User routes
	r.Post("/users", userHandler.CreateUser)
	r.Get("/users", userHandler.GetAllUsers)
	r.Get("/users/{id}", userHandler.GetUserByID)
	r.Delete("/users/{id}", userHandler.DeleteUser)

	// Job routes
	r.Post("/jobs", jobHandler.CreateJob)
	r.Get("/jobs", jobHandler.GetAllJobs)
	r.Get("/jobs/{id}", jobHandler.GetJobByID)
	r.Get("/companies/{companyId}/jobs", jobHandler.GetJobsByCompany)
	r.Delete("/jobs/{id}", jobHandler.DeleteJob)

	// Candidate routes
	r.Post("/candidates", candidateHandler.CreateCandidate)
	r.Get("/candidates", candidateHandler.GetAllCandidates)
	r.Get("/candidates/{id}", candidateHandler.GetCandidateByID)
	r.Delete("/candidates/{id}", candidateHandler.DeleteCandidate)

	// Recruiter routes
	r.Post("/recruiters", recruiterHandler.CreateRecruiter)
	r.Get("/recruiters", recruiterHandler.GetAllRecruiters)
	r.Get("/recruiters/{id}", recruiterHandler.GetRecruiterByID)
	r.Delete("/recruiters/{id}", recruiterHandler.DeleteRecruiter)

	// Company routes
	r.Post("/companies", companyHandler.CreateCompany)
	r.Get("/companies", companyHandler.GetAllCompanies)
	r.Get("/companies/{id}", companyHandler.GetCompanyByID)
	r.Delete("/companies/{id}", companyHandler.DeleteCompany)

	// Skill routes
	r.Post("/skills", skillHandler.CreateSkill)
	r.Get("/skills", skillHandler.GetAllSkills)
	r.Get("/skills/{id}", skillHandler.GetSkillByID)
	r.Delete("/skills/{id}", skillHandler.DeleteSkill)

	// Application routes
	r.Post("/applications", applicationHandler.CreateApplication)
	r.Get("/applications", applicationHandler.GetAllApplications)
	r.Get("/applications/{id}", applicationHandler.GetApplicationByID)
	r.Get("/jobs/{jobId}/applications", applicationHandler.GetApplicationsByJobID)
	r.Get("/candidates/{candidateId}/applications", applicationHandler.GetApplicationsByCandidateID)
	r.Put("/applications/{id}", applicationHandler.UpdateApplicationStatus)
	r.Delete("/applications/{id}", applicationHandler.DeleteApplication)

	// Job Category routes
	r.Post("/jobcategories", jobCategoryHandler.CreateJobCategory)
	r.Get("/jobcategories", jobCategoryHandler.GetAllJobCategories)
	r.Get("/jobcategories/{id}", jobCategoryHandler.GetJobCategoryByID)
	r.Delete("/jobcategories/{id}", jobCategoryHandler.DeleteJobCategory)

	// Candidate Skill routes
	r.Post("/candidateskills", candidateSkillHandler.CreateCandidateSkill)
	r.Get("/candidateskills", candidateSkillHandler.GetAllCandidateSkills)
	r.Get("/candidateskills/{id}", candidateSkillHandler.GetCandidateSkillByID)
	r.Get("/candidates/{candidateId}/skills", candidateSkillHandler.GetCandidateSkillsByCandidateID)
	r.Put("/candidateskills/{id}", candidateSkillHandler.UpdateCandidateSkillProficiencyLevel)
	r.Delete("/candidateskills/{id}", candidateSkillHandler.DeleteCandidateSkill)

	// Job Skill routes
	r.Post("/jobskills", jobSkillHandler.CreateJobSkill)
	r.Get("/jobskills", jobSkillHandler.GetAllJobSkills)
	r.Get("/jobskills/{id}", jobSkillHandler.GetJobSkillByID)
	r.Get("/jobs/{jobId}/skills", jobSkillHandler.GetJobSkillsByJobID)
	r.Put("/jobskills/{id}", jobSkillHandler.UpdateJobSkillProficiencyLevel)
	r.Delete("/jobskills/{id}", jobSkillHandler.DeleteJobSkill)

	// Start server
	fmt.Printf("Starting API server on port %s\n", cfg.Port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET /health")
	fmt.Println("  POST /users")
	fmt.Println("  GET /users")
	fmt.Println("  GET /users/{id}")
	fmt.Println("  DELETE /users/{id}")
	fmt.Println("  POST /jobs")
	fmt.Println("  GET /jobs")
	fmt.Println("  GET /jobs/{id}")
	fmt.Println("  GET /companies/{companyId}/jobs")
	fmt.Println("  DELETE /jobs/{id}")
	fmt.Println("  POST /candidates")
	fmt.Println("  GET /candidates")
	fmt.Println("  GET /candidates/{id}")
	fmt.Println("  DELETE /candidates/{id}")
	fmt.Println("  POST /recruiters")
	fmt.Println("  GET /recruiters")
	fmt.Println("  GET /recruiters/{id}")
	fmt.Println("  DELETE /recruiters/{id}")
	fmt.Println("  POST /companies")
	fmt.Println("  GET /companies")
	fmt.Println("  GET /companies/{id}")
	fmt.Println("  DELETE /companies/{id}")
	fmt.Println("  POST /skills")
	fmt.Println("  GET /skills")
	fmt.Println("  GET /skills/{id}")
	fmt.Println("  DELETE /skills/{id}")
	fmt.Println("  POST /applications")
	fmt.Println("  GET /applications")
	fmt.Println("  GET /applications/{id}")
	fmt.Println("  GET /jobs/{jobId}/applications")
	fmt.Println("  GET /candidates/{candidateId}/applications")
	fmt.Println("  PUT /applications/{id}")
	fmt.Println("  DELETE /applications/{id}")
	fmt.Println("  POST /jobcategories")
	fmt.Println("  GET /jobcategories")
	fmt.Println("  GET /jobcategories/{id}")
	fmt.Println("  DELETE /jobcategories/{id}")
	fmt.Println("  POST /candidateskills")
	fmt.Println("  GET /candidateskills")
	fmt.Println("  GET /candidateskills/{id}")
	fmt.Println("  GET /candidates/{candidateId}/skills")
	fmt.Println("  PUT /candidateskills/{id}")
	fmt.Println("  DELETE /candidateskills/{id}")
	fmt.Println("  POST /jobskills")
	fmt.Println("  GET /jobskills")
	fmt.Println("  GET /jobskills/{id}")
	fmt.Println("  GET /jobs/{jobId}/skills")
	fmt.Println("  PUT /jobskills/{id}")
	fmt.Println("  DELETE /jobskills/{id}")

	// Create HTTP server with proper shutdown handling
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Channel for shutdown errors
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		fmt.Println("Server is running. Press Ctrl+C to gracefully shutdown...")
		serverErrors <- server.ListenAndServe()
	}()

	// Wait for interrupt signal (Ctrl+C) or server error
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		fmt.Println("\nShutdown signal received, gracefully shutting down...")
		// Create context with 30-second timeout for graceful shutdown
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
