package server

import (
	"context"
	"fmt"
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

func userFromContext(ctx context.Context) (*model.User, error) {
	user, ok := ctx.Value(userKey).(*model.User)

	if !ok {
		return user, fmt.Errorf("error is occured when ctx.Value(userKey).(*model.User)")
	}

	return user, nil
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
