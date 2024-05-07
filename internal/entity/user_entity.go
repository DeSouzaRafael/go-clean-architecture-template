package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	CreatedAt time.Time  `json:"created_at" gorm:"<-:create"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

func (c *UserEntity) TableName() string {
	return "user"
}
