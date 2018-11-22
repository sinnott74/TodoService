package todo

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/stretchr/testify/require"
)

/**
 * In memory service tests
 **/

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
	require.Len(t, todos, 1, "Should be only 1 todo")
	require.Equal(t, addedTodo, todos[0], "Added Todo should be in list of Todos")

	err = todoService.Delete(context.Background(), addedTodo.ID)
	require.NoError(t, err, "Error deleting Todos")

	todos, err = todoService.GetAllForUser(context.Background(), username)
	require.NoError(t, err, "Error reading back Todos")
	require.Len(t, todos, 0, "Should be no Todos")

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
	require.Len(t, todos, 0, "No todos exist for testANOTHER@test.com")
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
	require.Len(t, todos, 1, "Should be only 1 todo")
	require.Equal(t, addedTodo, todos[0], "Added Todo should be in list of Todos")

	addedTodo.Completed = true
	err = todoService.Update(context.Background(), addedTodo.ID, addedTodo)
	require.NoError(t, err, "Error deleting Todos")

	todos, err = todoService.GetAllForUser(context.Background(), username)
	require.NoError(t, err, "Error reading back Todos")
	require.Len(t, todos, 1, "Should be no Todos")

	gottenTodo, err := todoService.GetByID(context.Background(), addedTodo.ID)
	require.NoError(t, err, "Error getting updated todo by ID")
	require.Equal(t, addedTodo, gottenTodo, "Added Todo should be in list of Todos")
}

// TestDeleteNotFound test deleting a todo by ID which doesn't exist
func TestDeleteNotFound(t *testing.T) {
	todoService := NewInmemTodoService()
	id := uuid.NewV4().String()
	err := todoService.Delete(context.Background(), id)
	require.EqualError(t, err, "Not found", "Not found error expected to be returned")
}

// TestUpdateNotFound tests updating a todo which doesn't exist
func TestUpdateNotFound(t *testing.T) {
	todoService := NewInmemTodoService()

	username := "test@test.com"

	todo := Todo{
		ID:        uuid.NewV4().String(),
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
		ID:        uuid.NewV4().String(),
		Username:  "test@test.com",
		Text:      "Finish off this microservice",
		Completed: false,
	}

	err := todoService.Update(context.Background(), uuid.NewV4().String(), todo)
	require.EqualError(t, err, "Inconsistent IDs", "Inconsistent IDs error expected to be returned")
}

/**
 * Postgres service tests
 **/
// TestHealth tests that the service is healthy when it can ping the DB
func TestHealth(t *testing.T) {
	ctx := context.Background()
	db, _, _ := sqlmock.New()
	defer db.Close()

	svc := NewPostgresService(db)
	require.NoError(t, svc.Health(ctx), "Expect service to be healthy")
}

// TestUnhealth tests that the service is unhealthy when it can't ping the DB
func TestUnhealth(t *testing.T) {
	ctx := context.Background()
	db, _, _ := sqlmock.New()
	db.Close() // close db

	svc := NewPostgresService(db)
	require.Error(t, svc.Health(ctx), "Expect service to not be healthy")
}

// TestAddTodoPostgres tests adding a todo with the Postgres service and verifies the SQL that is executed
func TestAddTodoPostgres(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	username := "test@test.com"
	todo := Todo{
		Username: username,
		Text:     "Finish off this microservice",
	}

	query := "INSERT INTO public.todo (id, username, text, completed, createdon, completedon, flagged) VALUES ($1,$2,$3,$4,$5,$6,$7)"

	expectedExec := mock.ExpectExec(regexp.QuoteMeta(query))
	expectedExec.WithArgs(AnyUUID{}, todo.Username, todo.Text, false, AnyTime{}, AnyTime{}, false)
	expectedExec.WillReturnResult(sqlmock.NewResult(1, 1))

	addedTodo, err := todoService.Add(context.Background(), todo)

	require.NoError(t, err, "Error adding a Todo")
	require.NotZero(t, addedTodo.ID, "Added todo should have an ID")
	require.NotZero(t, addedTodo.CreatedOn, "Added todo should have a CreatedOn")
}

// TestAddTodoFailedPostgres test a failed SQL call during the insert
func TestAddTodoFailedPostgres(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	username := "test@test.com"
	todo := Todo{
		Username: username,
		Text:     "Finish off this microservice",
	}

	query := "INSERT INTO public.todo (id, username, text, completed, createdon, completedon, flagged) VALUES ($1,$2,$3,$4,$5,$6,$7)"

	errDB := errors.New("Postgres error")
	expectedExec := mock.ExpectExec(regexp.QuoteMeta(query))
	expectedExec.WithArgs(AnyUUID{}, todo.Username, todo.Text, false, AnyTime{}, AnyTime{}, false)
	expectedExec.WillReturnError(errDB)

	addedTodo, err := todoService.Add(context.Background(), todo)

	require.Error(t, err, "Expected an error when doing sql insert")
	require.Zero(t, addedTodo.ID, "Added todo should have an ID")
}

// TestAddTodoNoRowsUpdatedPostgres tests that an error is returned when no rows are updated during a Todo add
func TestAddTodoNoRowsUpdatedPostgres(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	username := "test@test.com"
	todo := Todo{
		Username: username,
		Text:     "Finish off this microservice",
	}

	query := "INSERT INTO public.todo (id, username, text, completed, createdon, completedon, flagged) VALUES ($1,$2,$3,$4,$5,$6,$7)"

	expectedExec := mock.ExpectExec(regexp.QuoteMeta(query))
	expectedExec.WithArgs(AnyUUID{}, todo.Username, todo.Text, false, AnyTime{}, AnyTime{}, false)
	expectedExec.WillReturnResult(sqlmock.NewResult(1, 0))

	addedTodo, err := todoService.Add(context.Background(), todo)

	require.Error(t, err, "Expected error are no rows were updated adding a new Todo")
	require.Zero(t, addedTodo.ID, "Added todo should have an ID")
}

func TestDeletePostgres(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	id := uuid.NewV4().String()

	query := "DELETE FROM public.todo WHERE id = $1"

	expectedExec := mock.ExpectExec(regexp.QuoteMeta(query))
	expectedExec.WithArgs(id)
	expectedExec.WillReturnResult(sqlmock.NewResult(1, 1))

	err := todoService.Delete(context.Background(), id)

	require.NoError(t, err, "Error deleting a Todo")
}

func TestDeleteNotFoundPostgres(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	id := uuid.NewV4().String()

	query := "DELETE FROM public.todo WHERE id = $1"

	expectedExec := mock.ExpectExec(regexp.QuoteMeta(query))
	expectedExec.WithArgs(id)
	expectedExec.WillReturnResult(sqlmock.NewResult(1, 0))

	err := todoService.Delete(context.Background(), id)

	require.EqualError(t, err, ErrNotFound.Error(), "Expected Not Found Error when delete a Todo")
}

func TestDeleteFailedPostgres(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	id := uuid.NewV4().String()

	query := "DELETE FROM public.todo WHERE id = $1"

	errDB := errors.New("Postgres error")
	expectedExec := mock.ExpectExec(regexp.QuoteMeta(query))
	expectedExec.WithArgs(id)
	expectedExec.WillReturnError(errDB)

	err := todoService.Delete(context.Background(), id)

	require.Error(t, err, "Expected error when delete a Todo")
}

func TestGetByIDPostgres(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	todo := Todo{
		ID:          uuid.NewV4().String(),
		Username:    "test@test.com",
		Text:        "// TODO",
		Completed:   false,
		CreatedOn:   time.Now(),
		CompletedOn: time.Time{}, // zero time
		Flagged:     true,
	}

	query := "SELECT * FROM public.todo WHERE id = $1"

	mockedRow := sqlmock.NewRows([]string{"id", "username", "text", "completed", "create_on", "completed_on", "flagged"})
	mockedRow.AddRow(todo.ID, todo.Username, todo.Text, todo.Completed, todo.CreatedOn, todo.CompletedOn, todo.Flagged)

	expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(query))
	expectedQuery.WithArgs(todo.ID)
	expectedQuery.WillReturnRows(mockedRow)

	todoByID, err := todoService.GetByID(context.Background(), todo.ID)

	require.NoError(t, err, "Unexpected error when getting a Todo by ID")
	require.Equal(t, todo, todoByID, "Todos should have matched")
}

func TestGetByIDErrNotFoundPostgres(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	todo := Todo{
		ID: uuid.NewV4().String(),
	}

	query := "SELECT * FROM public.todo WHERE id = $1"

	mockedRow := sqlmock.NewRows([]string{"id", "username", "text", "completed", "create_on", "completed_on", "flagged"})

	expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(query))
	expectedQuery.WithArgs(todo.ID)
	expectedQuery.WillReturnRows(mockedRow)

	todoByID, err := todoService.GetByID(context.Background(), todo.ID)

	require.EqualError(t, err, ErrNotFound.Error(), "Expected ErrNotFound")
	require.Zero(t, todoByID, "Todo returned should be blank")
}

func TestGetByIDErrorPostgres(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	id := uuid.NewV4().String()

	query := "SELECT * FROM public.todo WHERE id = $1"

	errDB := errors.New("Postgres error")
	expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(query))
	expectedQuery.WithArgs(id)
	expectedQuery.WillReturnError(errDB)

	todoByID, err := todoService.GetByID(context.Background(), id)

	require.EqualError(t, err, errDB.Error(), "Expected mocked sql error")
	require.Zero(t, todoByID, "Todo returned should be blank")
}

// TestUpdatePostgres tests a successful update
func TestUpdatePostgres(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	todo := Todo{
		ID:          uuid.NewV4().String(),
		Username:    "test@test.com",
		Text:        "// TODO",
		Completed:   false,
		CreatedOn:   time.Now(),
		CompletedOn: time.Time{}, // zero time
		Flagged:     true,
	}

	query := "UPDATE public.todo SET id=$1, username=$2, text=$3, completed=$4, createdon=$5, completedon=$6, flagged=$7 WHERE id = $8"

	expectedQuery := mock.ExpectExec(regexp.QuoteMeta(query))
	expectedQuery.WithArgs(todo.ID, todo.Username, todo.Text, todo.Completed, todo.CreatedOn, todo.CompletedOn, todo.Flagged, todo.ID)
	expectedQuery.WillReturnResult(sqlmock.NewResult(0, 1))

	err := todoService.Update(context.Background(), todo.ID, todo)

	require.NoError(t, err, "Unexpected error when updating Todo")
}

// TestUpdateErrNotFoundPostgres tests that ErrNotFound is returned when no rows were updated
func TestUpdateErrNotFoundPostgres(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	todo := Todo{
		ID:          uuid.NewV4().String(),
		Username:    "test@test.com",
		Text:        "// TODO",
		Completed:   false,
		CreatedOn:   time.Now(),
		CompletedOn: time.Time{}, // zero time
		Flagged:     true,
	}

	query := "UPDATE public.todo SET id=$1, username=$2, text=$3, completed=$4, createdon=$5, completedon=$6, flagged=$7 WHERE id = $8"

	expectedQuery := mock.ExpectExec(regexp.QuoteMeta(query))
	expectedQuery.WithArgs(todo.ID, todo.Username, todo.Text, todo.Completed, todo.CreatedOn, todo.CompletedOn, todo.Flagged, todo.ID)
	expectedQuery.WillReturnResult(sqlmock.NewResult(0, 0))

	err := todoService.Update(context.Background(), todo.ID, todo)

	require.EqualError(t, err, ErrNotFound.Error(), "Expected ErrNotFound as no rows were updated")
}

// TestUpdateErrorPostgres tests that an error is returned when an error occurs during the SQL exec
func TestUpdateErrorPostgres(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	todo := Todo{
		ID:          uuid.NewV4().String(),
		Username:    "test@test.com",
		Text:        "// TODO",
		Completed:   false,
		CreatedOn:   time.Now(),
		CompletedOn: time.Time{}, // zero time
		Flagged:     true,
	}

	query := "UPDATE public.todo SET id=$1, username=$2, text=$3, completed=$4, createdon=$5, completedon=$6, flagged=$7 WHERE id = $8"

	errDB := errors.New("Postgres error")
	expectedQuery := mock.ExpectExec(regexp.QuoteMeta(query))
	expectedQuery.WithArgs(todo.ID, todo.Username, todo.Text, todo.Completed, todo.CreatedOn, todo.CompletedOn, todo.Flagged, todo.ID)
	expectedQuery.WillReturnError(errDB)

	err := todoService.Update(context.Background(), todo.ID, todo)

	require.EqualError(t, err, errDB.Error(), "Expected Update to return an error")
}

// TestGetAllForUserNoRowsPostgres tests a successful GetAllForUser which returns no Todos
func TestGetAllForUserNoRowsPostgres(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)
	username := "test@test.com"

	query := "SELECT * FROM public.todo WHERE username=$1"

	mockedRows := sqlmock.NewRows([]string{"id", "username", "text", "completed", "create_on", "completed_on", "flagged"})

	expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(query))
	expectedQuery.WithArgs(username)
	expectedQuery.WillReturnRows(mockedRows)

	todos, err := todoService.GetAllForUser(context.Background(), username)

	require.NoError(t, err, "Unexpected error when GetAllForUser")
	require.Empty(t, todos, "Expected no Todos to be returned")
}

func TestGetAllForUserSingleRowPostgres(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()
	todoService := NewPostgresService(db)

	todo := Todo{
		ID:          uuid.NewV4().String(),
		Username:    "test@test.com",
		Text:        "// TODO",
		Completed:   false,
		CreatedOn:   time.Now(),
		CompletedOn: time.Time{}, // zero time
		Flagged:     true,
	}

	query := "SELECT * FROM public.todo WHERE username=$1"

	mockedRows := sqlmock.NewRows([]string{"id", "username", "text", "completed", "create_on", "completed_on", "flagged"})
	mockedRows.AddRow(todo.ID, todo.Username, todo.Text, todo.Completed, todo.CreatedOn, todo.CompletedOn, todo.Flagged)

	expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(query))
	expectedQuery.WithArgs(todo.Username)
	expectedQuery.WillReturnRows(mockedRows)

	todos, err := todoService.GetAllForUser(context.Background(), todo.Username)

	require.NoError(t, err, "Unexpected error when GetAllForUser")
	require.NotEmpty(t, todos, "Expected Todos to be returned")
	require.Len(t, todos, 1, "Expected Todos to be returned")
}

/**
 * sqlmock.Arguments for sql variable matching
 **/
// AnyUUID is a sqlmock.Arguement to match uuid.UUID
type AnyUUID struct{}

func (AnyUUID) Match(v driver.Value) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, err := uuid.FromString(s)
	return err == nil
}

// AnyTime is a sqlmock.Arguement to match time.Time
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
