package server

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/msksgm/go-techblog-msksgm/model"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fid reflect.StructField) string {
		name := strings.SplitN(fid.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			name = ""
		}
		return name
	})
}

func userResponse(user *model.User, _token ...string) M {
	if user == nil {
		return nil
	}
	return M{
		"username": user.Username,
	}
}

func userTokenResponse(user *model.User, _token ...string) M {
	if user == nil {
		return nil
	}
	return M{
		"username": user.Username,
		"token":    user.Token,
	}
}

func (s *Server) createUser() http.HandlerFunc {
	type Input struct {
		User struct {
			Username string `json:"username" validate:"required,min=2"`
			Password string `json:"password" validate:"required,min=8,max=72"`
		} `json:"user" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		input := &Input{}

		if err := readJSON(r.Body, &input); err != nil {
			errorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, err)
			return
		}

		user := model.User{
			Username: input.User.Username,
		}

		err := user.SetPassword(input.User.Password)
		if err != nil {
			log.Fatalf("err is occured: %v", err)
		}

		if err := s.userService.CreateUser(r.Context(), &user); err != nil {
			switch {
			case errors.Is(err, model.ErrDuplicateUsername):
				err = ErrorM{"username": []string{"this username is already in use"}}
				errorResponse(w, http.StatusConflict, err)
			default:
				serverError(w, err)
			}
			return
		}

		writeJSON(w, http.StatusCreated, M{"user": user})
	}
}

func (s *Server) loginUser() http.HandlerFunc {
	type Input struct {
		User struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			errorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		user, err := s.userService.Authenticate(r.Context(), input.User.Username, input.User.Password)

		if err != nil || user == nil {
			invalidUserCredentialsError(w)
			return
		}

		token, err := generateUserToken(user)
		if err != nil {
			serverError(w, err)
			return
		}

		user.Token = token
		writeJSON(w, http.StatusOK, M{"user": user})
	}
}

func (s *Server) getCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, err := userFromContext(r.Context())
		if err != nil {
			log.Fatal(err)
		}
		user.Token = userTokenFromContext(ctx)

		writeJSON(w, http.StatusOK, M{"user": user})
	}
}

func (s *Server) updateUser() http.HandlerFunc {
	type Input struct {
		User struct {
			Username *string `json:"username,omitempty"`
			Password *string `json:"password,omitempty"`
		} `json:"user,omitempty" validate:"required"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		input := &Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, err)
			return
		}

		ctx := r.Context()
		user, err := userFromContext(r.Context())
		if err != nil {
			log.Fatal(err)
		}
		patch := model.UserPatch{
			Username: input.User.Username,
		}

		if v := input.User.Password; v != nil {
			err := user.SetPassword(*v)
			if err != nil {
				log.Fatalf("err is occured: %v", err)
			}
		}

		err = s.userService.UpdateUser(ctx, user, patch)
		if err != nil {
			serverError(w, err)
			return
		}

		user.Token = userTokenFromContext(ctx)

		writeJSON(w, http.StatusOK, M{"user": user})
	}
}
