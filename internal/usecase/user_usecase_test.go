// user_usecase_test.go
package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/usecase"
	"github.com/DeSouzaRafael/go-clean-architecture-template/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	uc := usecase.NewUser(mockRepo)

	user := entity.UserEntity{
		ID:   uuid.New(),
		Name: "User Name",
	}

	mockRepo.EXPECT().Create(gomock.Any(), user).Return(user, nil)

	createdUser, err := uc.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, user, createdUser)

	expectedErr := fmt.Errorf("some error")
	mockRepo.EXPECT().Create(gomock.Any(), user).Return(entity.UserEntity{}, expectedErr)

	createdUserError, err := uc.CreateUser(context.Background(), user)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("CreateUser: %v", expectedErr))
	assert.Equal(t, entity.UserEntity{}, createdUserError)
}

func TestGetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	uc := usecase.NewUser(mockRepo)

	user := entity.UserEntity{
		ID:   uuid.New(),
		Name: "User Name",
	}

	mockRepo.EXPECT().GetById(gomock.Any(), user).Return(user, nil)

	fetchedUser, err := uc.GetUserById(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, user, fetchedUser)

	expectedErr := fmt.Errorf("user not found")
	mockRepo.EXPECT().GetById(gomock.Any(), user).Return(entity.UserEntity{}, expectedErr)

	fetchedUserError, err := uc.GetUserById(context.Background(), user)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("GetUserById: %v", expectedErr))
	assert.Equal(t, entity.UserEntity{}, fetchedUserError)
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	uc := usecase.NewUser(mockRepo)

	user := entity.UserEntity{
		ID:   uuid.New(),
		Name: "User Name",
	}

	gomock.InOrder(
		mockRepo.EXPECT().GetById(gomock.Any(), user).Return(user, nil),
		mockRepo.EXPECT().Update(gomock.Any(), user).Return(nil),
	)

	err := uc.UpdateUser(context.Background(), user)

	assert.NoError(t, err)

	expectedErr := fmt.Errorf("user not found")
	mockRepo.EXPECT().GetById(gomock.Any(), user).Return(entity.UserEntity{}, expectedErr)

	err = uc.UpdateUser(context.Background(), user)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("UpdateUser: %v", expectedErr))

	gomock.InOrder(
		mockRepo.EXPECT().GetById(gomock.Any(), user).Return(user, nil),
		mockRepo.EXPECT().Update(gomock.Any(), user).Return(fmt.Errorf("failed to update")),
	)

	err = uc.UpdateUser(context.Background(), user)

	assert.Error(t, err)
	assert.EqualError(t, err, "UpdateUser: failed to update")
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	uc := usecase.NewUser(mockRepo)

	user := entity.UserEntity{
		ID:   uuid.New(),
		Name: "User Name",
	}

	gomock.InOrder(
		mockRepo.EXPECT().GetById(gomock.Any(), user).Return(user, nil),
		mockRepo.EXPECT().DeleteById(gomock.Any(), user).Return(nil),
	)

	err := uc.DeleteUser(context.Background(), user)

	assert.NoError(t, err)

	expectedErr := fmt.Errorf("user not found")
	mockRepo.EXPECT().GetById(gomock.Any(), user).Return(entity.UserEntity{}, expectedErr)

	err = uc.DeleteUser(context.Background(), user)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("DeleteUser: %v", expectedErr))

	gomock.InOrder(
		mockRepo.EXPECT().GetById(gomock.Any(), user).Return(user, nil),
		mockRepo.EXPECT().DeleteById(gomock.Any(), user).Return(fmt.Errorf("failed to delete")),
	)

	err = uc.DeleteUser(context.Background(), user)

	assert.Error(t, err)
	assert.EqualError(t, err, "DeleteUser: failed to delete")
}
