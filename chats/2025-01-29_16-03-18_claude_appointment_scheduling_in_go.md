# Appointment Scheduling App in Go...
Exported on 29/01/2025 at 16:03:18 [from Claude Chat](https://claude.ai/chat/0443ad81-1e51-4913-b717-5e4d2d5acef7) - with [SaveMyChatbot](https://save.hugocollin.com)

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

## User
let's call the project "cali" and the server will be called calid

Let's scaffold the server component.

## Claude
I'll help you scaffold the server component "calid". Let's update our structure and create the server implementation.

---
**Appointment Scheduler Project Structure:**
```plaintext
cali/
├── cmd/
│   ├── calid/
│   │   └── main.go         # Server daemon entry point
│   └── cali/
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

---
**Go Module Definition:**
```go
module github.com/yourusername/cali

go 1.21

require (
    github.com/mattn/go-sqlite3 v1.14.22
    github.com/gorilla/mux v1.8.1
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.2
)
```
---

Now let's create the server implementation starting with the configuration and main server code:

---
**Server Configuration:**
```go
package config

import (
    "github.com/spf13/viper"
    "path/filepath"
)

type Config struct {
    Server struct {
        Host string
        Port int
    }
    Database struct {
        Path string
    }
    Web struct {
        TemplatesDir string
        StaticDir    string
    }
}

func LoadConfig() (*Config, error) {
    viper.SetDefault("server.host", "127.0.0.1")
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("database.path", "./cali.db")
    viper.SetDefault("web.templatesdir", "./web/templates")
    viper.SetDefault("web.staticdir", "./web/static")

   // Look for config in standard locations
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("/etc/cali/")
    viper.AddConfigPath("$HOME/.config/cali/")
    viper.AddConfigPath(".")

   var config Config

   // Read config file if it exists
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }

   if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

   // Ensure database path is absolute
    if !filepath.IsAbs(config.Database.Path) {
        absPath, err := filepath.Abs(config.Database.Path)
        if err != nil {
            return nil, err
        }
        config.Database.Path = absPath
    }

   return &config, nil
}
```
---

---
**Server Main Implementation:**
```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

   "github.com/yourusername/cali/internal/api"
    "github.com/yourusername/cali/internal/config"
    "github.com/yourusername/cali/internal/db"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

   // Initialize database
    database, err := db.New(cfg.Database.Path)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer database.Close()

   // Initialize API server
    server := api.NewServer(database, cfg)

   // Create HTTP server
    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    srv := &http.Server{
        Addr:         addr,
        Handler:      server.Router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

   // Start server in a goroutine
    go func() {
        log.Printf("Starting server on %s", addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

   // Set up graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

   if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

   log.Println("Server exited properly")
}
```
---

Let's also create the basic API server structure:

---
**API Server Implementation:**
```go
package api

import (
    "encoding/json"
    "net/http"

   "github.com/gorilla/mux"
    "github.com/yourusername/cali/internal/config"
    "github.com/yourusername/cali/internal/db"
)

type Server struct {
    Router *mux.Router
    db     *db.Database
    config *config.Config
}

func NewServer(db *db.Database, cfg *config.Config) *Server {
    s := &Server{
        Router: mux.NewRouter(),
        db:     db,
        config: cfg,
    }
    s.routes()
    return s
}

func (s *Server) routes() {
    // API routes
    api := s.Router.PathPrefix("/api").Subrouter()
    api.HandleFunc("/appointments", s.handleListAppointments).Methods("GET")
    api.HandleFunc("/appointments", s.handleCreateAppointment).Methods("POST")
    api.HandleFunc("/appointments/{id}", s.handleGetAppointment).Methods("GET")
    api.HandleFunc("/appointments/{id}", s.handleUpdateAppointment).Methods("PUT")
    api.HandleFunc("/appointments/{id}", s.handleDeleteAppointment).Methods("DELETE")

   // Web interface routes
    s.Router.PathPrefix("/static/").Handler(
        http.StripPrefix("/static/", 
        http.FileServer(http.Dir(s.config.Web.StaticDir))))
    s.Router.HandleFunc("/", s.handleIndex).Methods("GET")
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if data != nil {
        if err := json.NewEncoder(w).Encode(data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
}

func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
    s.respondJSON(w, status, map[string]string{"error": message})
}
```
---

This gives us a solid foundation for the server component. Here's what we have:

1.  Configuration (`config.go`):
    *   Uses viper for flexible configuration
    *   Supports YAML config files in standard locations
    *   Has sensible defaults
    *   Handles both server and database configuration
2.  Server Main (`cmd/calid/main.go`):
    *   Graceful shutdown handling
    *   Proper timeout configuration
    *   Clean separation of concerns
3.  API Server (`internal/api/server.go`):
    *   RESTful API structure
    *   Static file serving for web interface
    *   Clean route organization
    *   Helper methods for JSON responses

Next steps could be:

1.  Implementing the database operations in `internal/db/db.go`
2.  Adding the individual handler implementations
3.  Setting up the web interface templates

Which would you like to tackle next?

