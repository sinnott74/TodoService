package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// JWTFunc defines a user supplied authorisation function.
// The func is given the current context and a valid MapClaims
// This is the point at which the user can do further validation / authorisation on the claims.JWTFunc
// The context returned will be used at the context for further chained http handlers.
// JWT authorisation fails if this returns an error, and further chained http handlers are not called.
type JWTFunc func(context.Context, jwt.MapClaims) (context.Context, error)

// TokenExtractor is a function that takes the authorization header value as input and returns
// either a token or an error.  An error should only be returned if an attempt
// to specify a token was found, but the information was somehow incorrectly
// formed.  In the case where a token is simply not present, this should not
// be treated as an error.  An empty string should be returned in that case.
type TokenExtractor func(authHeaderValue string) (string, error)

// defaultTokenExtractor is the default token extractor. It recieves the Authorisation Header value. It expects it to container a value in format of JWT {token}
func defaultTokenExtractor(authHeaderValue string) (string, error) {
	authHeaderParts := strings.Split(authHeaderValue, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "jwt" {
		return "", errors.New("Authorization header format must be JWT {token}")
	}
	return authHeaderParts[1], nil
}

// JWTOptions defines the user supplied JWT configuration options.
type JWTOptions struct {
	Secret   []byte
	AuthFunc JWTFunc
	// A function that extracts the token from the request
	// Default: FromAuthHeader (i.e., from Authorization header as bearer token)
	Extractor TokenExtractor
}

// JWT is middleware which handles authentication for JsonWebTokens
func JWT(options JWTOptions) func(next http.Handler) http.Handler {

	if options.Extractor == nil {
		options.Extractor = defaultTokenExtractor
	}

	return func(next http.Handler) http.Handler {
		authenticater := jwtAuth{
			secret:           options.Secret,
			userSuppliedFunc: options.AuthFunc,
			tokenExtractor:   options.Extractor,
		}

		return Auth(authenticater.authenticate)(next)
	}
}

// jwtAuth is the private version of JWTOptions which contains the authentication function passed to Auth middleware
type jwtAuth struct {
	secret           []byte
	userSuppliedFunc JWTFunc
	tokenExtractor   TokenExtractor
}

func (auth jwtAuth) authenticate(ctx context.Context, authHeaderValue string) (context.Context, error) {

	tokenString, err := auth.tokenExtractor(authHeaderValue)
	if err != nil {
		return ctx, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.secret, nil
	})
	if err != nil {
		return ctx, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// fmt.Printf("%+v\n", token)
		// fmt.Printf("%+v\n", claims)
		if auth.userSuppliedFunc != nil {
			return auth.userSuppliedFunc(ctx, claims)
		}
		return ctx, nil
	}

	// fmt.Println(err)
	return ctx, err
}
