# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-11

### Added

- Initial public release of Jobs API
- 9 complete entities (Job, Candidate, Recruiter, Company, Skill, Application, JobCategory, JobSkill, CandidateSkill)
- 35+ REST API endpoints with full CRUD operations
- Pagination support with configurable page size and limits
- Search and filtering capabilities for all entities with case-insensitive matching
- Sorting support with multiple sortable fields per entity
- Input validation with structured error responses
- Layered architecture (Handlers -> Services -> Repositories -> MongoDB)
- MongoDB Atlas integration with environment-based configuration
- Database seeding script with 200+ realistic test records
- Comprehensive README with API documentation and examples
- MIT License for open source usage

### Features

- Pagination: `?page=1&limit=10`
- Filtering: Entity-specific field queries with partial matching
- Sorting: `?sort=field&order=asc/desc`
- Validation: Field-level error messages with HTTP 400 responses

### Technology Stack

- Go 1.24+
- MongoDB with Atlas cloud support
- Chi v5 HTTP router
- go-playground/validator for input validation
- jaswdr/faker for test data generation
- godotenv for environment configuration

### Documentation

- Complete API endpoint documentation
- Architecture and design patterns explained
- Database schema documentation
- Setup and deployment instructions
- Examples for testing with curl

---

## Future Enhancements

- Swagger/OpenAPI documentation
- Authentication and authorization middleware
- Unit and integration tests
- More complex skill-based matching algorithms
- Application workflow management
- API rate limiting
- Comprehensive logging and monitoring
