package auth_jwt

import (
	"context"
	"errors"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

var (
	errInvalidScope = errors.New("invalid scope").Error()
)

type Authorizer struct {
	Jwks *keyfunc.JWKS
}

type MyCustomClaims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

func NewAuthorizer(jwksDomain string) (*Authorizer, error) {
	jwks, err := keyfunc.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", jwksDomain), keyfunc.Options{})
	if err != nil {
		return nil, err
	}

	return &Authorizer{
		Jwks: jwks,
	}, nil
}

func (a *Authorizer) EnsureValidToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, "Missing Authorization header", http.StatusBadRequest)
			return
		}
		accessToken := strings.TrimPrefix(authorization, "Bearer ")

		token, claims, err := a.ParseToken(accessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Authorizer) ParseToken(authorization string) (*jwt.Token, *MyCustomClaims, error) {
	if len(authorization) < 1 {
		return nil, nil, nil
	}

	accessToken := strings.TrimPrefix(authorization, "Bearer ")

	claimsStruct := MyCustomClaims{}
	token, err := jwt.ParseWithClaims(accessToken, &claimsStruct, a.Jwks.Keyfunc)
	if err != nil {
		return nil, nil, err
	}
	return token, &claimsStruct, nil
}

func (c MyCustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}
