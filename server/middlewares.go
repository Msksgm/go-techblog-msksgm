package server

import (
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/msksgm/go-techblog-msksgm/model"
)

func Logger(w io.Writer) func(h http.Handler) http.Handler {
	return (func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(w, h)
	})
}

func (s *Server) authenticate(mustAuth bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				if mustAuth {
					invalidAuthTokenError(w)
				} else {
					r = setContextUser(r, &model.AnonymousUser)
					h.ServeHTTP(w, r)
				}
				return
			}

			ss := strings.Split(authHeader, " ")

			if len(ss) < 2 {
				invalidAuthTokenError(w)
				return
			}

			token := ss[1]

			claims, err := parseUserToken(token)
			if err != nil {
				invalidAuthTokenError(w)
				return
			}

			username := claims["username"].(string)
			user, err := s.userService.UserByUsername(r.Context(), username)
			if err != nil {
				serverError(w, err)
				return
			}

			r = setContextUser(r, user)
			r = setContextUserToken(r, token)
			h.ServeHTTP(w, r)
		})
	}
}
