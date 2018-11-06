package todo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	httptransport "github.com/go-kit/kit/transport/http"
	middleware "github.com/sinnott74/go-http-middleware"
)

// ErrMissingParam is thrown when an http request is missing a URL Parameter
var ErrMissingParam = errors.New("Missing parameter")

// MakeHTTPHandler creates http transport layer for the Todo service
func MakeHTTPHandler(endpoints TodoEndpoints) http.Handler {

	options := []httptransport.ServerOption{
		// httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	jwtOptions := middleware.JWTOptions{
		Secret: JWTSecret(),
		AuthFunc: func(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
			// verify claims
			if username, ok := claims["username"]; ok {
				userCtx := context.WithValue(ctx, "username", username)
				return userCtx, nil
			}
			return ctx, errors.New("No username")
		},
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.StripSlashes)

	todoRouter := chi.NewRouter()
	todoRouter.Use(middleware.JWT(jwtOptions))
	todoRouter.Use(middleware.DefaultEtag)
	todoRouter.Use(chiMiddleware.DefaultCompress)

	todoRouter.Get("/", httptransport.NewServer(
		endpoints.GetAllForUserEndPoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Get("/{id}", httptransport.NewServer(
		endpoints.GetByIDEndpoint,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Post("/", httptransport.NewServer(
		endpoints.AddEndpoint,
		decodeAddRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Put("/{id}", httptransport.NewServer(
		endpoints.UpdateEndpoint,
		decodeUpdateRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Delete("/{id}", httptransport.NewServer(
		endpoints.DeleteEndpoint,
		decodeDeleteRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Mount("/api/todos", todoRouter)

	r.Get("/health", httptransport.NewServer(
		endpoints.HealthEndpoint,
		decodeHealthRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return GetAllForUserRequest{}, err
}

func decodeGetByIDRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrMissingParam
	}
	return GetByIDRequest{id}, err
}

func decodeAddRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var todo Todo
	err = render.Decode(r, &todo)
	if err != nil {
		return nil, err
	}
	return AddRequest{todo}, err
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrMissingParam
	}
	var todo Todo
	err = render.Decode(r, &todo)
	if err != nil {
		return nil, err
	}
	return UpdateRequest{id, todo}, err
}

func decodeDeleteRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrMissingParam
	}
	return DeleteRequest{id}, err
}

func decodeHealthRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return HealthRequest{}, err
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(error); ok && err != nil {
		encodeError(ctx, err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInconsistentIDs, ErrMissingParam:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
