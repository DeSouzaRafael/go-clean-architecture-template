package usecase

import (
	"context"

	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"

	"fmt"
)

type UserUseCase struct {
	repo UserRepo
}

func NewUser(c UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: c,
	}
}

func (uc *UserUseCase) GetUserById(ctx context.Context, user entity.UserEntity) (entity.UserEntity, error) {

	user, err := uc.repo.GetById(ctx, user)
	if err != nil {
		return entity.UserEntity{}, fmt.Errorf("GetUserById: %w", err)
	}

	return user, nil
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user entity.UserEntity) (entity.UserEntity, error) {

	user, err := uc.repo.Create(ctx, user)
	if err != nil {
		return entity.UserEntity{}, fmt.Errorf("CreateUser: %w", err)
	}

	return user, nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user entity.UserEntity) error {

	_, err := uc.repo.GetById(ctx, user)
	if err != nil {
		return fmt.Errorf("UpdateUser: %w", err)
	}

	err = uc.repo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("UpdateUser: %w", err)
	}

	return nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, user entity.UserEntity) error {

	_, err := uc.repo.GetById(ctx, user)
	if err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}

	err = uc.repo.DeleteById(ctx, user)
	if err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}

	return nil
}
