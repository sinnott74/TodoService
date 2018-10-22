package main

import (
	"database/sql"
	"net/http"
	"time"

	_ "github.com/lib/pq" // import postgres
	"github.com/sinnott74/TodoService/internal/todo"
)

func main() {

	db, err := initDB()
	if err != nil {
		panic(err)
	}

	service := todo.NewPostgresService(db)
	// service := todo.NewInmemTodoService()

	endpoints := todo.MakeTodoEndpoints(service)

	srv := &http.Server{
		Addr:         ":" + todo.Port(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      todo.MakeHTTPHandler(endpoints),
	}

	err = srv.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", todo.ConnectionURL())
	if err != nil {
		return db, err
	}
	err = db.Ping()
	return db, err
}
