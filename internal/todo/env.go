package todo

import (
	"os"
	"strings"
)

// JWTSecret to be used in during authentication
func JWTSecret() []byte {
	env := os.Getenv("JWT_SECRET")
	if env == "" {
		env = "SECRET_SSSHHHHHHH"
	}
	return []byte(env)
}

// ConnectionURL get the database connection string from ENV Vars or used a default
func ConnectionURL() string {
	connectionString := os.Getenv("POSTGRES_URL")
	if connectionString == "" {
		connectionString = "postgres://Sinnott@localhost:5432/tododb?sslmode=disable&timezone=UTC"
	}
	return connectionString
}

// Port retrieves the Port to start the server on
func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	return port
}

// Debug retrieves turns on debugging
func Debug() bool {
	debug := os.Getenv("DEBUG")
	if strings.ToLower(debug) == "true" {
		return true
	}
	return false
}
