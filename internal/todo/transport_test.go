package todo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

// TestCreatingATodoThenUpdatingIt test the Todo service's HTTP transport layer.
// It creates a Todo using the exposed HTTP API, then reads it back.
func TestCreatingATodoThenUpdatingTheDeleting(t *testing.T) {

	todoService := NewInmemTodoService()
	endpoints := MakeTodoEndpoints(todoService)
	server := httptest.NewServer(MakeHTTPHandler(endpoints))
	defer server.Close()

	// Create Todo
	todo := Todo{Username: "test@test.com", Text: "Get this service tested"}
	res := newHTTPServerCall(t, http.MethodPost, server.URL+"/api/todos", todo)
	defer res.Body.Close()
	require.Equalf(t, http.StatusOK, res.StatusCode, "Expecting StatusOK when creating Todo")

	var addResponse AddResponse
	json.NewDecoder(res.Body).Decode(&addResponse)
	require.Equalf(t, todo.Username, addResponse.Todo.Username, "Username isn't the same")
	require.Equalf(t, todo.Text, addResponse.Todo.Text, "Todo Text isn't the same")
	require.Equalf(t, todo.Completed, addResponse.Todo.Completed, "Todo completion status isn't the same")
	require.NotZerof(t, addResponse.Todo.ID, "Todo ID should be set")

	// Verify created correctly by Getting By ID
	res = newHTTPServerCall(t, http.MethodGet, server.URL+"/api/todos/"+addResponse.Todo.ID, nil)
	defer res.Body.Close()
	require.Equalf(t, http.StatusOK, res.StatusCode, "Expecting StatusOK when reading Todo by ID")

	var getByIDResponse GetByIDResponse
	json.NewDecoder(res.Body).Decode(&getByIDResponse)
	require.EqualValuesf(t, addResponse.Todo, getByIDResponse.Todo, "Created Todo isn't the same as the Todo By ID")

	// Updating Todo
	todo = getByIDResponse.Todo
	todo.Text = "Finish testing this service"
	res = newHTTPServerCall(t, http.MethodPut, server.URL+"/api/todos/"+todo.ID, todo)
	defer res.Body.Close()
	require.Equalf(t, http.StatusOK, res.StatusCode, "Expecting StatusOK when updating Todo")

	// Verify updated correctly by Getting By ID
	res = newHTTPServerCall(t, http.MethodGet, server.URL+"/api/todos/"+todo.ID, nil)
	defer res.Body.Close()
	require.Equalf(t, http.StatusOK, res.StatusCode, "Expecting StatusOK when reading Todo by ID")

	json.NewDecoder(res.Body).Decode(&getByIDResponse)
	require.EqualValuesf(t, todo, getByIDResponse.Todo, "Created Todo isn't the same as the Todo By ID")

	// Delete Todo
	res = newHTTPServerCall(t, http.MethodDelete, server.URL+"/api/todos/"+todo.ID, nil)
	defer res.Body.Close()
	require.Equalf(t, http.StatusOK, res.StatusCode, "Expecting StatusOK when deleting Todo")

	// Verify deleted correctly when GettingByID - 404
	res = newHTTPServerCall(t, http.MethodGet, server.URL+"/api/todos/"+todo.ID, nil)
	defer res.Body.Close()
	require.Equalf(t, http.StatusNotFound, res.StatusCode, "Expecting 404 when reading deleted Todo by ID")

	var errorMap map[string]string
	json.NewDecoder(res.Body).Decode(&errorMap)
	require.EqualValuesf(t, "Not found", errorMap["error"], "Expected Not found error reading deleted Todo")
}

// newJWTToken creates A JWT token to be used in a request
func newJWTToken(t *testing.T) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "test@test.com",
	})
	tokenString, err := token.SignedString(JWTSecret())
	require.NoError(t, err, "Error creating JWT token")
	return "JWT " + tokenString
}

// NewHTTPServerCall performs a http call
// It sets the request with all required headers. i.e. JWT token
func newHTTPServerCall(t *testing.T, httpMethod, url string, payload interface{}) *http.Response {
	var req *http.Request
	var err error

	if httpMethod == http.MethodGet || httpMethod == http.MethodDelete {
		req, err = http.NewRequest(httpMethod, url, nil)
	} else {
		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(payload)
		req, err = http.NewRequest(httpMethod, url, b)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", newJWTToken(t))
	require.NoErrorf(t, err, "Error creating %s request", httpMethod)
	res, err := http.DefaultClient.Do(req)
	require.NoErrorf(t, err, "Error doing %s request to %s with payload %v", httpMethod, url, payload)
	return res
}
