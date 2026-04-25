package httpserver

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew_DefaultsAndStart(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s := New(handler, Port("0"), ShutdownTimeout(1*time.Second))
	assert.NotNil(t, s)
	assert.NotNil(t, s.Notify())

	err := s.Shutdown()
	assert.NoError(t, err)
}

func TestPort(t *testing.T) {
	s := &Server{server: &http.Server{}}
	Port("9090")(s)
	assert.Equal(t, ":9090", s.server.Addr)
}

func TestReadTimeout(t *testing.T) {
	s := &Server{server: &http.Server{}}
	ReadTimeout(10 * time.Second)(s)
	assert.Equal(t, 10*time.Second, s.server.ReadTimeout)
}

func TestWriteTimeout(t *testing.T) {
	s := &Server{server: &http.Server{}}
	WriteTimeout(15 * time.Second)(s)
	assert.Equal(t, 15*time.Second, s.server.WriteTimeout)
}

func TestShutdownTimeout(t *testing.T) {
	s := &Server{server: &http.Server{}, notify: make(chan error, 1)}
	ShutdownTimeout(5 * time.Second)(s)
	assert.Equal(t, 5*time.Second, s.shutdownTimeout)
}

func TestNotify(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	s := New(handler, Port("0"), ShutdownTimeout(1*time.Second))
	ch := s.Notify()
	assert.NotNil(t, ch)
	_ = s.Shutdown()
}
