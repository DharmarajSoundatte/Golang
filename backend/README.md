# School CRM — Backend API

Go REST API for the School CRM platform.  
**Stack:** Go · Chi · GORM · PostgreSQL (Aiven) · JWT · Zap

---

## Team & Module Ownership

To prevent git conflicts every file in this project has a clear owner.  
**Do not edit files outside your assigned modules without first talking to the other person.**

| Developer | GitHub | Assigned Modules |
|-----------|--------|------------------|
| **Dharmaraj** | `DharmarajSoundatte` | Auth, Students, Parents, Fees |
| **Sanjana** | `sanjana` | Teachers, Classes, Attendance, Grades, Timetable, Announcements |

---

## Module Breakdown

### Dharmaraj's Modules

#### 1. Auth & User Management ✅ (done)
- Register, Login, JWT middleware
- Role-based access (`admin`, `teacher`, `student`, `parent`)

#### 2. Student Management
- Enroll / update / deactivate students
- Link student → class, parent
- Search & filter students

#### 3. Parent Management
- Register parents
- Link parent → one or more students
- Parent portal view (read-only)

#### 4. Fee Management
- Define fee structures per class/term
- Record payments
- Outstanding dues report
- Payment history per student

---

### Sanjana's Modules

#### 5. Teacher Management
- Add / update / deactivate teachers
- Assign subjects to teachers

#### 6. Class & Section Management
- Create classes (`Grade 1`, `Grade 2` …) and sections (`A`, `B`)
- Assign class teacher
- Assign students to class/section

#### 7. Attendance Management
- Mark daily attendance per student per class
- Teacher can mark / edit attendance
- Monthly attendance summary

#### 8. Grades & Marks
- Create exams / assessments
- Enter marks per student per subject
- Generate report cards

#### 9. Timetable
- Define periods per class per day
- Assign teacher + subject to each period

#### 10. Announcements
- Admin / teacher can post announcements
- Target audience: all / specific class / role

---

## File Ownership Map

```
backend/
├── cmd/server/main.go                        SHARED  (do not modify without discussion)
├── internal/
│   ├── config/config.go                      SHARED
│   ├── database/postgres.go                  SHARED
│   ├── middleware/                            SHARED
│   ├── routes/
│   │   ├── routes.go                         SHARED  (only registers sub-routers)
│   │   ├── dharmaraj_routes.go               DHARMARAJ
│   │   └── sanjana_routes.go                 SANJANA
│   │
│   ├── models/
│   │   ├── user.go                           DHARMARAJ ✅
│   │   ├── student.go                        DHARMARAJ
│   │   ├── parent.go                         DHARMARAJ
│   │   ├── fee.go                            DHARMARAJ
│   │   ├── teacher.go                        SANJANA
│   │   ├── class.go                          SANJANA
│   │   ├── attendance.go                     SANJANA
│   │   ├── grade.go                          SANJANA
│   │   ├── timetable.go                      SANJANA
│   │   └── announcement.go                   SANJANA
│   │
│   ├── repository/
│   │   ├── user_repository.go                DHARMARAJ ✅
│   │   ├── student_repository.go             DHARMARAJ
│   │   ├── parent_repository.go              DHARMARAJ
│   │   ├── fee_repository.go                 DHARMARAJ
│   │   ├── teacher_repository.go             SANJANA
│   │   ├── class_repository.go               SANJANA
│   │   ├── attendance_repository.go          SANJANA
│   │   ├── grade_repository.go               SANJANA
│   │   ├── timetable_repository.go           SANJANA
│   │   └── announcement_repository.go        SANJANA
│   │
│   ├── services/
│   │   ├── auth_service.go                   DHARMARAJ ✅
│   │   ├── user_service.go                   DHARMARAJ ✅
│   │   ├── student_service.go                DHARMARAJ
│   │   ├── parent_service.go                 DHARMARAJ
│   │   ├── fee_service.go                    DHARMARAJ
│   │   ├── teacher_service.go                SANJANA
│   │   ├── class_service.go                  SANJANA
│   │   ├── attendance_service.go             SANJANA
│   │   ├── grade_service.go                  SANJANA
│   │   ├── timetable_service.go              SANJANA
│   │   └── announcement_service.go           SANJANA
│   │
│   └── handlers/
│       ├── auth_handler.go                   DHARMARAJ ✅
│       ├── user_handler.go                   DHARMARAJ ✅
│       ├── student_handler.go                DHARMARAJ
│       ├── parent_handler.go                 DHARMARAJ
│       ├── fee_handler.go                    DHARMARAJ
│       ├── teacher_handler.go                SANJANA
│       ├── class_handler.go                  SANJANA
│       ├── attendance_handler.go             SANJANA
│       ├── grade_handler.go                  SANJANA
│       ├── timetable_handler.go              SANJANA
│       └── announcement_handler.go           SANJANA
│
├── migrations/
│   ├── 001_create_users.sql                  DHARMARAJ ✅
│   ├── 002_create_students.sql               DHARMARAJ
│   ├── 003_create_parents.sql                DHARMARAJ
│   ├── 004_create_fees.sql                   DHARMARAJ
│   ├── 005_create_teachers.sql               SANJANA
│   ├── 006_create_classes.sql                SANJANA
│   ├── 007_create_attendance.sql             SANJANA
│   ├── 008_create_grades.sql                 SANJANA
│   ├── 009_create_timetable.sql              SANJANA
│   └── 010_create_announcements.sql          SANJANA
│
└── pkg/                                      SHARED — do not modify
    ├── hash/
    ├── jwt/
    └── response/
```

---

## API Endpoints

### Dharmaraj — `/api/v1/auth`, `/api/v1/students`, `/api/v1/parents`, `/api/v1/fees`

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/auth/register` | Public | Register user |
| POST | `/api/v1/auth/login` | Public | Login, returns JWT |
| GET | `/api/v1/students` | Bearer | List students (paginated) |
| POST | `/api/v1/students` | Admin | Enroll new student |
| GET | `/api/v1/students/{id}` | Bearer | Get student by ID |
| PUT | `/api/v1/students/{id}` | Admin | Update student |
| DELETE | `/api/v1/students/{id}` | Admin | Deactivate student |
| GET | `/api/v1/parents` | Bearer | List parents |
| POST | `/api/v1/parents` | Admin | Add parent |
| GET | `/api/v1/parents/{id}` | Bearer | Get parent + linked students |
| GET | `/api/v1/fees` | Admin | List fee structures |
| POST | `/api/v1/fees` | Admin | Create fee structure |
| POST | `/api/v1/fees/payments` | Admin | Record payment |
| GET | `/api/v1/fees/students/{id}/dues` | Bearer | Outstanding dues for student |

---

### Sanjana — `/api/v1/teachers`, `/api/v1/classes`, `/api/v1/attendance`, `/api/v1/grades`, `/api/v1/timetable`, `/api/v1/announcements`

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/teachers` | Bearer | List teachers |
| POST | `/api/v1/teachers` | Admin | Add teacher |
| GET | `/api/v1/teachers/{id}` | Bearer | Get teacher |
| PUT | `/api/v1/teachers/{id}` | Admin | Update teacher |
| GET | `/api/v1/classes` | Bearer | List classes & sections |
| POST | `/api/v1/classes` | Admin | Create class |
| POST | `/api/v1/classes/{id}/students` | Admin | Assign student to class |
| GET | `/api/v1/attendance` | Bearer | Get attendance (filter by class/date) |
| POST | `/api/v1/attendance` | Teacher | Mark attendance |
| PUT | `/api/v1/attendance/{id}` | Teacher | Edit attendance record |
| GET | `/api/v1/grades` | Bearer | List grades/marks |
| POST | `/api/v1/grades` | Teacher | Enter marks |
| GET | `/api/v1/timetable/{class_id}` | Bearer | Get timetable for class |
| POST | `/api/v1/timetable` | Admin | Create timetable entry |
| GET | `/api/v1/announcements` | Bearer | List announcements |
| POST | `/api/v1/announcements` | Admin/Teacher | Post announcement |

---

## Routes Architecture (Conflict Prevention)

`routes/routes.go` **only calls sub-router functions** — never add routes directly to it:

```go
// routes/routes.go  — SHARED, minimal, never conflicts
func Setup(db *gorm.DB, cfg *config.Config, log *zap.Logger) http.Handler {
    r := chi.NewRouter()
    // ... global middleware ...

    r.Route("/api/v1", func(r chi.Router) {
        RegisterDharmarajRoutes(r, db, cfg, log)   // → dharmaraj_routes.go
        RegisterSanjanaRoutes(r, db, cfg, log)     // → sanjana_routes.go
    })
    return r
}
```

Each developer owns their own routes file and registers routes inside it.  
This means `routes.go` itself **will never need editing** once this structure is in place.

---

## Database AutoMigrate (Conflict Prevention)

`database/postgres.go` calls `AutoMigrate`. To avoid both people editing it at the same time:

- **Dharmaraj** adds his models in one block (lines 40–45)
- **Sanjana** adds her models in a separate block (lines 46–51)

```go
// Dharmaraj's models
if err := db.AutoMigrate(
    &models.User{},
    &models.Student{},
    &models.Parent{},
    &models.Fee{},
); err != nil { ... }

// Sanjana's models
if err := db.AutoMigrate(
    &models.Teacher{},
    &models.Class{},
    &models.Attendance{},
    &models.Grade{},
    &models.Timetable{},
    &models.Announcement{},
); err != nil { ... }
```

---

## Git Workflow

### Branch Naming
```
dharmaraj/<feature>    e.g.  dharmaraj/student-module
sanjana/<feature>      e.g.  sanjana/attendance-module
```

### Step-by-Step Flow
```bash
# 1. Always pull main before starting new work
git checkout main && git pull origin main

# 2. Create your feature branch
git checkout -b dharmaraj/student-module   # or sanjana/...

# 3. Work only on YOUR files (see ownership map above)

# 4. Commit with clear messages
git add backend/internal/handlers/student_handler.go
git add backend/internal/services/student_service.go
git commit -m "feat: add student enrollment handler and service"

# 5. Push your branch
git push -u origin dharmaraj/student-module

# 6. Open PR → main, request review from the other person

# 7. After merge, delete the feature branch
git branch -d dharmaraj/student-module
```

### Rules
1. **Never push directly to `main`** — always use a PR
2. **Never edit the other person's files** without a discussion first
3. **Coordinate before touching shared files:** `routes.go`, `postgres.go`, `go.mod`
4. **One module per PR** — keep PRs small and focused
5. **Always pull latest `main`** before starting a new branch to avoid merge conflicts

---

## Local Setup

```bash
# 1. Clone
git clone https://github.com/DharmarajSoundatte/Golang.git
cd Golang/backend

# 2. Create env file
cp .env.example .env
# Edit .env — fill in DATABASE_URL with the Aiven connection string

# 3. Download dependencies
go mod download

# 4. Run
go run ./cmd/server/main.go
```

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `APP_ENV` | `development` or `production` | `development` |
| `DATABASE_URL` | Full Aiven PostgreSQL URL | `postgres://avnadmin:pass@host:port/db?sslmode=require` |
| `SSL_ROOT_CERT` | Path to CA cert (optional) | `ca.pem` |
| `JWT_SECRET` | Secret for signing tokens | `change-in-production` |
| `JWT_EXPIRE_HOURS` | Token expiry in hours | `24` |
| `ALLOWED_ORIGINS` | CORS origins | `http://localhost:5173` |

---

## Project Status

| Module | Owner | Status |
|--------|-------|--------|
| Auth & User Management | Dharmaraj | ✅ Done |
| Student Management | Dharmaraj | 🔲 Pending |
| Parent Management | Dharmaraj | 🔲 Pending |
| Fee Management | Dharmaraj | 🔲 Pending |
| Teacher Management | Sanjana | 🔲 Pending |
| Class & Section Management | Sanjana | 🔲 Pending |
| Attendance Management | Sanjana | 🔲 Pending |
| Grades & Marks | Sanjana | 🔲 Pending |
| Timetable | Sanjana | 🔲 Pending |
| Announcements | Sanjana | 🔲 Pending |
