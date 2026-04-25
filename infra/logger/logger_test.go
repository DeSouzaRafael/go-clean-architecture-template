package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger_Levels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "unknown"}
	for _, level := range levels {
		l := NewLogger(level)
		assert.NotNil(t, l)
	}
}

func TestLogger_Info(t *testing.T) {
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf)
	l := &Logger{logger: &zl}
	l.Info("hello world")
	assert.Contains(t, buf.String(), "hello world")
}

func TestLogger_InfoWithArgs(t *testing.T) {
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf)
	l := &Logger{logger: &zl}
	l.Info("value: %d", 42)
	assert.Contains(t, buf.String(), "42")
}

func TestLogger_Warn(t *testing.T) {
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf)
	l := &Logger{logger: &zl}
	l.Warn("a warning")
	assert.Contains(t, buf.String(), "a warning")
}

func TestLogger_WarnWithArgs(t *testing.T) {
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf)
	l := &Logger{logger: &zl}
	l.Warn("warn value: %s", "test")
	assert.Contains(t, buf.String(), "test")
}

func TestLogger_Debug(t *testing.T) {
	buf := &bytes.Buffer{}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zl := zerolog.New(buf).Level(zerolog.DebugLevel)
	l := &Logger{logger: &zl}
	l.Debug("debug msg")
	assert.Contains(t, buf.String(), "debug msg")
}

func TestLogger_DebugWithString(t *testing.T) {
	buf := &bytes.Buffer{}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zl := zerolog.New(buf).Level(zerolog.DebugLevel)
	l := &Logger{logger: &zl}
	l.Debug("string message")
	assert.Contains(t, buf.String(), "string message")
}

func TestLogger_DebugWithError(t *testing.T) {
	buf := &bytes.Buffer{}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zl := zerolog.New(buf).Level(zerolog.DebugLevel)
	l := &Logger{logger: &zl}
	l.Debug(fmt.Errorf("error obj"))
	assert.Contains(t, buf.String(), "error obj")
}

func TestLogger_Error_AtDebugLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zl := zerolog.New(buf).Level(zerolog.DebugLevel)
	l := &Logger{logger: &zl}
	l.Error(fmt.Errorf("something went wrong"))
	assert.Contains(t, buf.String(), "something went wrong")
}

func TestLogger_Error_WithString(t *testing.T) {
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf)
	l := &Logger{logger: &zl}
	l.Error("string error")
	assert.Contains(t, buf.String(), "string error")
}

func TestLogger_Error_UnknownType(t *testing.T) {
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf)
	l := &Logger{logger: &zl}
	l.Error(12345)
	assert.Contains(t, buf.String(), "unknown type")
}

func TestLogger_Fatal(t *testing.T) {
	if !testing.Short() {
		t.Skip("Skipping Fatal test that calls os.Exit")
	}
}
