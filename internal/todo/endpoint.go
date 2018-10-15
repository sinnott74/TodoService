package todo

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all endpoints which compose the Todo service
type TodoEndpoints struct {
	GetAllForUserEndPoint endpoint.Endpoint
	GetByIDEndpoint       endpoint.Endpoint
	AddEndpoint           endpoint.Endpoint
	UpdateEndpoint        endpoint.Endpoint
	DeleteEndpoint        endpoint.Endpoint
}

// MakeTodoEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided Todo
func MakeTodoEndpoints(s TodoService) TodoEndpoints {
	return TodoEndpoints{
		GetAllForUserEndPoint: MakeGetAllForUserEndpoint(s),
		GetByIDEndpoint:       MakeGetByIDEndpoint(s),
		AddEndpoint:           MakeAddEndpoint(s),
		UpdateEndpoint:        MakeUpdateEndpoint(s),
		DeleteEndpoint:        MakeDeleteEndpoint(s),
	}
}

type GetAllForUserRequest struct {
}

type GetAllForUserResponse struct {
	Todos []Todo `json:"todos"`
}

func MakeGetAllForUserEndpoint(s TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		username := ctx.Value("username").(string)
		todos, err := s.GetAllForUser(ctx, username)
		return GetAllForUserResponse{todos}, err
	}
}

type GetByIDRequest struct {
	ID string
}

type GetByIDResponse struct {
	Todo Todo `json:"todo"`
}

func MakeGetByIDEndpoint(s TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		todo, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{todo}, err
	}
}

type AddRequest struct {
	Todo Todo
}

type AddResponse struct {
	Todo Todo `json:"todo"`
}

func MakeAddEndpoint(s TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		todo, err := s.Add(ctx, req.Todo)
		return AddResponse{todo}, err
	}
}

type UpdateRequest struct {
	ID   string
	Todo Todo
}

type UpdateResponse struct {
}

func MakeUpdateEndpoint(s TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		err := s.Update(ctx, req.ID, req.Todo)
		return UpdateResponse{}, err
	}
}

type DeleteRequest struct {
	ID string
}

type DeleteResponse struct {
}

func MakeDeleteEndpoint(s TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		err := s.Delete(ctx, req.ID)
		return DeleteResponse{}, err
	}
}
