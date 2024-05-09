// user_view_test.go
package v0

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/logger"
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/validator"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
	"github.com/DeSouzaRafael/go-clean-architecture-template/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewUserRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()
	NewUserRoutes(e, l, v, mockUseCase)

	// Verify that routes are configured correctly
	assertRouteExists(t, e, http.MethodGet, "/v0/user/:id")
	assertRouteExists(t, e, http.MethodPost, "/v0/user")
	assertRouteExists(t, e, http.MethodPut, "/v0/user/:id")
	assertRouteExists(t, e, http.MethodDelete, "/v0/user/:id")
}

func assertRouteExists(t *testing.T, e *echo.Echo, method, path string) {
	found := false
	for _, r := range e.Routes() {
		if r.Method == method && r.Path == path {
			found = true
			break
		}
	}
	assert.True(t, found, "expected route %s %s to exist", method, path)
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()
	userID := uuid.New()
	user := entity.UserEntity{
		ID:   userID,
		Name: "User Name",
	}

	mockUseCase.EXPECT().GetUserById(gomock.Any(), entity.UserEntity{ID: userID}).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/v0/user/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	r := userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.get(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"name":"User Name"`)
	}
}

func TestGetUser_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()
	userID := uuid.New()

	mockUseCase.EXPECT().GetUserById(gomock.Any(), entity.UserEntity{ID: userID}).Return(entity.UserEntity{}, fmt.Errorf("some error"))

	req := httptest.NewRequest(http.MethodGet, "/v0/user/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.get(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), `"error":"some error"`)
	}
}

func TestGetUser_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/v0/user/invalid-uuid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues("invalid-uuid")

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.get(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `"error":"invalid UUID format"`)
	}
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()
	reqBody := `{"name": "User Name", "phone": "+5511999999999"}`
	user := entity.UserEntity{
		ID:    uuid.New(),
		Name:  "User Name",
		Phone: "+5511999999999",
	}

	mockUseCase.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(user, nil)

	req := httptest.NewRequest(http.MethodPost, "/v0/user", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user")

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.create(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"name":"User Name"`)
	}
}

func TestCreateUser_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/v0/user", strings.NewReader("invalid-json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user")

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.create(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `"error":"invalid request body"`)
	}
}

func TestCreateUser_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")

	e := echo.New()
	reqBody := `{"name": "", "phone": "+5511999999999"}`
	validator := validator.NewValidator()

	req := httptest.NewRequest(http.MethodPut, "/v0/user/"+entity.UserEntity{}.ID.String(), strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues(entity.UserEntity{}.ID.String())

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: validator}
	err := r.create(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `{"error":"invalid request data: Field validation for 'Name' failed on the 'required' tag."}`+"\n")
	}
}

func TestCreateUser_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()
	reqBody := `{"name": "User Name", "phone": "+5511999999999"}`

	mockUseCase.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(entity.UserEntity{}, fmt.Errorf("some error"))

	req := httptest.NewRequest(http.MethodPost, "/v0/user", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user")

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.create(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), `"error":"some error"`)
	}
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()
	userID := uuid.New()
	reqBody := `{"name": "User Name", "phone": "+5511999999999"}`
	user := entity.UserEntity{
		ID:    userID,
		Name:  "User Name",
		Phone: "+5511999999999",
	}

	mockUseCase.EXPECT().UpdateUser(gomock.Any(), user).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/v0/user/"+userID.String(), strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.update(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"message":"Successfully updated"`)
	}
}

func TestUpdateUser_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()
	userID := uuid.New()
	reqBody := `{"name": "User Name", "phone": "+5511999999999"}`
	user := entity.UserEntity{
		ID:    userID,
		Name:  "User Name",
		Phone: "+5511999999999",
	}

	mockUseCase.EXPECT().UpdateUser(gomock.Any(), user).Return(fmt.Errorf("some error"))

	req := httptest.NewRequest(http.MethodPut, "/v0/user/"+userID.String(), strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.update(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), `"error":"some error"`)
	}
}

func TestUpdateUser_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")

	e := echo.New()
	reqBody := `{"name": "", "phone": "+5511999999999"}`
	validator := validator.NewValidator()

	req := httptest.NewRequest(http.MethodPut, "/v0/user/"+entity.UserEntity{}.ID.String(), strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues(entity.UserEntity{}.ID.String())

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: validator}
	err := r.update(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `{"error":"invalid request data: Field validation for 'Name' failed on the 'required' tag."}`+"\n")
	}
}

func TestUpdateUser_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()
	userID := uuid.New()

	req := httptest.NewRequest(http.MethodPut, "/v0/user/"+userID.String(), strings.NewReader("invalid-json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.update(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `"error":"invalid request body"`)
	}
}

func TestUpdateUser_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()

	req := httptest.NewRequest(http.MethodPut, "/v0/user/invalid-uuid", strings.NewReader(`{"name": "Jane Doe", "phone": "+551199999999"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues("invalid-uuid")

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.update(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `"error":"invalid UUID format"`)
	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()
	userID := uuid.New()

	mockUseCase.EXPECT().DeleteUser(gomock.Any(), entity.UserEntity{ID: userID}).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/v0/user/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.delete(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"message":"Successfully deleted"`)
	}
}

func TestDeleteUser_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/v0/user/invalid-uuid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues("invalid-uuid")

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.delete(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `"error":"Invalid UUID format"`)
	}
}

func TestDeleteUser_DeleteUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUser(ctrl)
	l := logger.NewLogger("info")
	v := validator.NewValidator()

	e := echo.New()
	userID := uuid.New()
	user := entity.UserEntity{ID: userID}

	mockUseCase.EXPECT().DeleteUser(gomock.Any(), user).Return(fmt.Errorf("failed to delete user"))

	req := httptest.NewRequest(http.MethodDelete, "/v0/user/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v0/user/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	r := &userRoutes{usecase: mockUseCase, logger: l, validator: v}
	err := r.delete(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), `"error":"failed to delete user"`)
	}
}
