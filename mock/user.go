package mock

import (
	"context"

	"github.com/msksgm/go-techblog-msksgm/model"
)

type UserService struct {
	CreateUserFn func(*model.User) error
}

func (m *UserService) CreateUser(_ context.Context, user *model.User) error {
	return m.CreateUserFn(user)
}
