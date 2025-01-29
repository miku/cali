# Appointment Scheduling App in Go...
Exported on 29/01/2025 at 15:58:08 [from Claude Chat](https://claude.ai/chat/0443ad81-1e51-4913-b717-5e4d2d5acef7) - with [SaveMyChatbot](https://save.hugocollin.com)

## User
I am going to write a simple appointment schduling application. Target demographics will be people familiar with the command line. It also has a web interface, that allows to display dates and available time slots, etc. We will write this in Go, we will write a small server that keep track of the currently scheduled appointments for a user in an sqlite database. A separate component, a cli tool will be able to talk to the server via an HTTP API or directly to the sqlite3 database for speed. let's start with a basic scaffold

## Claude
I'll help you create a scaffold for your appointment scheduling application. Let's break this down into components and create a basic structure.

---
**Appointment Scheduler Project Structure:**
```plaintext
appointment-scheduler/
├── cmd/
│   ├── server/
│   │   └── main.go         # Server entry point
│   └── cli/
│       └── main.go         # CLI tool entry point
├── internal/
│   ├── api/
│   │   ├── handlers.go     # HTTP handlers
│   │   └── routes.go       # Route definitions
│   ├── db/
│   │   ├── schema.sql      # Database schema
│   │   └── db.go          # Database operations
│   ├── models/
│   │   └── appointment.go  # Data models
│   └── config/
│       └── config.go       # Configuration handling
├── web/
│   ├── templates/
│   │   └── index.html     # Web interface templates
│   └── static/
│       ├── css/
│       └── js/
├── go.mod
└── README.md
```
---

Let's create the basic database schema first:

---
**Database Schema:**
```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS appointments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    CHECK (end_time > start_time)
);
```
---

Now let's create the basic Go module and the appointment model:

---
**Go Module Definition:**
```go
module github.com/yourusername/appointment-scheduler

go 1.21

require (
    github.com/mattn/go-sqlite3 v1.14.22
    github.com/gorilla/mux v1.8.1
    github.com/spf13/cobra v1.8.0
)
```
---

---
**Appointment Model:**
```go
package models

import (
    "time"
)

type User struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    CreatedAt time.Time `json:"created_at"`
}

type Appointment struct {
    ID          int64     `json:"id"`
    UserID      int64     `json:"user_id"`
    Title       string    `json:"title"`
    Description string    `json:"description,omitempty"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Validate checks if the appointment data is valid
func (a *Appointment) Validate() error {
    if a.Title == "" {
        return ErrEmptyTitle
    }
    if a.StartTime.IsZero() || a.EndTime.IsZero() {
        return ErrInvalidTime
    }
    if a.EndTime.Before(a.StartTime) {
        return ErrEndTimeBeforeStart
    }
    return nil
}

// Custom errors for appointment validation
var (
    ErrEmptyTitle        = errors.New("title cannot be empty")
    ErrInvalidTime       = errors.New("invalid time")
    ErrEndTimeBeforeStart = errors.New("end time cannot be before start time")
)
```
---

This gives us a good foundation to build upon. Let me explain the key components:

1.  Project Structure:
    *   `cmd/` contains the entry points for both the server and CLI
    *   `internal/` contains the private application code
    *   `web/` contains the web interface assets
2.  Database Schema:
    *   Two tables: `users` and `appointments`
    *   Proper foreign key constraints and timestamp handling
    *   Check constraint to ensure end\_time is after start\_time
3.  Go Module:
    *   Essential dependencies included
    *   Using popular libraries like gorilla/mux for routing and cobra for CLI
4.  Appointment Model:
    *   Structured with JSON tags for API responses
    *   Basic validation logic
    *   Clear error types

Would you like to focus on implementing the server or the CLI tool first? Also, I can help you set up the database operations next if you prefer to start there.

