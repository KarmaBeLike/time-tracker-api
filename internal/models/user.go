package models

import (
	"time"
)

// User represents the structure of a user in the system.
type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	PassportNumber string    `json:"passportNumber"`
	Surname        string    `json:"surname"`
	Patronymic     string    `json:"patronymic"`
	Address        string    `json:"address"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// TimeEntry represents the structure for tracking time spent on tasks.
type TimeEntry struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Task      string    `json:"task"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  int       `json:"duration"` // Duration in minutes
}
