package output

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ErrorResponse(c, http.StatusBadRequest, "something failed")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "something failed")
}

func TestUserOutput_Fields(t *testing.T) {
	id := uuid.New()
	out := UserOutput{ID: id, Name: "Jane", Phone: "+5511999999999"}
	assert.Equal(t, id, out.ID)
	assert.Equal(t, "Jane", out.Name)
	assert.Equal(t, "+5511999999999", out.Phone)
}
