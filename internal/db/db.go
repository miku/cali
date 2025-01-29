package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/miku/cali/internal/models"
)

type Database struct {
	db *sql.DB
}

func New(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

// InitSchema creates the database schema if it doesn't exist
func (d *Database) InitSchema() error {
	schema := `
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
        );`

	_, err := d.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// CreateAppointment inserts a new appointment into the database
func (d *Database) CreateAppointment(a *models.Appointment) error {
	query := `
        INSERT INTO appointments (
            user_id, title, description, start_time, end_time
        ) VALUES (?, ?, ?, ?, ?)
        RETURNING id, created_at, updated_at`

	err := d.db.QueryRow(
		query,
		a.UserID,
		a.Title,
		a.Description,
		a.StartTime,
		a.EndTime,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}

	return nil
}

// GetAppointment retrieves an appointment by ID
func (d *Database) GetAppointment(id int64) (*models.Appointment, error) {
	a := &models.Appointment{}
	query := `
        SELECT id, user_id, title, description, start_time, end_time,
               created_at, updated_at
        FROM appointments
        WHERE id = ?`

	err := d.db.QueryRow(query, id).Scan(
		&a.ID,
		&a.UserID,
		&a.Title,
		&a.Description,
		&a.StartTime,
		&a.EndTime,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return a, nil
}

// ListAppointments retrieves appointments for a user within a time range
func (d *Database) ListAppointments(userID int64, start, end time.Time) ([]*models.Appointment, error) {
	query := `
        SELECT id, user_id, title, description, start_time, end_time,
               created_at, updated_at
        FROM appointments
        WHERE user_id = ?
        AND start_time >= ?
        AND end_time <= ?
        ORDER BY start_time ASC`

	rows, err := d.db.Query(query, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to list appointments: %w", err)
	}
	defer rows.Close()

	var appointments []*models.Appointment
	for rows.Next() {
		a := &models.Appointment{}
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Title,
			&a.Description,
			&a.StartTime,
			&a.EndTime,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}
		appointments = append(appointments, a)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating appointments: %w", err)
	}

	return appointments, nil
}

// UpdateAppointment updates an existing appointment
func (d *Database) UpdateAppointment(a *models.Appointment) error {
	query := `
        UPDATE appointments
        SET title = ?, description = ?, start_time = ?, end_time = ?,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = ? AND user_id = ?
        RETURNING updated_at`

	err := d.db.QueryRow(
		query,
		a.Title,
		a.Description,
		a.StartTime,
		a.EndTime,
		a.ID,
		a.UserID,
	).Scan(&a.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	return nil
}

// DeleteAppointment removes an appointment
func (d *Database) DeleteAppointment(id, userID int64) error {
	query := `DELETE FROM appointments WHERE id = ? AND user_id = ?`

	result, err := d.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("appointment not found or unauthorized")
	}

	return nil
}
