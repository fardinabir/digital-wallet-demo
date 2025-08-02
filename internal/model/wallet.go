package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Wallet is the model for the wallet endpoint.
type Wallet struct {
	ID        int `gorm:"primaryKey"`
	Task      string
	Status    Status
	Priority  Priority
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// NewWallet returns a new instance of the wallet model.
func NewWallet(task string, priority Priority) *Wallet {
	return &Wallet{
		Task:     task,
		Priority: priority,
		Status:   Created,
	}
}

// NewUpdateWallet returns a new instance of the wallet model for updating.
func NewUpdateWallet(id int, task string, priority Priority, status Status) *Wallet {
	return &Wallet{
		ID:       id,
		Task:     task,
		Status:   status,
		Priority: priority,
	}
}

// Status is the status of the task.
type Status string

const (
	// Created is the status for a created task.
	Created = Status("created")
	// Processing is the status for a processing task.
	Processing = Status("processing")
	// Done is the status for a done task.
	Done = Status("done")
)

// StatusMap is a map of task status.
var StatusMap = map[Status]bool{
	Created:    true,
	Processing: true,
	Done:       true,
}

// IsValidStatus checks if the status is valid (Created, Processing, Done)
func IsValidStatus(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true // Skip validation for empty or nil fields
	}
	status := fl.Field().Interface().(Status)
	return status == Created || status == Processing || status == Done
}

// Priority is the priority of the task.
type Priority int

const (
	// Low is the lowest priority of task
	Low Priority = iota + 1
	// Medium is the medium priority of task
	Medium
	// High is the highest priority of task
	High
)

// IsValidPriority checks if the priority is valid (High, Medium, Low)
func IsValidPriority(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true // Skip validation for empty or nil fields
	}
	priority := fl.Field().Interface().(Priority)
	return priority == Low || priority == Medium || priority == High
}
