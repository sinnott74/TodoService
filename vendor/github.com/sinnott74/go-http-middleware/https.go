package middleware

import (
	"net/http"
)

// HTTPS middleware is responsible for redirecting the user to HTTPS
// It looks at the x-forward-proto header to determine the protocol used
// x-forward-proto is commonly set when behind load balancer which will terminate the ssl connection. e.g. AWS, Cloud Foundry, etc
func HTTPS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		proto := r.Header.Get("x-forwarded-proto")
		if proto == "http" {
			http.Redirect(w, r, "https://"+r.Host+r.URL.Path, http.StatusPermanentRedirect)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
