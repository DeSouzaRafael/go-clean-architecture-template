package usecase

import (
	"context"

	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
)

type (
	// Routers
	User interface {
		CreateUser(context.Context, entity.UserEntity) (entity.UserEntity, error)
		UpdateUser(context.Context, entity.UserEntity) error
		DeleteUser(context.Context, entity.UserEntity) error
		GetUserById(context.Context, entity.UserEntity) (entity.UserEntity, error)
	}

	// Repository
	UserRepo interface {
		Create(context.Context, entity.UserEntity) (entity.UserEntity, error)
		Update(context.Context, entity.UserEntity) error
		DeleteById(context.Context, entity.UserEntity) error
		GetById(context.Context, entity.UserEntity) (entity.UserEntity, error)
	}

	// Add all use cases for use in NewRouter
	UseCases interface {
		UserUseCase() User
	}
)

// Adjust the items below as new use cases are created
type AppUseCases struct {
	user User
}

func (a *AppUseCases) UserUseCase() User {
	return a.user
}

func NewAppUseCases(user User) *AppUseCases {
	return &AppUseCases{
		user: user,
	}
}
