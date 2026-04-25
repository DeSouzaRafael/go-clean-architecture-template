package input

import (
	"testing"

	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserInput_Validate_Valid(t *testing.T) {
	v := validator.NewValidator()
	input := UserInput{Name: "John Doe", Phone: "+5511999999999"}
	err := input.Validate(v)
	assert.NoError(t, err)
}

func TestUserInput_Validate_MissingName(t *testing.T) {
	v := validator.NewValidator()
	input := UserInput{Name: "", Phone: "+5511999999999"}
	err := input.Validate(v)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Name")
}

func TestUserInput_Validate_MissingPhone(t *testing.T) {
	v := validator.NewValidator()
	input := UserInput{Name: "John", Phone: ""}
	err := input.Validate(v)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Phone")
}

func TestUserInput_Validate_AllMissing(t *testing.T) {
	v := validator.NewValidator()
	input := UserInput{}
	err := input.Validate(v)
	require.Error(t, err)
}
