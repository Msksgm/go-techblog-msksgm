package postgres

import (
	"context"

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
