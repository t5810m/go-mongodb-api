# Database Schema

## Collections Overview

### Users
Site administrators and system users.
```
_id: ObjectID
email: string (unique)
password: string (hashed)
role: string (admin, moderator, user)
created_time: timestamp
updated_time: timestamp
```

### Companies
Hiring companies/organizations.
```
_id: ObjectID
name: string (required)
description: string
website: string
location: string
industry: string
employees_count: integer
created_time: timestamp
updated_time: timestamp
created_by: string
updated_by: string
```

### Job Categories
Job classification categories.
```
_id: ObjectID
name: string (required, unique)
description: string
created_time: timestamp
updated_time: timestamp
created_by: string
updated_by: string
```

### Skills
Technical skills and competencies.
```
_id: ObjectID
name: string (required, unique)
description: string
category: string
proficiency_levels: array of strings
created_time: timestamp
updated_time: timestamp
created_by: string
updated_by: string
```

### Recruiters
Job posters (independent or company-based).
```
_id: ObjectID
email: string (required, unique)
password: string (hashed)
first_name: string (required)
last_name: string (required)
phone: string
company_id: ObjectID (optional, references companies)
job_postings_count: integer
created_time: timestamp
updated_time: timestamp
created_by: string
updated_by: string
```

### Candidates
Job seekers.
```
_id: ObjectID
email: string (required, unique)
password: string (hashed)
first_name: string (required)
last_name: string (required)
phone: string
location: string
bio: string
years_experience: integer
created_time: timestamp
updated_time: timestamp
created_by: string
updated_by: string
```

### Jobs
Job postings.
```
_id: ObjectID
title: string (required, min: 5)
description: string (required, min: 20)
recruiter_id: ObjectID (required, references recruiters)
company_id: ObjectID (required, references companies)
category_id: ObjectID (required, references jobcategories)
location: string (required)
job_type: string (full-time, part-time, contract, freelance)
salary_min: integer (required)
salary_max: integer (required, >= salary_min)
status: string (active, closed, draft)
active: boolean
created_time: timestamp
updated_time: timestamp
created_by: string
updated_by: string
```

### Applications
Job applications.
```
_id: ObjectID
candidate_id: ObjectID (required, references candidates)
job_id: ObjectID (required, references jobs)
status: string (pending, accepted, rejected, withdrawn)
applied_time: timestamp (when application was submitted)
created_time: timestamp
updated_time: timestamp
created_by: string
updated_by: string
```

### Job Skills
Required skills for jobs.
```
_id: ObjectID
job_id: ObjectID (required, references jobs)
skill_id: ObjectID (required, references skills)
proficiency_level: string (beginner, intermediate, advanced, expert)
years_required: integer
created_time: timestamp
updated_time: timestamp
created_by: string
updated_by: string
```

### Candidate Skills
Skills of candidates.
```
_id: ObjectID
candidate_id: ObjectID (required, references candidates)
skill_id: ObjectID (required, references skills)
proficiency_level: string (beginner, intermediate, advanced, expert)
years_experience: integer
created_time: timestamp
updated_time: timestamp
created_by: string
updated_by: string
```

### Resumes
Candidate resumes/CVs.
```
_id: ObjectID
candidate_id: ObjectID (required, references candidates)
title: string
content: string
url: string
is_primary: boolean
created_time: timestamp
updated_time: timestamp
```

## Data Relationships

```
Companies (1) ──→ (many) Recruiters
Companies (1) ──→ (many) Jobs
Recruiters (1) ──→ (many) Jobs
Jobs (1) ──→ (many) Applications
Jobs (1) ──→ (many) Job Skills
Candidates (1) ──→ (many) Applications
Candidates (1) ──→ (many) Candidate Skills
Candidates (1) ──→ (many) Resumes
Skills (1) ──→ (many) Job Skills
Skills (1) ──→ (many) Candidate Skills
Job Categories (1) ──→ (many) Jobs
```
