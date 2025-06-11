package entity

import (
	"time"
)

// User is a struct that represents a user in database.
type User struct {
	ID        int64     `json:"id" gorm:"primaryKey;not null"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	// TODO: add token invalidation
	// AuthUUID string `json:"auth_uuid" gorm:"not null;type:uuid;default:gen_random_uuid()"`

	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	PhotoUrl  string `json:"photo_url"`
}
