package server

import (
	"context"
	"net/http"

	"github.com/msksgm/go-techblog-msksgm/model"
)

type contextKey string

const (
	userKey  contextKey = "user"
	tokenKey contextKey = "token"
)

func setContextUser(r *http.Request, u *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userKey, u)
	return r.WithContext(ctx)
}

func userFromContext(ctx context.Context) *model.User {
	user, ok := ctx.Value(userKey).(*model.User)

	if !ok {
		panic("missing user context key")
	}

	return user
}

func setContextUserToken(r *http.Request, token string) *http.Request {
	ctx := context.WithValue(r.Context(), tokenKey, token)
	return r.WithContext(ctx)
}

func userTokenFromContext(ctx context.Context) string {
	token, ok := ctx.Value(tokenKey).(string)

	if !ok {
		return ""
	}

	return token
}
