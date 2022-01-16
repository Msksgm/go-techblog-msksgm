package mock

import (
	"context"

	"github.com/msksgm/go-techblog-msksgm/model"
)

type UserService struct {
	CreateUserFn     func(*model.User) error
	AuthenticateFn   func() *model.User
	GetCurrentUserFn func() *model.User
}

func (m *UserService) CreateUser(_ context.Context, user *model.User) error {
	return m.CreateUserFn(user)
}

func (m *UserService) Authenticate(_ context.Context, username, password string) (*model.User, error) {
	return m.AuthenticateFn(), nil
}

func (m *UserService) UserByUsername(_ context.Context, username string) (*model.User, error) {
	return m.GetCurrentUserFn(), nil
}
