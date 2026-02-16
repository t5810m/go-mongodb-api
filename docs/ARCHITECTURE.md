# Architecture

## Layered Architecture Pattern

```
┌─────────────────────────────────────────┐
│  HTTP Handlers (Chi Router)             │
│  - Parse requests                       │
│  - Validate input                       │
│  - Return responses                     │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│  Service Layer                          │
│  - Business logic                       │
│  - Cross-entity validation              │
│  - Domain rules enforcement             │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│  Repository Layer                       │
│  - Database operations (CRUD)           │
│  - Query builders                       │
│  - Data transformation                  │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│  MongoDB                                │
│  - Data persistence                     │
│  - Indexes & aggregations               │
└─────────────────────────────────────────┘
```

## Benefits

- **Separation of Concerns** - Each layer has a single responsibility
- **Testability** - Easy to test each layer independently with mocks
- **Reusability** - Services can be used by different handlers
- **Maintainability** - Changes are isolated to specific layers
- **Database Independence** - Swapping MongoDB for another DB requires only repository changes

## Project Structure

```
api-m/
├── cmd/
│   ├── main.go                    # API server entry point, routing setup
│   └── seed/
│       └── seed.go                # Database seeding with test data
├── config/
│   ├── config.go                  # Configuration loading from .env
│   └── mongo.go                   # MongoDB connection management
├── models/
│   ├── user.go
│   ├── company.go
│   ├── recruiter.go
│   ├── candidate.go
│   ├── job.go
│   ├── application.go
│   ├── skill.go
│   ├── jobcategory.go
│   ├── jobskill.go
│   ├── candidateskill.go
│   ├── resume.go
│   └── skills.go                  # Skills list constant
├── handlers/
│   ├── user.go
│   ├── company.go
│   ├── recruiter.go
│   ├── candidate.go
│   ├── job.go
│   ├── application.go
│   ├── skill.go
│   ├── jobcategory.go
│   ├── jobskill.go
│   ├── candidateskill.go
│   └── auth.go                    # Authentication placeholder
├── services/
│   ├── user.go
│   ├── company.go
│   ├── recruiter.go
│   ├── candidate.go
│   ├── job.go
│   ├── application.go
│   ├── skill.go
│   ├── jobcategory.go
│   ├── jobskill.go
│   └── candidateskill.go
├── repositories/
│   ├── user.go
│   ├── company.go
│   ├── recruiter.go
│   ├── candidate.go
│   ├── job.go
│   ├── application.go
│   ├── skill.go
│   ├── jobcategory.go
│   ├── jobskill.go
│   └── candidateskill.go
├── helpers/
│   ├── pagination.go              # Pagination utilities
│   └── validator.go               # Request validation
├── docs/
│   ├── API_ENDPOINTS.md           # Complete endpoint reference
│   ├── DATABASE_SCHEMA.md         # MongoDB schema documentation
│   ├── ARCHITECTURE.md            # This file
│   ├── POSTMAN_TESTING.md         # Postman collection guide
│   └── API-M.postman_collection.json  # Postman collection
├── .env.example                   # Environment variables template
├── .gitignore                     # Git ignore rules
├── go.mod                         # Go module definition
├── go.sum                         # Go module checksums
├── CHANGELOG.md                   # Version history
└── README.md                      # Project overview
```

## Data Flow Example: Creating a Job

```
1. HTTP Request
   POST /jobs
   {
     "title": "Senior Backend Developer",
     "description": "...",
     "recruiter_id": "ObjectID",
     "company_id": "ObjectID",
     "category_id": "ObjectID",
     "location": "London",
     "job_type": "full-time",
     "salary_min": 60000,
     "salary_max": 80000,
     "status": "active"
   }

2. Handler (handlers/job.go)
   - Decode JSON body into Job struct
   - Validate required fields and constraints
   - Set timestamps and metadata
   - Call service.CreateJob()

3. Service (services/job.go)
   - Validate recruiter exists
   - Validate company exists
   - Validate category exists
   - Apply business rules
   - Call repository.Create()

4. Repository (repositories/job.go)
   - Convert IDs from string to ObjectID
   - Insert document into MongoDB
   - Capture generated _id
   - Return job with ID populated

5. HTTP Response
   201 Created
   {
     "id": "newly-generated-ObjectID",
     "title": "Senior Backend Developer",
     ...
     "created_time": "2026-02-15T21:30:00Z",
     "updated_time": "2026-02-15T21:30:00Z"
   }
```

## Key Technologies

- **Go 1.24+** - Fast, compiled language with excellent concurrency
- **Chi Router** - Lightweight, fast HTTP router
- **MongoDB Driver v2** - Official MongoDB driver for Go
- **Validator** - Input validation with struct tags
- **godotenv** - Environment variable loading
- **faker** - Realistic test data generation

## Design Patterns Used

- **Repository Pattern** - Abstract database operations
- **Service Layer Pattern** - Encapsulate business logic
- **Dependency Injection** - Services receive dependencies via constructor
- **Error Handling** - Consistent error propagation up the stack
- **Pagination** - Efficient data retrieval for large datasets
