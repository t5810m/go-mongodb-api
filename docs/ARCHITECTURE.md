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
│  MongoDB Atlas                          │
│  - Data persistence                     │
│  - Indexes                              │
└─────────────────────────────────────────┘
```

## Benefits

- **Separation of Concerns** — each layer has a single responsibility
- **Testability** — each layer can be tested independently using mocks
- **Reusability** — services can be used by different handlers
- **Maintainability** — changes are isolated to specific layers
- **Database Independence** — swapping MongoDB requires only repository changes

---

## Project Structure

```
go-mongodb-api/
├── cmd/
│   ├── main.go                        # Entry point, routing setup
│   ├── seed/
│   │   └── seed.go                    # Database seeder with realistic test data
│   └── migrate/                       # Database migration utilities
├── config/
│   ├── config.go                      # Configuration loading from .env
│   ├── mongo.go                       # MongoDB connection management
│   └── indexes.go                     # MongoDB index definitions
├── models/
│   ├── user.go                        # User (admin / candidate / recruiter)
│   ├── job.go
│   ├── application.go
│   ├── candidateskill.go
│   ├── jobskill.go
│   ├── skills.go
│   ├── jobcategory.go
│   ├── resume.go
│   ├── article.go
│   ├── country.go
│   ├── educationlevel.go
│   ├── jobtype.go
│   ├── knowledgelevel.go
│   └── locationavailability.go
├── handlers/
│   ├── auth.go                        # Login + Register
│   ├── user.go
│   ├── job.go
│   ├── application.go
│   ├── candidateskill.go
│   ├── jobskill.go
│   ├── skill.go
│   ├── jobcategory.go
│   ├── article.go
│   ├── country.go
│   ├── educationlevel.go
│   ├── jobtype.go
│   ├── knowledgelevel.go
│   └── locationavailability.go
├── services/
│   ├── auth.go
│   ├── user.go
│   ├── job.go
│   ├── application.go
│   ├── candidateskill.go
│   ├── jobskill.go
│   ├── skill.go
│   ├── jobcategory.go
│   ├── article.go
│   ├── country.go
│   ├── educationlevel.go
│   ├── jobtype.go
│   ├── knowledgelevel.go
│   └── locationavailability.go
├── repositories/
│   ├── user.go
│   ├── job.go
│   ├── application.go
│   ├── candidateskill.go
│   ├── jobskill.go
│   ├── skill.go
│   ├── jobcategory.go
│   ├── article.go
│   ├── country.go
│   ├── educationlevel.go
│   ├── jobtype.go
│   ├── knowledgelevel.go
│   └── locationavailability.go
├── interfaces/
│   ├── repository.go                  # Repository interfaces
│   └── service.go                     # Service interfaces
├── mocks/
│   ├── repository_mocks.go            # Testify mock implementations
│   └── service_mocks.go
├── middleware/
│   └── auth.go                        # JWT authentication + role enforcement
├── helpers/
│   ├── pagination.go                  # Pagination utilities
│   └── validator.go                   # Request validation
├── docs/
│   ├── API_ENDPOINTS.md
│   ├── DATABASE_SCHEMA.md
│   ├── ARCHITECTURE.md
│   ├── POSTMAN_TESTING.md
│   └── GO-MONGODB-API.postman_collection.json
├── .env                               # Environment variables (not committed)
├── .env.example                       # Environment variables template
├── go.mod
└── go.sum
```

---

## Authentication & Authorization

JWT-based authentication with role enforcement middleware.

```
POST /auth/login    → returns JWT token (valid 24h)
POST /auth/register → creates user, returns nothing (user must login)
```

### Roles

| Role | Permissions |
|------|-------------|
| `admin` | Full access to all endpoints |
| `recruiter` | Create/delete jobs, manage job skills, view/update applications |
| `candidate` | Submit/delete applications, manage own skills |

### Route Groups

```
Public          → no token required
Admin only      → requires role=admin
Admin+Recruiter → requires role=admin or recruiter
Admin+Candidate → requires role=admin or candidate
All roles       → any authenticated user
```

---

## Data Flow Example: Submitting a Job Application

```
1. POST /applications
   Authorization: Bearer <candidate_jwt>
   {
     "job_id": "ObjectID",
     "user_id": "ObjectID",
     "status": "applied"
   }

2. Middleware
   - Validates JWT signature
   - Extracts role from claims
   - Enforces role=candidate or admin

3. Handler (handlers/application.go)
   - Decodes JSON body
   - Validates required fields
   - Calls service.CreateApplication()

4. Service (services/application.go)
   - Validates job exists
   - Validates user exists and is a candidate
   - Calls repository.Create()

5. Repository (repositories/application.go)
   - Converts string IDs to ObjectIDs
   - Inserts document into MongoDB
   - Returns populated application

6. Response: 201 Created
   {
     "id": "ObjectID",
     "job_id": "ObjectID",
     "user_id": "ObjectID",
     "status": "applied",
     "applied_time": "2026-03-15T20:00:00Z",
     ...
   }
```

---

## Key Technologies

| Technology | Purpose |
|-----------|---------|
| Go 1.24 | Language |
| Chi Router | HTTP routing |
| MongoDB Driver v2 | Database driver |
| JWT (golang-jwt/jwt v5) | Authentication tokens |
| bcrypt | Password hashing |
| go-playground/validator | Struct validation |
| godotenv | Environment variables |
| testify | Unit testing + mocks |
| faker | Realistic seed data |

---

## Test Coverage

```
handlers    ~78%
middleware  ~97%
services    ~98%
```

Run tests:
```bash
go test ./...
go test ./... -cover
```
