package firewall

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-toschool/palermo"
	"github.com/go-toschool/palermo/auth"
	"github.com/go-toschool/sicily"
)

const (
	// AuthUserIDContextKey key for context
	AuthUserIDContextKey sicily.StringValueKey = "userIdContextKey"
	authTokenCookieName                        = "access_token"
	tokenTypePrefix                            = "Bearer "
	tokenHeaderKey                             = "Authorization"
	tokenMetaKey                               = "auth_token"
)

// CorsAndToken ...
type CorsAndToken struct {
	SessionService auth.AuthServiceClient
}

// CheckCorsAndAuth ...
func (ac *CorsAndToken) CheckCorsAndToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}

		cred, err := parseAuthCredentials(r)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.Background()
		session, err := ac.SessionService.Get(ctx, &auth.GetRequest{
			Data: &auth.SessionCredentials{
				ValidationToken: cred.ValidationToken,
				AuthToken:       cred.AuthToken,
			},
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx1 := r.Context()
		ctx1 = setUserIDToRequestContext(ctx1, session.Data.UserId)
		next.ServeHTTP(w, r.WithContext(ctx1))
	}
}

func NewAuth(ss auth.AuthServiceClient) *CorsAndToken {
	return &CorsAndToken{
		SessionService: ss,
	}
}

func setUserIDToRequestContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, AuthUserIDContextKey, userID)
}

func parseAuthToken(r *http.Request) (string, error) {
	header := r.Header.Get(tokenHeaderKey)
	if !strings.HasPrefix(header, tokenTypePrefix) {
		return "", errors.New("auth: no token authorization header present")
	}
	return header[len(tokenTypePrefix):], nil
}

func parseAuthCredentials(r *http.Request) (*palermo.SessionCredentials, error) {
	authToken, err := parseAuthToken(r)
	if err != nil {
		return nil, err
	}
	validationToken, err := parseValidationToken(r)
	if err != nil {
		return nil, err
	}

	return &palermo.SessionCredentials{
		AuthToken:       authToken,
		ValidationToken: validationToken,
	}, nil
}

func parseValidationToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(authTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
