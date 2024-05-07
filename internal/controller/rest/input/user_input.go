package input

type UserInput struct {
	Name  string `json:"name" binding:"required"  example:"user name"`
	Phone string `json:"phone" binding:"required"  example:"+551199999999"`
}
