package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/msksgm/go-techblog-msksgm/model"
)

var _ model.ArticleService = (*ArticleService)(nil)

type ArticleService struct {
	db *DB
}

func NewArticleService(db *DB) *ArticleService {
	return &ArticleService{db}
}

func (as *ArticleService) CreateArticle(ctx context.Context, article *model.Article) error {
	tx, err := as.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := createArticle(ctx, tx, article); err != nil {
		return err
	}

	return tx.Commit()
}

func createArticle(ctx context.Context, tx *sqlx.Tx, article *model.Article) error {
	query := `
	INSERT INTO articles (title, body, author_id, slug) 
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	args := []interface{}{
		article.Title,
		article.Body,
		article.AuthorID,
		article.Slug,
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&article.ID, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (as *ArticleService) Articles(ctx context.Context, filter model.ArticleFilter) ([]*model.Article, error) {
	tx, err := as.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	articles, err := findArticles(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	return articles, err
}

func findArticles(ctx context.Context, tx *sqlx.Tx, filter model.ArticleFilter) ([]*model.Article, error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.AuthorID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("author_id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Slug; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("slug = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Title; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("title = $%d", argPosition)), append(args, *v)
	}

	if v := filter.AuthorUsername; v != nil {
		argPosition++
		clause := "author_id = (select id from users where username = $%d)"
		where, args = append(where, fmt.Sprintf(clause, argPosition)), append(args, *v)
	}

	query := "SELECT * from articles" + formatWhereClause(where) + " ORDER BY created_at DESC"
	articles, err := queryArticles(ctx, tx, query, args...)
	if err != nil {
		return articles, err
	}

	return articles, nil
}

func queryArticles(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) ([]*model.Article, error) {
	articles := make([]*model.Article, 0)
	err := findMany(ctx, tx, &articles, query, args...)
	if err != nil {
		return articles, err
	}

	for _, article := range articles {
		if err := attachArticleAssociation(ctx, tx, article); err != nil {
			return nil, err
		}
	}

	return articles, nil
}

func attachArticleAssociation(ctx context.Context, tx *sqlx.Tx, article *model.Article) error {
	user, err := findUserByID(ctx, tx, article.AuthorID)
	if err != nil {
		return fmt.Errorf("cannot find article author: %w", err)
	}

	article.Author = user

	return nil
}
