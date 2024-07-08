package models

type Task struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	TotalHours   int    `json:"total_hours"`
	TotalMinutes int    `json:"total_minutes"`
	UserID       int    `json:"user_id"`
}
