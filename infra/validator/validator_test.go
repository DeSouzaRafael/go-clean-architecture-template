package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	assert.NotNil(t, v)
}

func TestValidate_Valid(t *testing.T) {
	v := NewValidator()
	input := testStruct{Name: "John", Email: "john@example.com"}
	err := v.Validate(input)
	assert.NoError(t, err)
}

func TestValidate_MissingRequired(t *testing.T) {
	v := NewValidator()
	input := testStruct{Name: "", Email: "john@example.com"}
	err := v.Validate(input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Name")
	assert.Contains(t, err.Error(), "required")
}

func TestValidate_MultipleErrors(t *testing.T) {
	v := NewValidator()
	input := testStruct{Name: "", Email: ""}
	err := v.Validate(input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Name")
	assert.Contains(t, err.Error(), "Email")
}

func TestValidate_InvalidEmail(t *testing.T) {
	v := NewValidator()
	input := testStruct{Name: "John", Email: "not-an-email"}
	err := v.Validate(input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Email")
}

func TestValidate_NotAStruct(t *testing.T) {
	v := NewValidator()
	input := "not a struct"
	err := v.Validate(input)
	require.Error(t, err)
}
