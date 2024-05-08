package input

import "github.com/DeSouzaRafael/go-clean-architecture-template/infra/validator"

type UserInput struct {
	Name  string `json:"name" binding:"required" validate:"required" example:"user name"`
	Phone string `json:"phone" binding:"required" validate:"required,phone_format" example:"+5511999999999"`
}

func (input *UserInput) Validate(v *validator.Validator) error {
	return v.Validate(input)
}
