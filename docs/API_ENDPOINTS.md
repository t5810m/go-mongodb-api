# API Endpoints

## Jobs

```
GET    /jobs                              # List all jobs (pagination, filtering, sorting)
GET    /jobs/{id}                         # Get job details
POST   /jobs                              # Create new job
GET    /companies/{companyId}/jobs        # Get company's jobs
DELETE /jobs/{id}                         # Delete job
```

## Candidates

```
GET    /candidates                        # List all candidates
GET    /candidates/{id}                   # Get candidate profile
POST   /candidates                        # Create new candidate
GET    /candidates/{candidateId}/skills   # Get candidate's skills
DELETE /candidates/{id}                   # Delete candidate
```

## Recruiters

```
GET    /recruiters                        # List all recruiters
GET    /recruiters/{id}                   # Get recruiter profile
POST   /recruiters                        # Create new recruiter
DELETE /recruiters/{id}                   # Delete recruiter
```

## Companies

```
GET    /companies                         # List all companies
GET    /companies/{id}                    # Get company details
POST   /companies                         # Create new company
DELETE /companies/{id}                    # Delete company
```

## Applications

```
GET    /applications                      # List all applications
GET    /applications/{id}                 # Get application details
POST   /applications                      # Submit job application
PUT    /applications/{id}                 # Update application status
GET    /jobs/{jobId}/applications         # Get applications for a job
GET    /candidates/{candidateId}/applications  # Get candidate's applications
DELETE /applications/{id}                 # Delete application
```

## Skills

```
GET    /skills                            # List all skills
GET    /skills/{id}                       # Get skill details
POST   /skills                            # Create new skill
GET    /jobs/{jobId}/skills               # Get required skills for a job
DELETE /skills/{id}                       # Delete skill
```

## Job Categories & Job Skills

```
GET    /jobcategories                     # List all job categories
POST   /jobcategories                     # Create job category
POST   /jobskills                         # Add required skill to job
GET    /jobskills/{id}                    # Get job skill details
DELETE /jobcategories/{id}                # Delete job category
DELETE /jobskills/{id}                    # Delete job skill
```

## Candidate Skills

```
GET    /candidateskills                   # List all candidate skills
POST   /candidateskills                   # Add skill to candidate profile
PUT    /candidateskills/{id}              # Update skill proficiency
GET    /candidates/{candidateId}/skills   # Get candidate's skills
DELETE /candidateskills/{id}              # Delete candidate skill
```

## Query Parameters

### Pagination
```
GET /jobs?page=1&limit=10
```

### Filtering
```
GET /jobs?title=React&location=Remote&status=active
GET /candidates?first_name=John&email=john@example.com
```

### Sorting
```
GET /jobs?sort=created_time&order=desc
GET /candidates?sort=first_name&order=asc
```

### Combined
```
GET /jobs?page=1&limit=10&title=Senior&sort=created_time&order=desc
```
