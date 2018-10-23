package todo

import (
	"time"
)

// Todo model
type Todo struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Text        string    `json:"text"`
	Completed   bool      `json:"completed"`
	CreatedOn   time.Time `json:"created_on"`
	CompletedOn time.Time `json:"completed_on"`
	Flagged     bool      `json:"flagged"`
}

// CREATE TABLE public.todos
// (
//     id uuid NOT NULL,
//     username character varying(50) COLLATE pg_catalog."default" NOT NULL,
//     text text COLLATE pg_catalog."default",
//     completed boolean DEFAULT false,
//     createdon timestamp with time zone NOT NULL,
//     completedon timestamp with time zone,
//		 flagged boolean DEFAULT false,
//     CONSTRAINT todos_pkey PRIMARY KEY (id)
// )
