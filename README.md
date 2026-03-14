## Jobs API - Go + MongoDB

Production-ready REST API for a job marketplace built with Go and MongoDB.

## Quick Start

**Prerequisites**
- Go 1.24+
- MongoDB (Atlas or local)

**Installation**

> 1. `git clone https://github.com/t5810m/go-mongodb-api.git`
> 2. `cd go-mongodb-api`

**Setup**

> 1. `cp .env.example .env`
> 2. Edit `.env` with your MongoDB URI
> 3. `go run cmd/seed/seed.go` — Create and seed database
> 4. `go run cmd/main.go` — Start server
>
> Server runs on `http://localhost:8080`

**Docker**

> 1. `cp .env.example .env`
> 2. Edit `.env` with your MongoDB URI
> 3. `docker-compose up --build`
>
> Server runs on `http://localhost:8080`

## Features

- 60+ REST endpoints for job marketplace
- Pagination, filtering, sorting on all list endpoints
- Layered architecture (Handler → Service → Repository → DB)
- 14 entities: Users, Jobs, Applications, Skills, JobCategories, JobSkills, CandidateSkills, Articles, Countries, EducationLevels, JobTypes, KnowledgeLevels, LocationAvailabilities, Resumes
- JWT authentication with role-based access (admin / recruiter / candidate)
- Input validation with field-level error messages
- MongoDB integration with indexes and clean patterns
- ~97% test coverage on middleware, ~98% on services

## Documentation

- [API Endpoints](./docs/API_ENDPOINTS.md)
- [Database Schema](./docs/DATABASE_SCHEMA.md)
- [Architecture](./docs/ARCHITECTURE.md)
- [Testing with Postman](./docs/POSTMAN_TESTING.md)

## Technology Stack

Go · MongoDB · Chi Router · godotenv · Validator · Docker

## Contact

**Zoran Makrevski** — zoran@makrevski.com

Feel free to reach out for questions, suggestions, or collaboration.

## License

MIT - See [LICENSE](./docs/LICENSE)
