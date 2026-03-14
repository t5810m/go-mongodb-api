# API Endpoints

Base URL: `http://localhost:8080`

## Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/auth/login` | Public | Login with email + password, returns JWT |
| POST | `/auth/register` | Public | Register a new user (candidate or recruiter) |

### Login request body
```json
{
  "email": "admin1@jobsapi.com",
  "password": "SeedPassword123!"
}
```

### Register request body
```json
{
  "first_name": "Jane",
  "last_name": "Doe",
  "email": "jane.doe@example.com",
  "password": "SecurePass123!",
  "phone": "+1234567890",
  "role": "candidate",
  "terms_accepted": true
}
```
> `role` must be one of: `admin`, `candidate`, `recruiter`

---

## Users

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/users` | Public | List all users |
| GET | `/users/{id}` | Public | Get user by ID |
| POST | `/users` | Admin | Create user |
| PUT | `/users/{id}` | Admin | Update user |
| DELETE | `/users/{id}` | Admin | Delete user |

### Query Parameters — GET /users
| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number (default: 1) |
| `limit` | int | Results per page (default: 10) |
| `sort` | string | Sort field: `first_name`, `last_name`, `email`, `role`, `created_time` |
| `order` | string | `asc` or `desc` (default: `desc`) |
| `role` | string | Filter by role: `admin`, `candidate`, `recruiter` |
| `first_name` | string | Partial match filter |
| `last_name` | string | Partial match filter |
| `email` | string | Partial match filter |

---

## Jobs

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/jobs` | Public | List all jobs |
| GET | `/jobs/{id}` | Public | Get job by ID |
| GET | `/users/{userId}/jobs` | Public | Get jobs posted by a user |
| POST | `/jobs` | Admin / Recruiter | Create job |
| DELETE | `/jobs/{id}` | Admin / Recruiter | Delete job |

### Query Parameters — GET /jobs
| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number (default: 1) |
| `limit` | int | Results per page (default: 10) |
| `sort` | string | Sort field: `title`, `location`, `salary_min`, `created_time` |
| `order` | string | `asc` or `desc` (default: `desc`) |
| `title` | string | Partial match filter |
| `location` | string | Partial match filter |
| `job_type` | string | `full-time`, `part-time`, `contract`, `freelance` |
| `status` | string | `active`, `closed`, `draft` |

---

## Job Skills

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/jobs/{jobId}/skills` | Public | Get skills required for a job |
| GET | `/jobskills` | Admin | List all job skills |
| GET | `/jobskills/{id}` | Admin / Recruiter | Get job skill by ID |
| POST | `/jobskills` | Admin / Recruiter | Add skill requirement to a job |
| PUT | `/jobskills/{id}` | Admin / Recruiter | Update required proficiency level |
| DELETE | `/jobskills/{id}` | Admin / Recruiter | Remove skill from job |

---

## Applications

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/applications` | Admin | List all applications |
| GET | `/applications/{id}` | Admin / Candidate / Recruiter | Get application by ID |
| GET | `/users/{userId}/applications` | Admin / Candidate | Get applications by user |
| GET | `/jobs/{jobId}/applications` | Admin / Recruiter | Get applications for a job |
| POST | `/applications` | Admin / Candidate | Submit application |
| PUT | `/applications/{id}` | Admin / Recruiter | Update application status |
| DELETE | `/applications/{id}` | Admin / Candidate | Delete application |

### Application statuses
`applied` → `under_review` → `accepted` / `rejected` / `withdrawn`

---

## Candidate Skills

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/candidateskills` | Admin | List all candidate skills |
| GET | `/candidateskills/{id}` | Admin / Candidate | Get candidate skill by ID |
| GET | `/users/{userId}/skills` | Admin / Candidate / Recruiter | Get skills of a user |
| POST | `/candidateskills` | Admin / Candidate | Add skill to profile |
| PUT | `/candidateskills/{id}` | Admin / Candidate | Update proficiency level |
| DELETE | `/candidateskills/{id}` | Admin / Candidate | Remove skill from profile |

---

## Skills

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/skills` | Public | List all skills |
| GET | `/skills/{id}` | Public | Get skill by ID |
| POST | `/skills` | Admin | Create skill |
| PUT | `/skills/{id}` | Admin | Update skill |
| DELETE | `/skills/{id}` | Admin | Delete skill |

### Query Parameters — GET /skills
| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `limit` | int | Results per page |
| `sort` | string | Sort field |
| `order` | string | `asc` or `desc` |
| `name` | string | Partial match filter |

---

## Job Categories

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/jobcategories` | Public | List all job categories |
| GET | `/jobcategories/{id}` | Public | Get job category by ID |
| POST | `/jobcategories` | Admin | Create job category |
| PUT | `/jobcategories/{id}` | Admin | Update job category |
| DELETE | `/jobcategories/{id}` | Admin | Delete job category |

---

## Articles

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/articles` | Public | List all articles |
| GET | `/articles/{id}` | Public | Get article by ID |
| POST | `/articles` | Admin | Create article |
| PUT | `/articles/{id}` | Admin | Update article |
| DELETE | `/articles/{id}` | Admin | Delete article |

---

## Countries

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/countries` | Public | List all countries |
| GET | `/countries/{id}` | Public | Get country by ID |
| POST | `/countries` | Admin | Create country |
| PUT | `/countries/{id}` | Admin | Update country |
| DELETE | `/countries/{id}` | Admin | Delete country |

---

## Education Levels

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/educationlevels` | Public | List all education levels |
| GET | `/educationlevels/{id}` | Public | Get education level by ID |
| POST | `/educationlevels` | Admin | Create education level |
| PUT | `/educationlevels/{id}` | Admin | Update education level |
| DELETE | `/educationlevels/{id}` | Admin | Delete education level |

---

## Job Types

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/jobtypes` | Public | List all job types |
| GET | `/jobtypes/{id}` | Public | Get job type by ID |
| POST | `/jobtypes` | Admin | Create job type |
| PUT | `/jobtypes/{id}` | Admin | Update job type |
| DELETE | `/jobtypes/{id}` | Admin | Delete job type |

---

## Knowledge Levels

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/knowledgelevels` | Public | List all knowledge levels |
| GET | `/knowledgelevels/{id}` | Public | Get knowledge level by ID |
| POST | `/knowledgelevels` | Admin | Create knowledge level |
| PUT | `/knowledgelevels/{id}` | Admin | Update knowledge level |
| DELETE | `/knowledgelevels/{id}` | Admin | Delete knowledge level |

---

## Location Availabilities

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/locationavailabilities` | Public | List all location availabilities |
| GET | `/locationavailabilities/{id}` | Public | Get location availability by ID |
| POST | `/locationavailabilities` | Admin | Create location availability |
| PUT | `/locationavailabilities/{id}` | Admin | Update location availability |
| DELETE | `/locationavailabilities/{id}` | Admin | Delete location availability |

---

## Health Check

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/health` | Public | Check API health status |
