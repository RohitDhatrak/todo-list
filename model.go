package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string
	Email     string
	DeletedAt *time.Time
}

type Todo struct {
	gorm.Model
	Title       string
	UserId      uint
	CompletedAt *time.Time
	DueDate     *time.Time
	Description string
	Status      TodoStatus
}

type TodoStatus string

const (
	PENDING     TodoStatus = "Pending"
	IN_PROGRESS TodoStatus = "In Progress"
	COMPLETED   TodoStatus = "Completed"
)
