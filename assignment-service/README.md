# assignment-service

Microservice for managing assignments (tugas) in the LMS platform.  
Built with **Go 1.22**, **Echo v4**, **MongoDB**, and **Redis Pub/Sub**.

---

## Tech Stack

| Layer        | Technology                        |
|--------------|-----------------------------------|
| Language     | Go 1.22                           |
| HTTP         | labstack/echo v4                  |
| Database     | MongoDB 7 (mongo-driver)          |
| Auth         | JWT validation (golang-jwt/jwt v5)|
| Events       | Redis Pub/Sub (go-redis v9)       |
| Container    | Docker + docker-compose           |

---

## Project Structure

```
assignment-service/
├── main.go                              ← Entry point
├── config/config.go                     ← Environment loading
├── handlers/assignment_handler.go       ← HTTP handlers
├── middleware/jwt.go                    ← JWT + role middleware
├── models/
│   ├── assignment.go                    ← Domain structs
│   └── request.go                       ← Request/validation structs
├── repositories/assignment_repository.go ← MongoDB queries
├── services/assignment_service.go       ← Business logic
├── utils/response.go                    ← Shared HTTP response helpers
├── pkg/redis/publisher.go               ← Redis event publisher
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

---

## Setup

### 1. Clone & configure

```bash
cp .env.example .env
# Edit .env — set JWT_SECRET to the same value used by user-service
```

### 2. Run with docker-compose (recommended)

```bash
docker-compose up --build
```

The service will be available at `http://localhost:8083`.

### 3. Run locally (without Docker)

Requires Go 1.22+, a running MongoDB instance, and a running Redis instance.

```bash
go mod download
go run .
```

### 4. Health check

```bash
curl http://localhost:8083/health
# {"status":"ok"}
```

---

## Environment Variables

| Variable       | Default                      | Description                     |
|----------------|------------------------------|---------------------------------|
| APP_PORT       | 8083                         | HTTP port                       |
| MONGO_URI      | mongodb://localhost:27017    | MongoDB connection URI           |
| MONGO_DB       | assignment_db                | MongoDB database name           |
| JWT_SECRET     | —                            | Shared secret with user-service |
| REDIS_URL      | (empty)                      | Redis.io URL, e.g. `rediss://...` |
| REDIS_ADDR     | localhost:6379               | Redis address                   |
| REDIS_USERNAME | (empty)                      | Redis username, e.g. `default`   |
| REDIS_PASSWORD | (empty)                      | Redis password (if any)         |
| REDIS_DB       | 0                            | Redis logical database          |

If `REDIS_URL` is set, the service uses that connection string and ignores `REDIS_ADDR`, `REDIS_USERNAME`, `REDIS_PASSWORD`, and `REDIS_DB`.

---

## Authentication

All endpoints require a valid JWT in the `Authorization` header:

```
Authorization: Bearer <token>
```

Tokens are issued by **user-service**. This service only **validates** them.

JWT claims expected:

- `user_id` (string) – UUID of the user
- `role` (string) – `"teacher"` | `"student"` | `"admin"`

---

## API Reference

Base path: `/api/v1/assignments`

---

### GET /api/v1/assignments?class_id=xxx

List all assignments for a class.  
**Roles:** teacher, student

**Query params:** `class_id` (required)

**Response 200:**

```json
{
  "success": true,
  "message": "success",
  "data": [
    {
      "id": "665f1a2b3c4d5e6f7a8b9c0d",
      "class_id": "class-uuid-001",
      "teacher_id": "teacher-uuid-001",
      "title": "Tugas 1: Essay",
      "description": "Tulis essay minimal 500 kata",
      "file_url": null,
      "deadline": "2024-12-01T23:59:59Z",
      "created_at": "2024-11-01T10:00:00Z",
      "updated_at": "2024-11-01T10:00:00Z"
    }
  ]
}
```

---

### GET /api/v1/assignments/:id

Get assignment detail.  
**Roles:** teacher, student

**Response 200:**

```json
{
  "success": true,
  "message": "success",
  "data": {
    "id": "665f1a2b3c4d5e6f7a8b9c0d",
    "class_id": "class-uuid-001",
    "teacher_id": "teacher-uuid-001",
    "title": "Tugas 1: Essay",
    "description": "Tulis essay minimal 500 kata",
    "file_url": null,
    "deadline": "2024-12-01T23:59:59Z",
    "created_at": "2024-11-01T10:00:00Z",
    "updated_at": "2024-11-01T10:00:00Z"
  }
}
```

---

### POST /api/v1/assignments

Create a new assignment.  
**Roles:** teacher only

**Request body:**

```json
{
  "class_id": "class-uuid-001",
  "title": "Tugas 1: Essay",
  "description": "Tulis essay minimal 500 kata tentang Pancasila",
  "file_url": null,
  "deadline": "2024-12-01T23:59:59Z"
}
```

**Response 201:**

```json
{
  "success": true,
  "message": "success",
  "data": {
    "id": "665f1a2b3c4d5e6f7a8b9c0d",
    "class_id": "class-uuid-001",
    "teacher_id": "teacher-uuid-001",
    "title": "Tugas 1: Essay",
    "description": "Tulis essay minimal 500 kata tentang Pancasila",
    "file_url": null,
    "deadline": "2024-12-01T23:59:59Z",
    "created_at": "2024-11-01T10:00:00Z",
    "updated_at": "2024-11-01T10:00:00Z"
  }
}
```

---

### PUT /api/v1/assignments/:id

Update an assignment (only the creating teacher can do this).  
**Roles:** teacher only

**Request body:**

```json
{
  "title": "Tugas 1: Essay (Revisi)",
  "description": "Tulis essay minimal 700 kata",
  "file_url": "https://storage.example.com/rubrik.pdf",
  "deadline": "2024-12-05T23:59:59Z"
}
```

**Response 200:** updated assignment object (same shape as GET detail)

---

### DELETE /api/v1/assignments/:id

Delete an assignment (only the creating teacher can do this).  
**Roles:** teacher only

**Response 200:**

```json
{
  "success": true,
  "message": "success",
  "data": null
}
```

---

### POST /api/v1/assignments/:id/submit

Submit an assignment.  
**Roles:** student only  
**Rules:**

- Only one submission per student per assignment
- Rejected if current time is past the deadline

**Request body:**

```json
{
  "file_url": "https://storage.example.com/essay-saya.pdf",
  "content": null
}
```

**Response 201:**

```json
{
  "success": true,
  "message": "success",
  "data": {
    "id": "665f2b3c4d5e6f7a8b9c0d1e",
    "assignment_id": "665f1a2b3c4d5e6f7a8b9c0d",
    "student_id": "student-uuid-001",
    "file_url": "https://storage.example.com/essay-saya.pdf",
    "content": null,
    "grade": null,
    "status": "submitted",
    "submitted_at": "2024-11-15T09:30:00Z",
    "graded_at": null,
    "created_at": "2024-11-15T09:30:00Z",
    "updated_at": "2024-11-15T09:30:00Z"
  }
}
```

---

### GET /api/v1/assignments/:id/my-submission

Get your own submission for an assignment.  
**Roles:** student only

**Response 200:** submission object (same shape as submit response)

---

### GET /api/v1/assignments/:id/submissions

List all submissions for an assignment.  
**Roles:** teacher only (and only the teacher who created the assignment)

**Response 200:**

```json
{
  "success": true,
  "message": "success",
  "data": [
    {
      "id": "665f2b3c4d5e6f7a8b9c0d1e",
      "assignment_id": "665f1a2b3c4d5e6f7a8b9c0d",
      "student_id": "student-uuid-001",
      "file_url": "https://storage.example.com/essay-saya.pdf",
      "content": null,
      "grade": null,
      "status": "submitted",
      "submitted_at": "2024-11-15T09:30:00Z",
      "graded_at": null,
      "created_at": "2024-11-15T09:30:00Z",
      "updated_at": "2024-11-15T09:30:00Z"
    }
  ]
}
```

---

### POST /api/v1/assignments/:id/grade/:submission_id

Grade a student's submission.  
**Roles:** teacher only (and only the teacher who created the assignment)

**Request body:**

```json
{
  "grade": 87.5
}
```

**Response 200:**

```json
{
  "success": true,
  "message": "success",
  "data": {
    "id": "665f2b3c4d5e6f7a8b9c0d1e",
    "assignment_id": "665f1a2b3c4d5e6f7a8b9c0d",
    "student_id": "student-uuid-001",
    "file_url": "https://storage.example.com/essay-saya.pdf",
    "content": null,
    "grade": 87.5,
    "status": "graded",
    "submitted_at": "2024-11-15T09:30:00Z",
    "graded_at": "2024-11-20T14:00:00Z",
    "created_at": "2024-11-15T09:30:00Z",
    "updated_at": "2024-11-20T14:00:00Z"
  }
}
```

---

## Error Responses

All errors follow this shape:

```json
{
  "success": false,
  "message": "error description",
  "data": null
}
```

| HTTP Status | Scenario                                  |
|-------------|-------------------------------------------|
| 400         | Missing / invalid request body or params  |
| 401         | Missing or invalid JWT                    |
| 403         | Insufficient role or not the owner        |
| 404         | Resource not found                        |
| 409         | Student already submitted                 |
| 422         | Submission past deadline                  |
| 500         | Internal server error                     |

---

## Redis Events

Events are published via Redis Pub/Sub. Downstream consumers (e.g. progress-service) subscribe to these channels.

| Channel                     | Trigger                      |
|-----------------------------|------------------------------|
| `lms:assignment:created`    | Teacher creates assignment   |
| `lms:assignment:submitted`  | Student submits              |
| `lms:assignment:graded`     | Teacher grades submission    |

### lms:assignment:created

```json
{
  "event": "ASSIGNMENT_CREATED",
  "assignment_id": "665f1a2b3c4d5e6f7a8b9c0d",
  "class_id": "class-uuid-001",
  "title": "Tugas 1: Essay",
  "deadline": "2024-12-01T23:59:59Z",
  "teacher_id": "teacher-uuid-001",
  "timestamp": "2024-11-01T10:00:00Z"
}
```

### lms:assignment:submitted

```json
{
  "event": "ASSIGNMENT_SUBMITTED",
  "submission_id": "665f2b3c4d5e6f7a8b9c0d1e",
  "assignment_id": "665f1a2b3c4d5e6f7a8b9c0d",
  "student_id": "student-uuid-001",
  "class_id": "class-uuid-001",
  "timestamp": "2024-11-15T09:30:00Z"
}
```

### lms:assignment:graded

```json
{
  "event": "ASSIGNMENT_GRADED",
  "submission_id": "665f2b3c4d5e6f7a8b9c0d1e",
  "assignment_id": "665f1a2b3c4d5e6f7a8b9c0d",
  "student_id": "student-uuid-001",
  "grade": 87.5,
  "timestamp": "2024-11-20T14:00:00Z"
}
```

---

## MongoDB Indexes

Created automatically on startup:

**assignments:**

- `class_id` (ascending)
- `teacher_id` (ascending)

**submissions:**

- `assignment_id` (ascending)
- `student_id` (ascending)
- `(assignment_id, student_id)` — **unique compound index** (prevents duplicate submissions)
