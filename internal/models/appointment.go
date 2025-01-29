package models

import (
	"errors"
	"time"
)

// Custom errors for appointment validation
var (
	ErrEmptyTitle         = errors.New("title cannot be empty")
	ErrInvalidTime        = errors.New("invalid time")
	ErrEndTimeBeforeStart = errors.New("end time cannot be before start time")
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
