package v0

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/logger"
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/validator"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/controller/rest/input"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/controller/rest/output"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/usecase"
)

type userRoutes struct {
	usecase   usecase.User
	logger    logger.Interface
	validator *validator.Validator
}

func NewUserRoutes(handler *echo.Group, l logger.Interface, v *validator.Validator, uc usecase.User) {

	ur := &userRoutes{uc, l, v}

	g := handler.Group("/user")
	g.GET("/:id", ur.get)
	g.POST("", ur.create)
	g.PUT("/:id", ur.update)
	g.DELETE("/:id", ur.delete)
}

// @Summary     Get User
// @Description Search for a user by its UUID.
// @ID          getUser
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       id   path   string  true  "User ID"
// @Success     200  {object} output.UserOutput  "Returns the found user"
// @Failure     400  {object} output.ResponseError  "Invalid UUID format"
// @Failure     500  {object} output.ResponseError  "Internal server error"
// @Router      /v0/user/{id} [get]
func (ur *userRoutes) get(c echo.Context) error {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ur.logger.Error(err, "http - v0 - get")
		return output.ErrorResponse(c, http.StatusBadRequest, "invalid UUID format")
	}

	user, err := ur.usecase.GetUserById(c.Request().Context(), entity.UserEntity{ID: id})
	if err != nil {
		ur.logger.Error(err, "http - v0 - update")
		return output.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	response := output.UserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Phone: user.Phone,
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary     Update User
// @Description update existing user details
// @ID          updateUser
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path string true "User ID"
// @Param       request body input.UserInput true "Update user details"
// @Success     200 "User Successfully updated"
// @Failure     400 {object} output.ResponseError
// @Failure     404 {object} output.ResponseError "User not found"
// @Failure     500 {object} output.ResponseError
// @Router      /v0/user/{id} [put]
func (ur *userRoutes) update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ur.logger.Error(err, "http - v0 - update")
		return output.ErrorResponse(c, http.StatusBadRequest, "invalid UUID format")
	}

	var input input.UserInput
	if err := c.Bind(&input); err != nil {
		ur.logger.Error(err, "http - v0 - update")
		return output.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := input.Validate(ur.validator); err != nil {
		ur.logger.Error(err, "http - v0 - update")
		return output.ErrorResponse(c, http.StatusBadRequest, "invalid request data: "+err.Error())
	}

	err = ur.usecase.UpdateUser(
		c.Request().Context(),
		entity.UserEntity{
			ID:    id,
			Name:  input.Name,
			Phone: input.Phone,
		},
	)
	if err != nil {
		ur.logger.Error(err, "http - v0 - update")
		return output.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully updated"})
}

// @Summary     Create User
// @Description add new user
// @ID          create
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       user body input.UserInput true "Set up users"
// @Success     200 {object} output.UserOutput
// @Failure     400 {object} output.ResponseError
// @Failure     500 {object} output.ResponseError
// @Router      /v0/user [post]
func (ur *userRoutes) create(c echo.Context) error {

	var input input.UserInput
	if err := c.Bind(&input); err != nil {
		ur.logger.Error(err, "http - v0 - create")
		return output.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := input.Validate(ur.validator); err != nil {
		ur.logger.Error(err, "http - v0 - update validation")
		return output.ErrorResponse(c, http.StatusBadRequest, "invalid request data: "+err.Error())
	}

	user, err := ur.usecase.CreateUser(
		c.Request().Context(),
		entity.UserEntity{
			Name:  input.Name,
			Phone: input.Phone,
		},
	)
	if err != nil {
		ur.logger.Error(err, "http - v0 - create")
		return output.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	response := output.UserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Phone: user.Phone,
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary     Delete User
// @Description Deletes a user by its UUID.
// @ID          deleteUser
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id   path   string  true  "User ID"
// @Success     200  "User successfully deleted"
// @Failure     400  {object} output.ResponseError  "Invalid UUID format"
// @Failure     500  {object} output.ResponseError  "Internal server error"
// @Router      /v0/user/{id} [delete]
func (ur *userRoutes) delete(c echo.Context) error {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ur.logger.Error(err, "http - v0 - delete")
		return output.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID format")
	}

	err = ur.usecase.DeleteUser(c.Request().Context(), entity.UserEntity{ID: id})
	if err != nil {
		ur.logger.Error(err, "http - v0 - delete")
		return output.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully deleted"})
}
