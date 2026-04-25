package repository

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepo[T any] struct {
	DB *gorm.DB
}

type PageOptions struct {
	Offset int
	Limit  int
}

func NewBaseRepo[T any](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{DB: db}
}

func (repo *BaseRepo[T]) Get(ctx context.Context, id interface{}) (T, error) {
	var entity T
	err := repo.DB.WithContext(ctx).First(&entity, id).Error
	return entity, err
}

func (repo *BaseRepo[T]) List(ctx context.Context, opts PageOptions) ([]T, error) {
	var items []T
	err := repo.DB.WithContext(ctx).Offset(opts.Offset).Limit(opts.Limit).Find(&items).Error
	return items, err
}

func (repo *BaseRepo[T]) Create(ctx context.Context, entity T) (T, error) {
	err := repo.DB.WithContext(ctx).Create(&entity).Error
	return entity, err
}

func (repo *BaseRepo[T]) Update(ctx context.Context, entity T) error {
	return repo.DB.WithContext(ctx).Save(&entity).Error
}

func (repo *BaseRepo[T]) Delete(ctx context.Context, id interface{}) error {
	var entity T
	return repo.DB.WithContext(ctx).Delete(&entity, id).Error
}
