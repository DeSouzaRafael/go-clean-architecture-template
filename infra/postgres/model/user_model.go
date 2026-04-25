package model

import (
	"time"

	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string
	Phone     string
	CreatedAt time.Time      `gorm:"<-:create"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string { return "user" }

func ToUserModel(e entity.UserEntity) UserModel {
	var deletedAt gorm.DeletedAt
	if e.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *e.DeletedAt, Valid: true}
	}
	return UserModel{
		ID:        e.ID,
		Name:      e.Name,
		Phone:     e.Phone,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func ToUserEntity(m UserModel) entity.UserEntity {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		deletedAt = &t
	}
	return entity.UserEntity{
		ID:        m.ID,
		Name:      m.Name,
		Phone:     m.Phone,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
