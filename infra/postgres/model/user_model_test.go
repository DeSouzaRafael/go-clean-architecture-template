package model

import (
	"testing"
	"time"

	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestToUserModel(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	id := uuid.New()

	e := entity.UserEntity{
		ID:        id,
		Name:      "Alice",
		Phone:     "+5511999999999",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	m := ToUserModel(e)

	assert.Equal(t, id, m.ID)
	assert.Equal(t, "Alice", m.Name)
	assert.Equal(t, "+5511999999999", m.Phone)
	assert.Equal(t, now, m.CreatedAt)
	assert.Equal(t, now, m.UpdatedAt)
	assert.False(t, m.DeletedAt.Valid)
}

func TestToUserModel_WithDeletedAt(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	e := entity.UserEntity{DeletedAt: &now}

	m := ToUserModel(e)

	assert.True(t, m.DeletedAt.Valid)
	assert.Equal(t, now, m.DeletedAt.Time.UTC().Truncate(time.Second))
}

func TestToUserEntity(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	id := uuid.New()

	m := UserModel{
		ID:        id,
		Name:      "Bob",
		Phone:     "+5511888888888",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: gorm.DeletedAt{Valid: false},
	}

	e := ToUserEntity(m)

	assert.Equal(t, id, e.ID)
	assert.Equal(t, "Bob", e.Name)
	assert.Equal(t, "+5511888888888", e.Phone)
	assert.Equal(t, now, e.CreatedAt)
	assert.Equal(t, now, e.UpdatedAt)
	assert.Nil(t, e.DeletedAt)
}

func TestToUserEntity_WithDeletedAt(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	m := UserModel{
		DeletedAt: gorm.DeletedAt{Time: now, Valid: true},
	}

	e := ToUserEntity(m)

	require.NotNil(t, e.DeletedAt)
	assert.Equal(t, now, e.DeletedAt.UTC().Truncate(time.Second))
}

func TestRoundTrip(t *testing.T) {
	id := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)

	original := entity.UserEntity{
		ID:        id,
		Name:      "Round Trip",
		Phone:     "+5511777777777",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ToUserEntity(ToUserModel(original))

	assert.Equal(t, original.ID, result.ID)
	assert.Equal(t, original.Name, result.Name)
	assert.Equal(t, original.Phone, result.Phone)
	assert.Nil(t, result.DeletedAt)
}

func TestTableName(t *testing.T) {
	m := UserModel{}
	assert.Equal(t, "user", m.TableName())
}
