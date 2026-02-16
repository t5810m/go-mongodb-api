# Testing the API with Postman

This guide helps you test the api-m API using the pre-built Postman collection.

## Setup

### 1. Start the API

Make sure the API is running on your machine.
In command line after you enter the folder where you have downloaded the project, enter:

```bash
go run ./cmd/main.go
```

The API will start on `http://localhost:8080`

### 2. Import the Collection in Postman
Click on Import button and browse to the docs folder in this project, and import the .json file.


## Troubleshooting

**"Connection refused" error:**
- Make sure the API is running with `go run ./cmd/main.go`
- Check that you're using `localhost:8080`

**"Failed to create" error:**
- Check the response for validation errors
- Make sure you're following the testing order (dependencies matter)
