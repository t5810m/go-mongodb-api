# Database Schema

Database: `job_board` (MongoDB Atlas)

## Collections Overview

### users
Unified identity collection. A single user can be an `admin`, `candidate`, or `recruiter`.

```
_id:                 ObjectID
first_name:          string (required, min: 2, max: 100)
last_name:           string (required, min: 2, max: 100)
email:               string (required, unique)
password:            string (bcrypt hashed)
phone:               string
role:                string (admin | candidate | recruiter)
company_name:        string (recruiters only)
verified:            boolean
active:              boolean
terms_accepted:      boolean
last_terms_accepted: timestamp (nullable)
last_login_time:     timestamp (nullable)
created_time:        timestamp
updated_time:        timestamp
created_by:          string
updated_by:          string
```
**Indexes:** `email` (unique), `role`

---

### jobs
Job postings created by recruiters.

```
_id:          ObjectID
title:        string (required, min: 5, max: 255)
description:  string (required, min: 20)
user_id:      ObjectID (references users — recruiter)
category_id:  ObjectID (references jobcategories)
location:     string (required, min: 3)
job_type:     string (full-time | part-time | contract | freelance)
salary_min:   integer (required, > 0)
salary_max:   integer (required, > 0, >= salary_min)
status:       string (active | closed | draft)
active:       boolean
created_time: timestamp
updated_time: timestamp
created_by:   string
updated_by:   string
```
**Indexes:** `user_id`, `category_id`, `status`, `created_time` (desc)

---

### applications
Job applications submitted by candidates.

```
_id:            ObjectID
job_id:         ObjectID (references jobs)
user_id:        ObjectID (references users — candidate)
status:         string (applied | under_review | accepted | rejected | withdrawn)
recruiter_note: string
applied_time:   timestamp
updated_time:   timestamp
created_by:     string
updated_by:     string
```
**Indexes:** `job_id`, `user_id`, `status`, `{job_id + user_id}` (unique)

---

### candidateskills
Skills on a candidate's profile.

```
_id:               ObjectID
user_id:           ObjectID (references users — candidate)
skill_id:          ObjectID (references skills)
proficiency_level: string (beginner | intermediate | advanced | expert)
created_time:      timestamp
updated_time:      timestamp
created_by:        string
updated_by:        string
```
**Indexes:** `user_id`, `{user_id + skill_id}` (unique)

---

### jobskills
Skills required for a job posting.

```
_id:                        ObjectID
job_id:                     ObjectID (references jobs)
skill_id:                   ObjectID (references skills)
proficiency_level_required: string (beginner | intermediate | advanced | expert)
is_required:                boolean
created_time:               timestamp
updated_time:               timestamp
created_by:                 string
updated_by:                 string
```
**Indexes:** `job_id`, `{job_id + skill_id}` (unique)

---

### skills
Master list of technical skills and competencies.

```
_id:          ObjectID
name:         string (required, unique, min: 2, max: 100)
description:  string
created_time: timestamp
updated_time: timestamp
created_by:   string
updated_by:   string
```
**Indexes:** `name` (unique)

---

### jobcategories
Job classification categories.

```
_id:          ObjectID
name:         string (required, unique, min: 3, max: 100)
description:  string (required, min: 10, max: 500)
created_time: timestamp
updated_time: timestamp
created_by:   string
updated_by:   string
```
**Indexes:** `name` (unique)

---

### articles
Platform articles and blog posts.

```
_id:          ObjectID
title:        string (required, min: 2, max: 255)
content:      string (required)
slug:         string (required)
active:       boolean
created_time: timestamp
updated_time: timestamp
created_by:   string
updated_by:   string
```

---

### countries
Country reference data.

```
_id:          ObjectID
name:         string (required, min: 2, max: 100)
created_time: timestamp
updated_time: timestamp
created_by:   string
updated_by:   string
```

---

### educationlevels
Education level reference data.

```
_id:          ObjectID
title:        string (required, min: 2, max: 100)
created_time: timestamp
updated_time: timestamp
created_by:   string
updated_by:   string
```

---

### jobtypes
Job type reference data.

```
_id:          ObjectID
title:        string (required, min: 2, max: 100)
created_time: timestamp
updated_time: timestamp
created_by:   string
updated_by:   string
```

---

### knowledgelevels
Knowledge level reference data.

```
_id:          ObjectID
title:        string (required, min: 2, max: 100)
created_time: timestamp
updated_time: timestamp
created_by:   string
updated_by:   string
```

---

### locationavailabilities
Location/remote availability options.

```
_id:          ObjectID
title:        string (required, min: 2, max: 100)
created_time: timestamp
updated_time: timestamp
created_by:   string
updated_by:   string
```

---

### resumes
Candidate resume file references.

```
_id:           ObjectID
user_id:       ObjectID (references users — candidate)
file_url:      string (required, valid URL)
file_name:     string (required, min: 3)
uploaded_time: timestamp
updated_time:  timestamp
created_by:    string
updated_by:    string
```

---

## Data Relationships

```
Users (role=recruiter) (1) ──→ (many) Jobs
Users (role=candidate) (1) ──→ (many) Applications
Users (role=candidate) (1) ──→ (many) CandidateSkills
Users (role=candidate) (1) ──→ (many) Resumes
Jobs             (1) ──→ (many) Applications
Jobs             (1) ──→ (many) JobSkills
JobCategories    (1) ──→ (many) Jobs
Skills           (1) ──→ (many) JobSkills
Skills           (1) ──→ (many) CandidateSkills
```
