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
	TaskIDs        []int     `json:"taskIds"`
}
