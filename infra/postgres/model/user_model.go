package model

import (
	"time"

	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string
	Phone     string
	CreatedAt time.Time  `gorm:"<-:create"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
}

func (UserModel) TableName() string { return "user" }

func ToUserModel(e entity.UserEntity) UserModel {
	return UserModel{
		ID:        e.ID,
		Name:      e.Name,
		Phone:     e.Phone,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}

func ToUserEntity(m UserModel) entity.UserEntity {
	return entity.UserEntity{
		ID:        m.ID,
		Name:      m.Name,
		Phone:     m.Phone,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}
