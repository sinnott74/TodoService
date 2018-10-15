package todo

import (
	"context"
	"testing"

	"github.com/rs/xid"

	"github.com/stretchr/testify/require"
)

// TestAddTodo tests adding a Todo and reading it back
func TestAddTodo(t *testing.T) {
	todoService := NewInmemTodoService()

	username := "test@test.com"

	todo := Todo{
		Username: username,
		Text:     "Finish off this microservice",
	}

	addedTodo, err := todoService.Add(context.Background(), todo)

	require.NoError(t, err, "Error adding a Todo")
	require.NotZero(t, addedTodo.ID, "Added todo should have an ID")
	require.NotZero(t, addedTodo.CreatedOn, "Added todo should have a CreatedOn")

	todos, err := todoService.GetAllForUser(context.Background(), username)
	require.NoError(t, err, "Error reading back Todos")
	require.Equal(t, addedTodo, todos[0], "Added Todo should be in list of Todos")

	gottenTodo, err := todoService.GetByID(context.Background(), addedTodo.ID)
	require.NoError(t, err, "Error getting Todo by ID")
	require.Equal(t, addedTodo, gottenTodo, "Added Todo should be in list of Todos")
}

// TestAddTodo tests deleting a Todo
func TestDeleteTodo(t *testing.T) {
	todoService := NewInmemTodoService()

	username := "test@test.com"

	todo := Todo{
		Username: username,
		Text:     "Finish off this microservice",
	}

	addedTodo, err := todoService.Add(context.Background(), todo)

	require.NoError(t, err, "Error adding a Todo")
	require.NotZero(t, addedTodo.ID, "Added todo should have an ID")

	todos, err := todoService.GetAllForUser(context.Background(), "test@test.com")
	require.NoError(t, err, "Error reading back Todos")
	require.Equal(t, 1, len(todos), "Should be only 1 todo")
	require.Equal(t, addedTodo, todos[0], "Added Todo should be in list of Todos")

	err = todoService.Delete(context.Background(), addedTodo.ID)
	require.NoError(t, err, "Error deleting Todos")

	todos, err = todoService.GetAllForUser(context.Background(), username)
	require.NoError(t, err, "Error reading back Todos")
	require.Equal(t, 0, len(todos), "Should be no Todos")

	_, err = todoService.GetByID(context.Background(), addedTodo.ID)
	require.Error(t, err, "ErrNotFound expected")
}

// TestGetAllForUserOnlyReturnsTodosForUser tests that only a user's Todo's are returned
func TestGetAllForUserOnlyReturnsTodosForUser(t *testing.T) {
	todoService := NewInmemTodoService()

	username := "test@test.com"

	todo := Todo{
		Username: username,
		Text:     "Finish off this microservice",
	}

	addedTodo, err := todoService.Add(context.Background(), todo)

	require.NoError(t, err, "Error adding a Todo")
	require.NotZero(t, addedTodo.ID, "Added todo should have an ID")

	todos, err := todoService.GetAllForUser(context.Background(), "testANOTHER@test.com")
	require.NoError(t, err, "Error reading back Todos")
	require.Equal(t, 0, len(todos), "No todos exist for testANOTHER@test.com")
}

// TestUpdateTodo tests updating a Todo
func TestUpdateTodo(t *testing.T) {
	todoService := NewInmemTodoService()

	username := "test@test.com"

	todo := Todo{
		Username:  username,
		Text:      "Finish off this microservice",
		Completed: false,
	}

	addedTodo, err := todoService.Add(context.Background(), todo)
	require.NoError(t, err, "Error adding a Todo")
	require.NotZero(t, addedTodo.ID, "Added todo should have an ID")

	todos, err := todoService.GetAllForUser(context.Background(), username)
	require.NoError(t, err, "Error reading back Todos")
	require.Equal(t, 1, len(todos), "Should be only 1 todo")
	require.Equal(t, addedTodo, todos[0], "Added Todo should be in list of Todos")

	addedTodo.Completed = true
	err = todoService.Update(context.Background(), addedTodo.ID, addedTodo)
	require.NoError(t, err, "Error deleting Todos")

	todos, err = todoService.GetAllForUser(context.Background(), username)
	require.NoError(t, err, "Error reading back Todos")
	require.Equal(t, 1, len(todos), "Should be no Todos")

	gottenTodo, err := todoService.GetByID(context.Background(), addedTodo.ID)
	require.NoError(t, err, "Error getting updated todo by ID")
	require.Equal(t, addedTodo, gottenTodo, "Added Todo should be in list of Todos")
}

// TestDeleteNotFound test deleting a todo by ID which doesn't exist
func TestDeleteNotFound(t *testing.T) {
	todoService := NewInmemTodoService()
	id := xid.New().String()
	err := todoService.Delete(context.Background(), id)
	require.EqualError(t, err, "Not found", "Not found error expected to be returned")
}

// TestUpdateNotFound tests updating a todo which doesn't exist
func TestUpdateNotFound(t *testing.T) {
	todoService := NewInmemTodoService()

	username := "test@test.com"

	todo := Todo{
		ID:        xid.New().String(),
		Username:  username,
		Text:      "Finish off this microservice",
		Completed: false,
	}

	err := todoService.Update(context.Background(), todo.ID, todo)
	require.EqualError(t, err, "Not found", "Not found error expected to be returned")
}

// TestUpdateInconsistentIDs tests thats during an update, when the ID doesn't match the Todos ID an error is returned.
func TestUpdateInconsistentIDs(t *testing.T) {
	todoService := NewInmemTodoService()

	todo := Todo{
		ID:        xid.New().String(),
		Username:  "test@test.com",
		Text:      "Finish off this microservice",
		Completed: false,
	}

	err := todoService.Update(context.Background(), xid.New().String(), todo)
	require.EqualError(t, err, "Inconsistent IDs", "Inconsistent IDs error expected to be returned")
}
