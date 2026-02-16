package main

import (
	"context"
	"errors"
	"fmt"
	"go-mongodb-api/config"
	"go-mongodb-api/models"
	"log"
	"math/rand"
	"time"

	"github.com/jaswdr/faker/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var fake faker.Faker

func init() {
	fake = faker.New()
}

func generatePhoneNumber() string {
	return fmt.Sprintf("(%03d) %03d-%04d", rand.Intn(900)+100, rand.Intn(900)+100, rand.Intn(9000)+1000)
}

func main() {
	// Initialize configuration
	_, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize MongoDB
	_, err = config.InitMongo()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer func() {
		if err := config.DisconnectMongo(); err != nil {
			log.Printf("error disconnecting MongoDB: %v", err)
		}
	}()

	db, err := config.GetDatabase("jobs_db")
	if err != nil {
		log.Fatalf("Failed to get database: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Clear existing data
	clearCollections(ctx, db)

	// Create data
	fmt.Println("Seeding database...")

	skillIDs := createSkills(ctx, db)
	jobCategoryIDs := createJobCategories(ctx, db)
	createUsers(ctx, db)
	companyIDs := createCompanies(ctx, db)
	recruiterIDs := createRecruiters(ctx, db, companyIDs)
	candidateIDs := createCandidates(ctx, db)
	createCandidateSkills(ctx, db, candidateIDs, skillIDs)
	jobIDs := createJobs(ctx, db, recruiterIDs, companyIDs, jobCategoryIDs)
	createJobSkills(ctx, db, jobIDs, skillIDs)
	createApplications(ctx, db, jobIDs, candidateIDs)
	createResumes(ctx, db, candidateIDs)

	fmt.Println("✓ Database seeding completed successfully!")
}

func clearCollections(ctx context.Context, db *mongo.Database) {
	collections := []string{"users", "companies", "recruiters", "candidates", "jobs", "applications", "jobcategories", "skills", "candidateskills", "jobskills", "resumes"}
	for _, col := range collections {
		err := db.Collection(col).Drop(ctx)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Printf("Warning: could not drop collection %s: %v\n", col, err)
		}
	}
}

func createUsers(ctx context.Context, db *mongo.Database) []bson.ObjectID {
	fmt.Println("Creating users...")
	col := db.Collection("users")

	adminNames := []string{"Alice Johnson", "Bob Smith", "Charlie Brown"}
	var userIDs []bson.ObjectID

	for i := 0; i < len(adminNames); i++ {
		user := models.User{
			ID:                bson.NewObjectID(),
			FirstName:         fake.Person().FirstName(),
			LastName:          fake.Person().LastName(),
			Email:             fmt.Sprintf("admin%d@jobsapi.com", i+1),
			Password:          "hashed_password_123",
			Verified:          true,
			Active:            true,
			TermsAccepted:     true,
			LastTermsAccepted: time.Now().AddDate(0, -1, 0),
			LastLoginTime:     time.Now().AddDate(0, 0, -rand.Intn(30)),
			CreatedTime:       time.Now().AddDate(-1, 0, 0),
			UpdatedTime:       time.Now(),
			CreatedBy:         "system",
			UpdatedBy:         "system",
		}
		result, err := col.InsertOne(ctx, user)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			continue
		}
		userIDs = append(userIDs, result.InsertedID.(bson.ObjectID))
	}

	fmt.Printf("Created %d users\n", len(userIDs))
	return userIDs
}

func createCompanies(ctx context.Context, db *mongo.Database) []bson.ObjectID {
	fmt.Println("Creating companies...")
	col := db.Collection("companies")
	var companyIDs []bson.ObjectID

	companyNames := []string{
		"TechCorp", "CloudSoft", "DataSystems", "WebInnovate", "DevOps Plus",
		"AILabs", "SecureNet", "MobileFirst", "APIGateway", "CloudStack",
		"BigDataCo", "FinTechHub", "HealthTech", "EcommerceIO", "GameEngine",
		"MetaVerse", "BlockChainCo", "RoboTech", "SpaceX", "QuantumCompute",
		"OpenSourceLabs", "StartupHub", "EnterpriseSoft", "SaaScend", "DevCloud",
	}

	for i := 0; i < rand.Intn(11)+15; i++ {
		company := models.Company{
			ID:          bson.NewObjectID(),
			Name:        companyNames[i%len(companyNames)] + fmt.Sprintf(" %d", i/len(companyNames)+1),
			Description: "Building innovative tech solutions for the future",
			Website:     fmt.Sprintf("https://company%d.com", i+1),
			Email:       fmt.Sprintf("hr@company%d.com", i+1),
			Phone:       generatePhoneNumber(),
			Address:     fake.Address().Address(),
			City:        fake.Address().City(),
			PostalCode:  fmt.Sprintf("%05d", rand.Intn(99999)),
			Country:     "US",
			LogoUrl:     fmt.Sprintf("https://logo.company%d.com/logo.png", i+1),
			Verified:    true,
			Active:      true,
			CreatedTime: time.Now().AddDate(-1, rand.Intn(12), -rand.Intn(30)),
			UpdatedTime: time.Now(),
			CreatedBy:   "system",
			UpdatedBy:   "system",
		}
		result, err := col.InsertOne(ctx, company)
		if err != nil {
			log.Printf("Error inserting company: %v", err)
			continue
		}
		companyIDs = append(companyIDs, result.InsertedID.(bson.ObjectID))
	}

	fmt.Printf("Created %d companies\n", len(companyIDs))
	return companyIDs
}

func createRecruiters(ctx context.Context, db *mongo.Database, companyIDs []bson.ObjectID) []bson.ObjectID {
	fmt.Println("Creating recruiters...")
	col := db.Collection("recruiters")
	var recruiterIDs []bson.ObjectID

	recruiterCount := rand.Intn(11) + 30
	for i := 0; i < recruiterCount; i++ {
		isIndependent := rand.Float32() > 0.7 // 30% independent
		var companyID *bson.ObjectID

		if !isIndependent && len(companyIDs) > 0 {
			cid := companyIDs[rand.Intn(len(companyIDs))]
			companyID = &cid
		}

		recruiter := models.Recruiter{
			ID:          bson.NewObjectID(),
			FirstName:   fake.Person().FirstName(),
			LastName:    fake.Person().LastName(),
			Email:       fake.Internet().Email(),
			Password:    "hashed_password_123",
			Phone:       generatePhoneNumber(),
			CompanyID:   companyID,
			Verified:    true,
			Active:      true,
			CreatedTime: time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30)),
			UpdatedTime: time.Now(),
			CreatedBy:   "system",
			UpdatedBy:   "system",
		}
		result, err := col.InsertOne(ctx, recruiter)
		if err != nil {
			log.Printf("Error inserting recruiter: %v", err)
			continue
		}
		recruiterIDs = append(recruiterIDs, result.InsertedID.(bson.ObjectID))
	}

	fmt.Printf("Created %d recruiters\n", len(recruiterIDs))
	return recruiterIDs
}

func createCandidates(ctx context.Context, db *mongo.Database) []bson.ObjectID {
	fmt.Println("Creating candidates...")
	col := db.Collection("candidates")
	var candidateIDs []bson.ObjectID

	candidateCount := rand.Intn(21) + 40
	for i := 0; i < candidateCount; i++ {
		candidate := models.Candidate{
			ID:          bson.NewObjectID(),
			FirstName:   fake.Person().FirstName(),
			LastName:    fake.Person().LastName(),
			Email:       fake.Internet().Email(),
			Password:    "hashed_password_123",
			Verified:    rand.Float32() > 0.2,
			Phone:       generatePhoneNumber(),
			Active:      true,
			CreatedTime: time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30)),
			UpdatedTime: time.Now(),
			CreatedBy:   "system",
			UpdatedBy:   "system",
		}
		result, err := col.InsertOne(ctx, candidate)
		if err != nil {
			log.Printf("Error inserting candidate: %v", err)
			continue
		}
		candidateIDs = append(candidateIDs, result.InsertedID.(bson.ObjectID))
	}

	fmt.Printf("Created %d candidates\n", len(candidateIDs))
	return candidateIDs
}

func createJobs(ctx context.Context, db *mongo.Database, recruiterIDs, companyIDs []bson.ObjectID, categoryIDs []bson.ObjectID) []bson.ObjectID {
	fmt.Println("Creating jobs...")
	col := db.Collection("jobs")
	var jobIDs []bson.ObjectID

	jobTitles := []string{
		"Software Engineer", "Senior Backend Developer", "Frontend Developer", "Full Stack Developer",
		"QA Engineer", "DevOps Engineer", "Mobile Developer (iOS/Android)", "React Developer",
		"Go Developer", "Python Developer", "Java Developer", "Node.js Developer",
		"Database Administrator", "Solutions Architect", "Tech Lead", "Engineering Manager",
		"Data Scientist", "Machine Learning Engineer", "Security Engineer", "Cloud Architect",
		"SEO Specialist", "Web Developer", "UI/UX Designer", "Product Manager",
	}

	jobCount := rand.Intn(16) + 20
	for i := 0; i < jobCount; i++ {
		recruiterID := recruiterIDs[rand.Intn(len(recruiterIDs))]
		companyID := companyIDs[rand.Intn(len(companyIDs))]
		categoryID := categoryIDs[rand.Intn(len(categoryIDs))]

		// Generate valid salary range (min < max)
		salaryMin := 50000 + rand.Intn(50000)
		salaryMax := salaryMin + 50000 + rand.Intn(100000)

		job := models.Job{
			ID:          bson.NewObjectID(),
			Title:       jobTitles[rand.Intn(len(jobTitles))],
			Description: "We are looking for an experienced professional to join our team.",
			RecruiterID: recruiterID,
			CompanyID:   companyID,
			CategoryID:  categoryID,
			Location:    fmt.Sprintf("%s, %s", fake.Address().City(), "US"),
			JobType:     []string{"Full-time", "Part-time", "Contract", "Remote"}[rand.Intn(4)],
			SalaryMin:   salaryMin,
			SalaryMax:   salaryMax,
			Status:      []string{"Active", "Active", "Active", "Closed"}[rand.Intn(4)],
			Active:      rand.Float32() > 0.2,
			CreatedTime: time.Now().AddDate(0, -rand.Intn(6), -rand.Intn(30)),
			UpdatedTime: time.Now(),
			CreatedBy:   "system",
			UpdatedBy:   "system",
		}
		result, err := col.InsertOne(ctx, job)
		if err != nil {
			log.Printf("Error inserting job: %v", err)
			continue
		}
		jobIDs = append(jobIDs, result.InsertedID.(bson.ObjectID))
	}

	fmt.Printf("Created %d jobs\n", len(jobIDs))
	return jobIDs
}

func createSkills(ctx context.Context, db *mongo.Database) []bson.ObjectID {
	fmt.Println("Creating skills...")
	col := db.Collection("skills")
	var skillIDs []bson.ObjectID

	skills := []string{
		// Languages
		"Go", "Python", "JavaScript", "TypeScript", "Java", "C++", "C#", "PHP", "Ruby", "Rust",
		"Kotlin", "Swift", "Objective-C", "Scala", "Elixir", "Clojure", "Haskell", "R", "MATLAB", "Perl",

		// Frameworks & Libraries
		"React", "Vue.js", "Angular", "Next.js", "Svelte", "Ember.js", "Backbone.js", "jQuery",
		"Node.js", "Express.js", "Django", "Flask", "FastAPI", "Spring Boot", "Spring Framework",
		"Laravel", "Symfony", "Ruby on Rails", "ASP.NET Core", "ASP.NET", "Phoenix",

		// Databases
		"MongoDB", "PostgreSQL", "MySQL", "MariaDB", "SQLite", "Oracle", "SQL Server",
		"Redis", "Elasticsearch", "Cassandra", "DynamoDB", "Firebase", "Couchbase", "Neo4j",

		// Cloud & DevOps
		"AWS", "Azure", "Google Cloud", "Docker", "Kubernetes", "Jenkins", "GitLab CI", "GitHub Actions",
		"Terraform", "Ansible", "CloudFormation", "Helm", "Docker Compose",

		// Tools & Platforms
		"Git", "GitHub", "GitLab", "Bitbucket", "Jira", "Confluence", "Slack", "REST API",
		"GraphQL", "gRPC", "SOAP", "Microservices", "Linux", "Windows Server", "macOS",

		// Testing & QA
		"Jest", "Mocha", "Chai", "Pytest", "JUnit", "TestNG", "Selenium", "Cypress",
		"Postman", "SoapUI", "LoadRunner", "JMeter", "TDD", "BDD",

		// Data & Analytics
		"SQL", "Data Analysis", "Statistics", "Machine Learning", "TensorFlow", "PyTorch",
		"Pandas", "NumPy", "Scikit-learn", "Matplotlib", "Tableau", "Power BI",

		// Other
		"Agile", "Scrum", "Kanban", "JIRA", "AWS Lambda", "Serverless", "Web Security",
		"RESTful API Design", "API Development", "System Design", "Architecture", "Leadership",
	}

	for _, skillName := range skills {
		skill := models.Skill{
			ID:          bson.NewObjectID(),
			Name:        skillName,
			Description: fmt.Sprintf("Proficiency in %s", skillName),
			CreatedTime: time.Now().AddDate(-2, 0, 0),
			UpdatedTime: time.Now(),
			CreatedBy:   "system",
			UpdatedBy:   "system",
		}
		result, err := col.InsertOne(ctx, skill)
		if err != nil {
			log.Printf("Error inserting skill: %v", err)
			continue
		}
		skillIDs = append(skillIDs, result.InsertedID.(bson.ObjectID))
	}

	fmt.Printf("Created %d skills\n", len(skillIDs))
	return skillIDs
}

func createJobCategories(ctx context.Context, db *mongo.Database) []bson.ObjectID {
	fmt.Println("Creating job categories...")
	col := db.Collection("jobcategories")
	var categoryIDs []bson.ObjectID

	categories := []string{
		"Backend Development",
		"Frontend Development",
		"Full Stack Development",
		"Mobile Development",
		"DevOps & Infrastructure",
		"Data Science & Analytics",
		"QA & Testing",
		"Security",
	}

	for _, catName := range categories {
		category := models.JobCategory{
			ID:          bson.NewObjectID(),
			Name:        catName,
			Description: fmt.Sprintf("Jobs in %s", catName),
			CreatedTime: time.Now().AddDate(-1, 0, 0),
			UpdatedTime: time.Now(),
			CreatedBy:   "system",
			UpdatedBy:   "system",
		}
		result, err := col.InsertOne(ctx, category)
		if err != nil {
			log.Printf("Error inserting job category: %v", err)
			continue
		}
		categoryIDs = append(categoryIDs, result.InsertedID.(bson.ObjectID))
	}

	fmt.Printf("Created %d job categories\n", len(categoryIDs))
	return categoryIDs
}

func createCandidateSkills(ctx context.Context, db *mongo.Database, candidateIDs, skillIDs []bson.ObjectID) {
	fmt.Println("Creating candidate skills...")
	col := db.Collection("candidateskills")

	proficiencyLevels := []string{"Beginner", "Intermediate", "Advanced", "Expert"}
	skillCount := 0

	for _, candidateID := range candidateIDs {
		// Each candidate has 3-8 skills
		numSkills := rand.Intn(6) + 3
		selectedSkills := make(map[string]bool)
		var skills []bson.ObjectID

		// Keep selecting until we have enough unique skills
		for len(skills) < numSkills {
			skillID := skillIDs[rand.Intn(len(skillIDs))]
			skillKey := skillID.Hex()
			if !selectedSkills[skillKey] {
				selectedSkills[skillKey] = true
				skills = append(skills, skillID)
			}
		}

		for _, skillID := range skills {
			candidateSkill := models.CandidateSkill{
				ID:               bson.NewObjectID(),
				CandidateID:      candidateID,
				SkillID:          skillID,
				ProficiencyLevel: proficiencyLevels[rand.Intn(len(proficiencyLevels))],
				CreatedTime:      time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30)),
				UpdatedTime:      time.Now(),
				CreatedBy:        "system",
				UpdatedBy:        "system",
			}
			_, err := col.InsertOne(ctx, candidateSkill)
			if err != nil {
				log.Printf("Error inserting candidate skill: %v", err)
			}
			skillCount++
		}
	}

	fmt.Printf("Created %d candidate skills\n", skillCount)
}

func createJobSkills(ctx context.Context, db *mongo.Database, jobIDs, skillIDs []bson.ObjectID) {
	fmt.Println("Creating job skills...")
	col := db.Collection("jobskills")

	proficiencyLevels := []string{"Beginner", "Intermediate", "Advanced", "Expert"}
	skillCount := 0

	for _, jobID := range jobIDs {
		// Each job requires 2-5 skills
		numSkills := rand.Intn(4) + 2
		selectedSkills := make(map[string]bool)
		var skills []bson.ObjectID

		// Keep selecting until we have enough unique skills
		for len(skills) < numSkills {
			skillID := skillIDs[rand.Intn(len(skillIDs))]
			skillKey := skillID.Hex()
			if !selectedSkills[skillKey] {
				selectedSkills[skillKey] = true
				skills = append(skills, skillID)
			}
		}

		for _, skillID := range skills {
			jobSkill := models.JobSkill{
				ID:                       bson.NewObjectID(),
				JobID:                    jobID,
				SkillID:                  skillID,
				ProficiencyLevelRequired: proficiencyLevels[rand.Intn(len(proficiencyLevels))],
				IsRequired:               rand.Float32() > 0.3, // 70% required, 30% nice-to-have
				CreatedTime:              time.Now().AddDate(0, -rand.Intn(6), -rand.Intn(30)),
				UpdatedTime:              time.Now(),
				CreatedBy:                "system",
				UpdatedBy:                "system",
			}
			_, err := col.InsertOne(ctx, jobSkill)
			if err != nil {
				log.Printf("Error inserting job skill: %v", err)
			}
			skillCount++
		}
	}

	fmt.Printf("Created %d job skills\n", skillCount)
}

func createApplications(ctx context.Context, db *mongo.Database, jobIDs, candidateIDs []bson.ObjectID) {
	fmt.Println("Creating applications...")
	col := db.Collection("applications")

	statuses := []string{"Applied", "Under Review", "Rejected", "Accepted", "Withdrawn"}

	// Create 60-100 applications
	appCount := rand.Intn(41) + 60
	for i := 0; i < appCount; i++ {
		jobID := jobIDs[rand.Intn(len(jobIDs))]
		candidateID := candidateIDs[rand.Intn(len(candidateIDs))]

		application := models.Application{
			ID:            bson.NewObjectID(),
			JobID:         jobID,
			CandidateID:   candidateID,
			Status:        statuses[rand.Intn(len(statuses))],
			RecruiterNote: fake.Lorem().Sentence(5),
			AppliedTime:   time.Now().AddDate(0, -rand.Intn(3), -rand.Intn(30)),
			UpdatedTime:   time.Now(),
			CreatedBy:     "system",
			UpdatedBy:     "system",
		}
		_, err := col.InsertOne(ctx, application)
		if err != nil {
			log.Printf("Error inserting application: %v", err)
		}
	}

	fmt.Printf("Created %d applications\n", appCount)
}

func createResumes(ctx context.Context, db *mongo.Database, candidateIDs []bson.ObjectID) {
	fmt.Println("Creating resumes...")
	col := db.Collection("resumes")

	for _, candidateID := range candidateIDs {
		// 80% of candidates have a resume
		if rand.Float32() > 0.2 {
			resume := models.Resume{
				ID:           bson.NewObjectID(),
				CandidateID:  candidateID,
				FileUrl:      fmt.Sprintf("https://resumes.example.com/%s.pdf", candidateID.Hex()),
				FileName:     fmt.Sprintf("%s_resume.pdf", candidateID.Hex()[:8]),
				UploadedTime: time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30)),
				UpdatedTime:  time.Now(),
				CreatedBy:    "system",
				UpdatedBy:    "system",
			}
			_, err := col.InsertOne(ctx, resume)
			if err != nil {
				log.Printf("Error inserting resume: %v", err)
			}
		}
	}

	fmt.Printf("Created resumes for candidates\n")
}
