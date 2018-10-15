package middleware

import (
	"context"
	"net/http"
)

// AuthFunc defines the user supplied function to implement Authorisation
// It is given the current request context and the Authorization header value
// and returns the context object to use with further chained http handlers.
// If an err is returned chained http handlers are not called
type AuthFunc func(context.Context, string) (context.Context, error)

// Auth middleware is responsible handling request authentication
// The authentication is handled by the supplied AuthFunc
func Auth(authFunc AuthFunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				// missing header
				w.WriteHeader(http.StatusUnauthorized)
				// w.Write(errors.New("unauthorized: no authentication provided").Error())
				return
			}
			ctx, err := authFunc(r.Context(), auth)
			if err != nil {
				// unauthorised
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
