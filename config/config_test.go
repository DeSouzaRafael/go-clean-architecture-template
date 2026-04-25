package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("APP_NAME", "test-app")
	t.Setenv("APP_VERSION", "1.0.0")
	t.Setenv("ENV", "local")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("PG_URL", "postgres://localhost/test")
	t.Setenv("PG_MAX_OPEN_CONNS", "20")
	t.Setenv("PG_MAX_IDLE_CONNS", "10")
	t.Setenv("PG_CONN_MAX_LIFETIME_SEC", "7200")

	cfg, err := NewConfig()
	require.NoError(t, err)

	assert.Equal(t, "test-app", cfg.App.Name)
	assert.Equal(t, "1.0.0", cfg.App.Version)
	assert.Equal(t, "local", cfg.App.Env)
	assert.Equal(t, "8080", cfg.HTTP.Port)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "postgres://localhost/test", cfg.PG.URL)
	assert.Equal(t, 20, cfg.PG.MaxOpenConns)
	assert.Equal(t, 10, cfg.PG.MaxIdleConns)
	assert.Equal(t, 7200*time.Second, cfg.PG.ConnMaxLifetime)
}

func TestNewConfig_Defaults(t *testing.T) {
	os.Unsetenv("PG_MAX_OPEN_CONNS")
	os.Unsetenv("PG_MAX_IDLE_CONNS")
	os.Unsetenv("PG_CONN_MAX_LIFETIME_SEC")

	cfg, err := NewConfig()
	require.NoError(t, err)

	assert.Equal(t, 10, cfg.PG.MaxOpenConns)
	assert.Equal(t, 5, cfg.PG.MaxIdleConns)
	assert.Equal(t, 3600*time.Second, cfg.PG.ConnMaxLifetime)
}

func TestGetEnvInt_InvalidValue(t *testing.T) {
	t.Setenv("TEST_INT", "not-a-number")
	result := getEnvInt("TEST_INT", 42)
	assert.Equal(t, 42, result)
}

func TestGetEnvInt_ValidValue(t *testing.T) {
	t.Setenv("TEST_INT", "99")
	result := getEnvInt("TEST_INT", 42)
	assert.Equal(t, 99, result)
}

func TestGetEnvInt_MissingKey(t *testing.T) {
	os.Unsetenv("TEST_INT_MISSING")
	result := getEnvInt("TEST_INT_MISSING", 42)
	assert.Equal(t, 42, result)
}
