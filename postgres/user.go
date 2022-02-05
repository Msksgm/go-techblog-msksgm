package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/msksgm/go-techblog-msksgm/model"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db}
}

func (us *UserService) CreateUser(ctx context.Context, user *model.User) error {
	tx, err := us.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := createUser(ctx, tx, user); err != nil {
		return err
	}
	return tx.Commit()
}

func createUser(ctx context.Context, tx *sqlx.Tx, user *model.User) error {
	query := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2) RETURNING id, created_at, updated_at
	`
	args := []interface{}{user.Username, user.PasswordHash}
	err := tx.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return model.ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (us *UserService) Authenticate(ctx context.Context, username, password string) (*model.User, error) {
	user, err := us.UserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if !user.VerifyPassword(password) {
		return nil, model.ErrUnAuthorized
	}

	return user, nil
}

func (us *UserService) UserByUsername(ctx context.Context, username string) (*model.User, error) {
	tx, err := us.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	tx.Rollback()

	user, err := findOneUser(ctx, tx, model.UserFilter{Username: &username})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func findOneUser(ctx context.Context, tx *sqlx.Tx, filter model.UserFilter) (*model.User, error) {
	users, err := findUsers(ctx, tx, filter)

	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, model.ErrNotFound
	}

	return users[0], nil
}

func findUsers(ctx context.Context, tx *sqlx.Tx, filter model.UserFilter) ([]*model.User, error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Username; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("username = $%d", argPosition)), append(args, *v)
	}

	query := "SELECT * from users" + formatWhereClause(where) + " ORDER BY id ASC" + formatLimitOffset(filter.Limit, filter.Offset)

	users, err := queryUsers(ctx, tx, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func queryUsers(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) ([]*model.User, error) {
	users := make([]*model.User, 0)

	if err := findMany(ctx, tx, &users, query, args...); err != nil {
		return users, err
	}

	return users, nil
}

func (us *UserService) UpdateUser(ctx context.Context, user *model.User, patch model.UserPatch) error {
	tx, err := us.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println(err)
		return model.ErrInternal
	}

	tx.Rollback()

	if err := updateUser(ctx, tx, user, patch); err != nil {
		log.Println(err)
		return model.ErrInternal
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return model.ErrInternal
	}

	return nil
}

func updateUser(ctx context.Context, tx *sqlx.Tx, user *model.User, patch model.UserPatch) error {
	if v := patch.Username; v != nil {
		user.Username = *v
	}

	if v := patch.PasswordHash; v != nil {
		user.PasswordHash = *v
	}

	args := []interface{}{
		user.Username,
		user.PasswordHash,
		user.ID,
	}

	query := `
	UPDATE users
	SET username = $1, password_hash=$2, updated_at=NOW()
	WHERE id = $3
	RETURNING updated_at`

	if err := tx.QueryRowxContext(ctx, query, args...).Scan(&user.UpdatedAt); err != nil {
		log.Printf("error updating record: %v", err)
		return model.ErrInternal
	}

	return nil
}

func findUserByID(ctx context.Context, tx *sqlx.Tx, id uint) (*model.User, error) {
	return findOneUser(ctx, tx, model.UserFilter{ID: &id})
}
