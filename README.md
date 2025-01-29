# Cali: A prototype

> I am going to write a simple appointment schduling application. Target
> demographics will be people familiar with the command line. It also has a web
> interface, that allows to display dates and available time slots, etc.

> We will write this in Go, we will write a small server that keep track of the
> currently scheduled appointments for a user in an sqlite database. A separate
> component, a cli tool will be able to talk to the server via an HTTP API or
> directly to the sqlite3 database for speed.

> let's start with a basic scaffold

```
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
