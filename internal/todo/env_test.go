package todo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPortDefault checks that the default PORT is returned when not set
func TestPortDefault(t *testing.T) {
	port := Port()
	assert.Equal(t, "8000", port)
}

// TestPortEnvSet checks that the correct PORT is returned when set
func TestPortEnvSet(t *testing.T) {
	expectedPort := "7474"
	os.Setenv("PORT", expectedPort)
	port := Port()
	assert.Equal(t, expectedPort, port)
	os.Unsetenv("PORT")
}

// TestPortDefault checks that the default PORT is returned when not set
func TestJWTSecretDefault(t *testing.T) {
	secret := JWTSecret()
	assert.Equal(t, []byte("SECRET_SSSHHHHHHH"), secret)
}

// TestPortEnvSet checks that the correct PORT is returned when set
func TestJWTSecretEnvSet(t *testing.T) {
	expectedSecret := "3rd secret of fatima"
	os.Setenv("JWT_SECRET", expectedSecret)
	secret := JWTSecret()
	assert.Equal(t, []byte(expectedSecret), secret)
	os.Unsetenv("JWT_SECRET")
}

// TestPortDefault checks that the default PORT is returned when not set
func TestConnectionURLDefault(t *testing.T) {
	connURL := ConnectionURL()
	assert.Equal(t, "postgres://Sinnott@localhost:5432/tododb?sslmode=disable&timezone=UTC", connURL)
}

// TestPortEnvSet checks that the correct PORT is returned when set
func TestConnectionURLEnvSet(t *testing.T) {
	expectedConnURL := "postgres://test"
	os.Setenv("POSTGRES_URL", expectedConnURL)
	connURL := ConnectionURL()
	assert.Equal(t, expectedConnURL, connURL)
	os.Unsetenv("POSTGRES_URL")
}

// TestPortDefault checks that the default PORT is returned when not set
func TestDebugDefault(t *testing.T) {
	debug := Debug()
	assert.Equal(t, false, debug)
}

// TestPortEnvSet checks that the correct PORT is returned when set
func TestDebugEnvSet(t *testing.T) {
	os.Setenv("DEBUG", "true")
	debug := Debug()
	assert.Equal(t, true, debug)
	os.Unsetenv("DEBUG")
}
