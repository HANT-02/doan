# API Endpoints Map
This file serves as the source-of-truth for the mapping between frontend calls and the backend API routes.

## Base URL
Prefix: `/v1` or `/v2` for specific user routes.

## Authentication
- `POST /v1/auth/login`: Login and receive access_token + refresh_token
- `POST /v1/auth/logout`: Logout user
- `POST /v1/auth/refresh`: Refresh access token
- `POST /v1/auth/register`: Register a new user
- `POST /v1/auth/forgot-password`: Request password reset
- `POST /v1/auth/reset-password`: Reset password
- `POST /v1/auth/change-password`: Change password
- `POST /v1/auth/verify-otp`: Verify OTP
- `GET /v1/auth/me`: Get current logged-in user profile
- `POST /v2/auth/login`: Login V2
- `POST /v2/auth/logout`: Logout V2
- `POST /v2/auth/refresh`: Refresh V2

## Teachers
- `GET /v1/teachers`: List all teachers (supports pagination/filtering)
- `GET /v1/teachers/:id`: Get teacher details
- `POST /v1/teachers`: Create a new teacher
- `PUT /v1/teachers/:id`: Update a teacher
- `DELETE /v1/teachers/:id`: Soft delete a teacher
- `GET /v1/teachers/:id/timetable`: Get teacher timetable
- `GET /v1/teachers/:id/stats/teaching-hours`: Get teaching hours stats

## Rooms
- `GET /v1/rooms`: List all rooms
- `GET /v1/rooms/:id`: Get room details
- `POST /v1/rooms`: Create a new room
- `PUT /v1/rooms/:id`: Update a room
- `DELETE /v1/rooms/:id`: Soft delete a room

## Classes
- `GET /v1/classes`: List all classes
- `GET /v1/classes/:id`: Get class details
- `GET /v1/classes/:id/students`: Get class roster and capacity snapshot
- `POST /v1/classes`: Create a new class
- `PUT /v1/classes/:id`: Update a class
- `PUT /v1/classes/:id/teacher`: Assign teacher for a class
- `DELETE /v1/classes/:id`: Soft delete a class
- `POST /v1/classes/:id/students`: Add students into a class
- `DELETE /v1/classes/:id/students`: Remove students from a class

## Students
- `GET /v1/students`: List all students
- `GET /v1/students/:id`: Get student details
- `POST /v1/students`: Create a new student
- `PUT /v1/students/:id`: Update a student
- `DELETE /v1/students/:id`: Soft delete a student

## Courses
- `GET /v1/courses`: List all courses
- `GET /v1/courses/:id`: Get course details
- `POST /v1/courses`: Create a new course
- `PUT /v1/courses/:id`: Update a course
- `DELETE /v1/courses/:id`: Soft delete a course

## Programs
- `GET /v1/programs`: List all programs
- `GET /v1/programs/:id`: Get program details
- `POST /v1/programs`: Create a new program
- `PUT /v1/programs/:id`: Update a program
- `DELETE /v1/programs/:id`: Soft delete a program
- `POST /v1/programs/:id/courses`: Add courses to program
- `DELETE /v1/programs/:id/courses`: Remove courses from program

## Scheduling
- `POST /v1/scheduling/preview`: Trigger CSP scheduling preview
- `GET /v1/scheduling/preview/latest`: Get the latest preview result
- `GET /v1/scheduling/preview/:id`: Get a preview result by run ID
- `POST /v1/scheduling/commit`: Commit scaffold for a preview run

## Materials / AI Audit
- `POST /v1/materials/upload`: Upload teaching material and trigger stub AI audit
- `GET /v1/materials`: List materials with optional `teacher_id`, `status`, `queue`
- `GET /v1/materials/flagged`: List materials waiting for compliance review
- `GET /v1/materials/:id`: Get material detail with audit trail
- `POST /v1/materials/:id/review`: Approve or reject a material

## Data Structures Notes
- Reponses typically format as `{ "message": "...", "data": ... }`
- List endpoints return: `{ "data": { "teachers": [...], "pagination": { "current_page": 1, "total_pages": x, "total_items": x, "items_per_page": x } } }`
- `GET /v1/classes` currently returns PascalCase payload from Go usecase wrapper: `{ "data": { "Classes": [...], "Pagination": { "CurrentPage": 1, "ItemsPerPage": 10, "TotalItems": x, "TotalPages": x } } }`
- `GET /v1/rooms` currently returns PascalCase payload from Go usecase wrapper: `{ "data": { "Rooms": [...], "Pagination": { "CurrentPage": 1, "ItemsPerPage": 10, "TotalItems": x, "TotalPages": x } } }`
- `GET /v1/programs` returns lowercase DTO payload: `{ "data": { "programs": [...], "pagination": { "current_page": 1, "items_per_page": 10, "total_items": x, "total_pages": x } } }`
- `GET /v1/programs/:id` returns the program object directly in `data`, with embedded `courses` when linked.
- Authentication endpoints MUST return Tokens with `Authorization: Bearer <token>` in the future requests.
