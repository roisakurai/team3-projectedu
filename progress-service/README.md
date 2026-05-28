# progress-service

A microservice for the LMS (Learning Management System) that tracks and calculates student learning progress. It consumes events from other services via Redis Pub/Sub and exposes a REST API for dashboards and reports.

---

## Tech Stack

| Layer         | Technology                          |
|---------------|-------------------------------------|
| Language      | Go 1.22                             |
| Framework     | Echo v4 (labstack/echo)             |
| Database      | MongoDB 7 (mongo-driver)            |
| Auth          | JWT validation only (golang-jwt/v5) |
| Events        | Redis 7 Pub/Sub (go-redis/v9)       |
| Container     | Docker + docker-compose             |

---

## Setup & Run

### Prerequisites
- Docker & docker-compose installed
- (Optional) Go 1.22+ for local development

### 1. Clone & configure

```bash
git clone <repo-url>
cd progress-service
cp .env.example .env
# Edit .env and set JWT_SECRET to match your user-service secret
```

### 2. Run with docker-compose

```bash
docker-compose up --build
```

This starts three containers:
- `progress-service` on port **8085**
- `progress-mongodb` on port **27017**
- `progress-redis` on port **6379**

### 3. Run locally (without Docker)

Make sure MongoDB and Redis are running locally, then:

```bash
go mod download
go run ./cmd/main.go
```

---

## API Endpoints

All endpoints require a valid JWT in the `Authorization: Bearer <token>` header.

### Health Check

```
GET /health
Response: { "status": "ok" }
```

---

### Student Endpoints (role: `student`)

#### GET /api/v1/progress/me
Get own progress across all classes.

```json
// Response
{
  "success": true,
  "message": "success",
  "data": [
    {
      "id": "665a1b2c3d4e5f6789abcdef",
      "student_id": "stu_001",
      "class_id": "cls_001",
      "total_materials": 10,
      "completed_materials": 7,
      "total_assignments": 5,
      "submitted_assignments": 4,
      "graded_assignments": 3,
      "average_grade": 85.5,
      "last_activity_at": "2024-06-01T10:00:00Z",
      "updated_at": "2024-06-01T10:00:00Z",
      "created_at": "2024-01-15T08:00:00Z"
    }
  ]
}
```

#### GET /api/v1/progress/me/class/:class_id
Get own progress in a specific class.

```json
// Response
{
  "success": true,
  "message": "success",
  "data": {
    "class_id": "cls_001",
    "material_progress": {
      "total": 10,
      "completed": 7,
      "percentage": 70.0
    },
    "assignment_progress": {
      "total": 5,
      "submitted": 4,
      "graded": 3,
      "average_grade": 85.5,
      "submission_percentage": 80.0
    },
    "overall_percentage": 75.0
  }
}
```

#### GET /api/v1/progress/me/dashboard
Full dashboard with summary across all classes.

```json
// Response
{
  "success": true,
  "message": "success",
  "data": {
    "student_id": "stu_001",
    "classes": [
      {
        "class_id": "cls_001",
        "material_progress": {
          "total": 10,
          "completed": 7,
          "percentage": 70.0
        },
        "assignment_progress": {
          "total": 5,
          "submitted": 4,
          "graded": 3,
          "average_grade": 85.5,
          "submission_percentage": 80.0
        },
        "overall_percentage": 75.0
      }
    ],
    "global_average_grade": 85.5,
    "last_activity_at": "2024-06-01T10:00:00Z"
  }
}
```

---

### Teacher Endpoints (role: `teacher`)

#### GET /api/v1/progress/class/:class_id
Get progress of all students in a class.

```json
// Response
{
  "success": true,
  "message": "success",
  "data": [
    {
      "id": "665a1b2c3d4e5f6789abcdef",
      "student_id": "stu_001",
      "class_id": "cls_001",
      "total_materials": 10,
      "completed_materials": 7,
      "total_assignments": 5,
      "submitted_assignments": 4,
      "graded_assignments": 3,
      "average_grade": 85.5,
      "last_activity_at": "2024-06-01T10:00:00Z",
      "updated_at": "2024-06-01T10:00:00Z",
      "created_at": "2024-01-15T08:00:00Z"
    }
  ]
}
```

#### GET /api/v1/progress/class/:class_id/student/:student_id
Get detailed progress of a specific student in a class.

```json
// Response
{
  "success": true,
  "message": "success",
  "data": {
    "class_id": "cls_001",
    "material_progress": { "total": 10, "completed": 7, "percentage": 70.0 },
    "assignment_progress": {
      "total": 5, "submitted": 4, "graded": 3,
      "average_grade": 85.5, "submission_percentage": 80.0
    },
    "overall_percentage": 75.0
  }
}
```

---

### Admin Endpoints (role: `admin`)

#### GET /api/v1/progress/summary
Platform-wide statistics.

```json
// Response
{
  "success": true,
  "message": "success",
  "data": {
    "total_students": 150,
    "total_classes": 12,
    "platform_average_grade": 82.3,
    "total_submissions": 847
  }
}
```

---

## Error Response Format

```json
{
  "success": false,
  "message": "error description",
  "data": null
}
```

| HTTP Code | Reason                          |
|-----------|---------------------------------|
| 400       | Bad request / missing params    |
| 401       | Missing or invalid JWT          |
| 403       | Insufficient role               |
| 500       | Internal server error           |

---

## Redis Event Channel Reference

The service subscribes to 6 channels on startup, each in its own goroutine.

### `lms:class:student_joined`
Initializes a new `student_progress` record for the student in the class.
```json
{
  "event": "STUDENT_JOINED_CLASS",
  "student_id": "stu_001",
  "class_id": "cls_001",
  "timestamp": "2024-01-15T08:00:00Z"
}
```

### `lms:material:created`
Increments `total_materials` for all enrolled students; inserts `material_progress` (pending) for each.
```json
{
  "event": "MATERIAL_CREATED",
  "material_id": "mat_001",
  "class_id": "cls_001",
  "title": "Introduction to Go",
  "teacher_id": "tea_001",
  "timestamp": "2024-01-16T09:00:00Z"
}
```

### `lms:material:completed`
Marks `material_progress` as completed; increments `completed_materials`; updates `last_activity_at`.
```json
{
  "event": "MATERIAL_COMPLETED",
  "student_id": "stu_001",
  "material_id": "mat_001",
  "class_id": "cls_001",
  "timestamp": "2024-01-17T11:00:00Z"
}
```

### `lms:assignment:created`
Increments `total_assignments` for all enrolled students; inserts `assignment_progress` with status `pending` for each.
```json
{
  "event": "ASSIGNMENT_CREATED",
  "assignment_id": "asg_001",
  "class_id": "cls_001",
  "title": "Go Concurrency Exercise",
  "deadline": "2024-02-01T23:59:59Z",
  "teacher_id": "tea_001",
  "timestamp": "2024-01-18T10:00:00Z"
}
```

### `lms:assignment:submitted`
Upserts `assignment_progress` to `submitted`; increments `submitted_assignments`.
```json
{
  "event": "ASSIGNMENT_SUBMITTED",
  "submission_id": "sub_001",
  "assignment_id": "asg_001",
  "student_id": "stu_001",
  "class_id": "cls_001",
  "timestamp": "2024-01-30T14:00:00Z"
}
```

### `lms:assignment:graded`
Updates `assignment_progress` to `graded`; inserts `grade_record`; increments `graded_assignments`; recalculates `average_grade`; updates `last_activity_at`.
```json
{
  "event": "ASSIGNMENT_GRADED",
  "submission_id": "sub_001",
  "assignment_id": "asg_001",
  "student_id": "stu_001",
  "class_id": "cls_001",
  "grade": 92.5,
  "timestamp": "2024-02-03T16:00:00Z"
}
```

---

## MongoDB Collections & Indexes

| Collection           | Unique Index                                |
|----------------------|---------------------------------------------|
| `student_progress`   | `{ student_id: 1, class_id: 1 }`           |
| `material_progress`  | `{ student_id: 1, material_id: 1 }`        |
| `assignment_progress`| `{ student_id: 1, assignment_id: 1 }`      |
| `grade_records`      | `{ student_id: 1 }` (non-unique)           |

All indexes are created automatically at service startup.

---

## Project Structure

```
progress-service/
├── cmd/main.go                          ← entry point
├── config/config.go                     ← env loading
├── internal/
│   ├── handler/progress_handler.go      ← HTTP handlers
│   ├── service/progress_service.go      ← business logic
│   ├── repository/progress_repository.go← MongoDB queries
│   ├── model/                           ← data structs
│   ├── middleware/jwt_middleware.go      ← JWT auth + role guard
│   └── consumer/                        ← Redis event consumers
├── pkg/redis/subscriber.go              ← Redis client helper
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```
