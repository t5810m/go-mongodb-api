# Testing the API with Postman

## 1. Start the API

```bash
go run ./cmd/main.go
```

The API starts on `http://localhost:8080`. You should see:
```
Connected to database: job_board
Starting API server on port 8080
Server is running. Press Ctrl+C to gracefully shutdown...
```

## 2. Seed the Database (first time only)

```bash
go run ./cmd/seed/seed.go
```

This populates the database with realistic test data. Seeded user credentials:

| Role | Email pattern | Password |
|------|--------------|----------|
| Admin | `admin1@jobsapi.com` – `admin3@jobsapi.com` | `SeedPassword123!` |
| Recruiter | `recruiter1@jobsapi.com` – `recruiter{N}@jobsapi.com` | `SeedPassword123!` |
| Candidate | `candidate1@jobsapi.com` – `candidate{N}@jobsapi.com` | `SeedPassword123!` |

> The seeder preserves existing reference data (skills, jobcategories, countries, etc.) and only re-seeds transactional collections (users, jobs, applications, etc.).

## 3. Import the Collection

In Postman: **Import** → select `docs/GO-MONGODB-API.postman_collection.json`

## 4. Set Up the Environment

1. In Postman go to **Environments** → create new environment named `GO-MONGODB-API (local)`
2. Add the following variables:

| Variable | Initial Value |
|----------|--------------|
| `baseUrl` | `http://localhost:8080` |
| `token` | *(leave empty)* |
| `userId` | *(leave empty)* |
| `jobId` | *(leave empty)* |
| `skillId` | *(leave empty)* |
| `categoryId` | *(leave empty)* |
| `applicationId` | *(leave empty)* |
| `candidateSkillId` | *(leave empty)* |
| `jobSkillId` | *(leave empty)* |
| `articleId` | *(leave empty)* |
| `countryId` | *(leave empty)* |
| `educationLevelId` | *(leave empty)* |
| `jobTypeId` | *(leave empty)* |
| `knowledgeLevelId` | *(leave empty)* |
| `locationAvailabilityId` | *(leave empty)* |

3. Select the environment from the top-right dropdown in Postman

## 5. Login and Capture Token

The **Auth → Login** request has a built-in Post-response script that automatically saves the JWT to the `token` environment variable after a successful login.

1. Open **Auth → Login**
2. The body is pre-filled with `admin1@jobsapi.com` / `SeedPassword123!`
3. Hit **Send**
4. The `token` variable is now set — all protected requests will use it automatically

To log in as a different role, change the email in the body:
- Recruiter: `recruiter1@jobsapi.com`
- Candidate: `candidate1@jobsapi.com`

## 6. Using Pagination, Sorting and Filters

All "Get All" requests include pre-filled query params in the **Params** tab, disabled by default. To use them:

1. Open any "Get All" request (e.g. **Users → Get All Users**)
2. Click the **Params** tab
3. Tick the checkbox next to any parameter to enable it
4. Set the value and hit **Send**

### Available params by endpoint

**GET /users**
- `page`, `limit`, `sort`, `order`
- `role` — `admin`, `candidate`, `recruiter`
- `first_name`, `last_name`, `email` — partial match

**GET /jobs**
- `page`, `limit`, `sort`, `order`
- `title`, `location` — partial match
- `job_type` — `full-time`, `part-time`, `contract`, `freelance`
- `status` — `active`, `closed`, `draft`

**GET /applications**
- `page`, `limit`, `sort`, `order`
- `status` — `applied`, `under_review`, `accepted`, `rejected`, `withdrawn`
- `job_id`, `user_id` — exact match

**GET /skills, /jobcategories, /articles, /countries, etc.**
- `page`, `limit`, `sort`, `order`
- `name` — partial match

## 7. Using ID Variables

When you need to test endpoints that require an ID (e.g. `GET /users/{userId}`):

1. Run **Get All Users** to see the list
2. Copy an `id` from the response
3. Set it as the `userId` environment variable
4. Run **Get User by ID** — it will use `{{userId}}` automatically

You can also use the Post-response script pattern on any Create request to auto-capture IDs:

```javascript
const json = pm.response.json();
if (json.id) {
    pm.environment.set("userId", json.id);
}
```

## 8. Suggested Testing Order

For full flow testing:

1. **Health Check** — verify the API is running
2. **Auth → Login** — get JWT token
3. **Skills → Get All Skills** — verify reference data
4. **Job Categories → Get All** — verify categories
5. **Users → Get All Users** — verify seeded users
6. **Jobs → Get All Jobs** — verify seeded jobs
7. **Applications → Get All** (admin token required)
8. **Candidate Skills → Get All** (admin token required)
9. Test CRUD operations for each resource

## Troubleshooting

**"Connection refused"**
- Make sure the API is running: `go run ./cmd/main.go`

**"Invalid credentials"**
- Check the email format: `admin1@jobsapi.com`
- Password: `SeedPassword123!` (note the `!`)

**"Unauthorized" on protected endpoints**
- Run Login first to get a fresh token
- Check the `token` environment variable is set
- Tokens expire after 24 hours

**"User not found" / "Not found" on ID endpoints**
- Make sure the `userId` (or other ID variable) is set in the environment
- Verify the ID exists by running the corresponding "Get All" first

**Duplicate key error on startup**
- The database has duplicate `{job_id, user_id}` pairs in applications
- Re-run the seeder: `go run ./cmd/seed/seed.go`
