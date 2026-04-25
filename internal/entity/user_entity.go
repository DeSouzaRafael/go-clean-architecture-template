package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	ID        uuid.UUID
	Name      string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
