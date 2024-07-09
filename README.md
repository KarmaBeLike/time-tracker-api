# Time Tracker API

This repository implements a Time Tracker API using Go (Golang) with the Gin framework and PostgreSQL database.


## Features

### User Management
- Create a new user based on passport number.
- Retrieve a list of users with optional filters and pagination.
- Update and delete users by ID.

### Task Management
- Retrieve worklogs (tasks) for a user with optional date range.
- Start and stop tasks for a user by ID and task ID.

## API Endpoints

### Users
- POST /users/ - Create a new user
- GET /users/ - Get list of users
- DELETE /users/{userId} - Delete a user
- PUT /users/{userId} - Update a user

### Tasks
- GET /users/{userId}/tasks?startDate={startDate}&endDate={endDate} - Get list of tasks for a user
- POST /users/{userId}/tasks/{taskId}/start - Start a task for a user
- POST /users/{userId}/tasks/{taskId}/stop - End a task for a user

## Getting Started
 **install dependencies:**
```
go mod tidy
```
- **run project:**
```
go run ./cmd
```