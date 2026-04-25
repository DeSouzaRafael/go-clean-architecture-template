package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/logger"
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/validator"
	"github.com/DeSouzaRafael/go-clean-architecture-template/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewRouter_HealthEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	l := logger.NewLogger("info")
	v := validator.NewValidator()
	uc := mocks.NewMockUseCases(ctrl)
	mockUser := mocks.NewMockUser(ctrl)
	uc.EXPECT().UserUseCase().Return(mockUser)

	NewRouter(e, l, v, uc, "dev")

	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
}

func TestCorsConfig_Production(t *testing.T) {
	cc := corsConfig("prd")
	assert.Equal(t, []string{"https://*.your.domain.com"}, cc.AllowOrigins)
}

func TestCorsConfig_NonProduction(t *testing.T) {
	for _, env := range []string{"dev", "local", ""} {
		cc := corsConfig(env)
		assert.Equal(t, []string{"*"}, cc.AllowOrigins)
	}
}
