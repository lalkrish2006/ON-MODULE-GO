# OD System - Go Backend

A complete Go backend migration for the On-Duty Application Module, replacing the legacy PHP implementation.

## 🚀 Features Implemented
- **Authentication**: Role-based login (Student, Mentor, HOD, Principal, LabTech, Admin) with Session Management.
- **Dashboards**:
  - **Student**: Apply for OD, View History, track status.
  - **Mentor**: Approve/Reject Student ODs.
  - **HOD**: Approve/Reject Mentor-accepted ODs.
  - **Principal**: Approve/Reject External ODs.
  - **LabTech**: Approve Lab access for Internal ODs.
- **API**: Endpoints for fetching student details and mentors via AJAX.
- **Email Notifications**: Stub implementation (logs to console).

## 🛠 Prerequisites
- **Go 1.21+** installed.
- **MySQL Database** running on `localhost:3307` (or update `internal/config/config.go`).
- Database schema (`college_db`) must already exist (re-uses existing PHP DB).

## 📂 Project Structure
```
od_go/
├── cmd/
│   └── server/
│       └── main.go       # Entry point
├── internal/
│   ├── config/           # Configuration
│   ├── database/         # DB Connection
│   ├── handlers/         # HTTP Handlers (Logic)
│   ├── middleware/       # Auth Middleware
│   ├── models/           # DB Structs
│   ├── services/         # Session & Email Services
│   └── utils/            # Hashing & Helpers
├── templates/            # HTML Templates (Go-enabled)
└── go.mod                # Go Module definition
```

## ▶️ How to Run
1. Open a terminal in the `od_go` directory.
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```
4. Access the application at:
   [http://localhost:8080](http://localhost:8080)

## 🔑 Default Credentials (from DB)
- **Student**: Use your Register Number & Password.
- **Staff (Mentor/HOD)**: Use your Name/Email & Password.
- **Admin**: Standard admin credentials.

## 📝 Notes
- **Static files**: Templates use CDNs for Bootstrap and Icons. No local CSS/JS files are required for the MVP.
- **Email**: Emails are currently logged to the console terminal instead of being sent.
