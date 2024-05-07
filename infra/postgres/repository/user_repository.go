package repository

import (
	"context"

	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
	"gorm.io/gorm"
)

type UserRepo struct {
	*BaseRepo[entity.UserEntity]
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		BaseRepo: NewBaseRepo[entity.UserEntity](db),
	}
}

func (r *UserRepo) GetById(ctx context.Context, c entity.UserEntity) (entity.UserEntity, error) {
	return r.Get(ctx, c.ID)
}

func (r *UserRepo) DeleteById(ctx context.Context, c entity.UserEntity) error {
	return r.Delete(ctx, c.ID)
}
