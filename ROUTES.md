# Routes Summary

This file lists HTTP routes across the project services for quick reference.

---

## user-service

Base: <http://localhost:8080> (default)

- GET /swagger/index.html — Swagger UI
- POST /register — Register new user
- POST /login — Login (returns JWT token)
- GET /verify-email/:token — Email verification
- GET /test — Health/test
- GET /profile — Get profile (protected: requires JWT)
- PUT /users/:id — Update user (protected)
- DELETE /users/:id — Delete user (protected)

---

## class-service

Base: <http://localhost:8082> (default)
All protected by JWT (middleware).

- GET /swagger/index.html — Swagger UI
- POST /classes — Create class (role: teacher)
  - Body example: { "name":"Biology 101", "description":"Intro", "join_code":"BIO101" }
- GET /classes — Get all classes (protected)
- POST /classes/join — Join class (role: student)
  - Body example: { "join_code":"BIO101" }
- GET /classes/:id/students — Get students by class

---

## material-service

Base: (see service .env)

- GET /swagger/index.html — Swagger UI
- POST /materials — Create material (protected)
- PUT /materials/:id/read — Mark material as read (protected)
- GET /materials — Get all materials (public)
- GET /materials/class/:class_id — Get materials by class (public)

---

## assignment-service

Base: /api/v1/assignments (protected with JWT middleware)

- GET /swagger/index.html — Swagger UI
- GET /health — Health
- GET /api/v1/assignments?class_id=... — List assignments (requires class_id query)
- GET /api/v1/assignments/:id — Get assignment
- POST /api/v1/assignments — Create assignment (role: teacher)
- PUT /api/v1/assignments/:id — Update assignment (role: teacher)
- DELETE /api/v1/assignments/:id — Delete assignment (role: teacher)
- GET /api/v1/assignments/:id/submissions — List submissions (role: teacher)
- POST /api/v1/assignments/:id/grade/:submission_id — Grade submission (role: teacher)
- POST /api/v1/assignments/:id/submit — Submit assignment (role: student)
- GET /api/v1/assignments/:id/my-submission — Get my submission (role: student)

---

## progress-service

Base: /api/v1/progress (protected with JWT middleware)

- GET /swagger/index.html — Swagger UI
- GET /health — Health
- Student routes (role: student):
  - GET /api/v1/progress/me
  - GET /api/v1/progress/me/class/:class_id
  - GET /api/v1/progress/me/dashboard
- Teacher routes (role: teacher):
  - GET /api/v1/progress/class/:class_id
  - GET /api/v1/progress/class/:class_id/student/:student_id
- Admin routes (role: admin):
  - GET /api/v1/progress/summary

---

## notification-service

Base: <http://localhost:8086> (default)

- GET /swagger/index.html — Swagger UI
- GET /health — Health
- POST /notifications — Create notification
- POST /notifications/email — Create notification and send email
- GET /notifications/me — Get my notifications (protected)
- PATCH /notifications/:id/read — Mark notification as read (protected)

---

## Notes

- All services expect `Authorization: Bearer <token>` for protected routes.
- `user-service` issues tokens signed with `JWT_SECRET` environment variable — other services must use the same secret for local token validation when applicable.
- If you want, I can also generate a Postman collection JSON for these routes.
