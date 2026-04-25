package repository

import (
	"context"

	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/postgres/model"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
	"gorm.io/gorm"
)

type UserRepo struct {
	*BaseRepo[model.UserModel]
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{BaseRepo: NewBaseRepo[model.UserModel](db)}
}

func (r *UserRepo) Create(ctx context.Context, e entity.UserEntity) (entity.UserEntity, error) {
	m, err := r.BaseRepo.Create(ctx, model.ToUserModel(e))
	return model.ToUserEntity(m), err
}

func (r *UserRepo) GetById(ctx context.Context, e entity.UserEntity) (entity.UserEntity, error) {
	m, err := r.BaseRepo.Get(ctx, e.ID)
	return model.ToUserEntity(m), err
}

func (r *UserRepo) Update(ctx context.Context, e entity.UserEntity) error {
	return r.BaseRepo.Update(ctx, model.ToUserModel(e))
}

func (r *UserRepo) DeleteById(ctx context.Context, e entity.UserEntity) error {
	return r.BaseRepo.Delete(ctx, e.ID)
}
