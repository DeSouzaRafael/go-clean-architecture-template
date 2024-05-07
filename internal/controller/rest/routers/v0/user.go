package v0

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/logger"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/controller/rest"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/controller/rest/input"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/controller/rest/output"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/usecase"
)

type userRoutes struct {
	usecase usecase.User
	logger  logger.Interface
}

func NewUserRoutes(handler *echo.Group, l logger.Interface, usecase usecase.User) {
	r := &userRoutes{usecase, l}

	g := handler.Group("/user")
	g.GET("/:id", r.get)
	g.POST("", r.create)
	g.PUT("/:id", r.update)
	g.DELETE("/:id", r.delete)
}

// @Summary     Get User
// @Description Search for a user by its UUID.
// @ID          getUser
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       id   path   string  true  "User ID"
// @Success     200  {object} output.UserOutput  "Returns the found user"
// @Failure     400  {object} rest.response  "Invalid UUID format"
// @Failure     500  {object} rest.response  "Internal server error"
// @Router      /v0/user/{id} [get]
func (r *userRoutes) get(c echo.Context) error {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		r.logger.Error(err, "http - v0 - get")
		return rest.ErrorResponse(c, http.StatusBadRequest, "invalid UUID format")
	}

	user, err := r.usecase.GetUserById(c.Request().Context(), entity.UserEntity{ID: id})
	if err != nil {
		r.logger.Error(err, "http - v0 - update")
		return rest.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
// @Param       request body request.UserRequest true "Update user details"
// @Success     200 "User Successfully updated"
// @Failure     400 {object} rest.response
// @Failure     404 {object} rest.response "User not found"
// @Failure     500 {object} rest.response
// @Router      /v0/user/{id} [put]
func (r *userRoutes) update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		r.logger.Error(err, "http - v0 - update")
		return rest.ErrorResponse(c, http.StatusBadRequest, "invalid UUID format")
	}

	var input input.UserInput
	if err := c.Bind(&input); err != nil {
		r.logger.Error(err, "http - v0 - update")
		return rest.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	err = r.usecase.UpdateUser(
		c.Request().Context(),
		entity.UserEntity{
			ID:    id,
			Name:  input.Name,
			Phone: input.Phone,
		},
	)
	if err != nil {
		r.logger.Error(err, "http - v0 - update")
		return rest.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully updated"})
}

// @Summary     Create User
// @Description add new user
// @ID          create
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body request.UserRequest true "Set up users"
// @Success     200 {object} output.UserOutput
// @Failure     400 {object} rest.response
// @Failure     500 {object} rest.response
// @Router      /v0/user [post]
func (r *userRoutes) create(c echo.Context) error {
	var input input.UserInput
	if err := c.Bind(&input); err != nil {
		r.logger.Error(err, "http - v0 - create")
		return rest.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	user, err := r.usecase.CreateUser(
		c.Request().Context(),
		entity.UserEntity{
			Name:  input.Name,
			Phone: input.Phone,
		},
	)
	if err != nil {
		r.logger.Error(err, "http - v0 - create")
		return rest.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
// @Failure     400  {object} rest.response  "Invalid UUID format"
// @Failure     500  {object} rest.response  "Internal server error"
// @Router      /v0/user/{id} [delete]
func (r *userRoutes) delete(c echo.Context) error {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		r.logger.Error(err, "http - v0 - delete")
		return rest.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID format")
	}

	err = r.usecase.DeleteUser(c.Request().Context(), entity.UserEntity{ID: id})
	if err != nil {
		r.logger.Error(err, "http - v0 - delete")
		return rest.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully deleted"})
}
