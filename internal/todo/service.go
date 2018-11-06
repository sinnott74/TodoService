package todo

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/rs/xid"
)

// TodoService for Todos
type TodoService interface {
	GetAllForUser(ctx context.Context, username string) ([]Todo, error)
	GetByID(ctx context.Context, id string) (Todo, error)
	Add(ctx context.Context, todo Todo) (Todo, error)
	Update(ctx context.Context, id string, todo Todo) error
	Delete(ctx context.Context, id string) error
	Health(ctx context.Context) error
}

// *** Implementation ***

var (
	// ErrInconsistentIDs is when the ID of the Entity you are updating differs from the ID given
	ErrInconsistentIDs = errors.New("Inconsistent IDs")
	// ErrNotFound is when the Entity doesn't exist
	ErrNotFound = errors.New("Not found")
)

// NewPostgresService creates a Todo service which uses Postgres for persistence
func NewPostgresService(db *sql.DB) TodoService {
	return &postgresService{db}
}

// psqlService is a Postgres implementation of the service
type postgresService struct {
	db *sql.DB
}

// Get all Todos from the database
func (s *postgresService) GetAllForUser(ctx context.Context, username string) ([]Todo, error) {

	todos := []Todo{}
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM public.todos WHERE username=$1", username)
	if err != nil {
		return todos, err
	}
	defer rows.Close()

	for rows.Next() {
		todo := Todo{}
		err = rows.Scan(&todo.ID, &todo.Username, &todo.Text, &todo.Completed, &todo.CreatedOn, &todo.CompletedOn, &todo.Flagged)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}

	err = rows.Err()

	return todos, err
}

// Get an Todos from the database
func (s *postgresService) GetByID(ctx context.Context, id string) (Todo, error) {
	todo := Todo{}
	row := s.db.QueryRowContext(ctx, "SELECT * FROM public.todos WHERE id = $1", id)
	err := row.Scan(&todo.ID, &todo.Username, &todo.Text, &todo.Completed, &todo.CreatedOn, &todo.CompletedOn, &todo.Flagged)
	return todo, err
}

// Add a Todo to the database
func (s *postgresService) Add(ctx context.Context, todo Todo) (Todo, error) {

	todo.ID = uuid.NewV4().String()
	todo.CreatedOn = time.Now().UTC()
	_, err := s.db.ExecContext(ctx, "INSERT INTO public.todos (id, username, text, completed, createdon, completedon, flagged) VALUES ($1,$2,$3,$4,$5,$6,$7)", todo.ID, todo.Username, todo.Text, todo.Completed, todo.CreatedOn, todo.CompletedOn, todo.Flagged)
	return todo, err
}

// Update a Todo in the database
func (s *postgresService) Update(ctx context.Context, id string, todo Todo) error {
	_, err := s.db.ExecContext(ctx, "UPDATE public.todos SET id=$1, username=$2, text=$3, completed=$4, createdon=$5, completedon=$6, flagged=$7 WHERE id = $8", todo.ID, todo.Username, todo.Text, todo.Completed, todo.CreatedOn, todo.CompletedOn, todo.Flagged, id)
	return err
}

// Delete a Todo in the database
func (s *postgresService) Delete(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM public.todos WHERE id = $1", id)
	if err != nil {
		return err
	}

	numRowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numRowAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Health perform service health check. verifies that DB is accessible
func (s *postgresService) Health(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// NewInmemTodoService creates an in memory Todo service
func NewInmemTodoService() TodoService {
	s := &inmemService{
		m: map[string]Todo{},
	}
	rand.Seed(time.Now().UnixNano())
	return s
}

// inmemService is a In Memory implementation of the service
type inmemService struct {
	sync.RWMutex
	m map[string]Todo
}

// GetAllForUser gets Todos from memory for a user
func (s *inmemService) GetAllForUser(ctx context.Context, username string) ([]Todo, error) {
	s.RLock()
	defer s.RUnlock()

	todos := make([]Todo, 0, len(s.m))
	for _, todo := range s.m {
		if todo.Username == username {
			todos = append(todos, todo)
		}
	}

	return todos, nil
}

// Get an Todos from the database
func (s *inmemService) GetByID(ctx context.Context, id string) (Todo, error) {
	s.Lock()
	defer s.Unlock()

	if todo, ok := s.m[id]; ok {
		return todo, nil
	}

	return Todo{}, ErrNotFound
}

// Add a Todo to memory
func (s *inmemService) Add(ctx context.Context, todo Todo) (Todo, error) {
	s.Lock()
	defer s.Unlock()

	todo.ID = xid.New().String()
	todo.CreatedOn = time.Now()

	s.m[todo.ID] = todo
	return todo, nil
}

// Update a Todo in memory
func (s *inmemService) Update(ctx context.Context, id string, todo Todo) error {
	s.Lock()
	defer s.Unlock()

	if id != todo.ID {
		return ErrInconsistentIDs
	}

	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}

	s.m[todo.ID] = todo
	return nil
}

// Delete a Todo from memory
func (s *inmemService) Delete(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}

	delete(s.m, id)
	return nil
}

// Health check the In memory TodoService business process
func (s *inmemService) Health(ctx context.Context) error {
	todo := Todo{}
	addedTodo, err := s.Add(ctx, todo)
	if err != nil {
		return err
	}
	retrievedTodo, err := s.GetByID(ctx, addedTodo.ID)
	if err != nil {
		return err
	}
	if addedTodo != retrievedTodo {
		return errors.New("health check error retrieving todo")
	}
	return s.Delete(ctx, addedTodo.ID)
}
