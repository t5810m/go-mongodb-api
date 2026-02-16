## Jobs API - Go + MongoDB

Production-ready REST API for a job marketplace built with Go and MongoDB.

## Quick Start

**Prerequisites**
- Go 1.24+
- MongoDB (Atlas or local)

**Installation**

> 1. `git clone https://github.com/yourusername/go-mongodb-api.git`
> 2. `cd go-mongodb-api`

**Setup**

> 1. `cp .env.example .env`
> 2. Edit `.env` with your MongoDB URI
> 3. `go run cmd/seed/seed.go` — Create and seed database
> 4. `go run cmd/main.go` — Start server
>
> Server runs on `http://localhost:8080`

## Features

- 35+ REST endpoints for job marketplace
- Pagination, filtering, sorting
- Layered architecture (Handler → Service → Repository → DB)
- 11 entities: Users, Companies, Recruiters, Candidates, Jobs, Skills, JobCategories, Applications, JobSkills, CandidateSkills, Resumes
- Input validation with field-level error messages
- MongoDB integration with clean patterns

## Documentation

- [API Endpoints](./docs/API_ENDPOINTS.md)
- [Database Schema](./docs/DATABASE_SCHEMA.md)
- [Architecture](./docs/ARCHITECTURE.md)
- [Testing with Postman](./docs/POSTMAN_TESTING.md)

## Technology Stack

Go · MongoDB · Chi Router · godotenv · Validator

## License

MIT - See [LICENSE](./LICENSE)

---

**Zoran Makrevski** | zoran@makrevski.com
