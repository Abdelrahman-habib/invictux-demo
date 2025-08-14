package database

import "time"

// AppSetting represents application configuration settings
type AppSetting struct {
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// Common database model interfaces
type Model interface {
	GetID() string
}

// Timestamps for database models
type Timestamps struct {
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
