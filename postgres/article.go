package postgres

import (
	"context"
	"fmt"
	"log"

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

func (as *ArticleService) ArticleBySlug(ctx context.Context, slug string) (*model.Article, error) {
	tx, err := as.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	article, err := findArticleBySlug(ctx, tx, slug)
	if err != nil {
		return nil, err
	}

	return article, tx.Commit()
}

func findArticleBySlug(ctx context.Context, tx *sqlx.Tx, slug string) (*model.Article, error) {
	filter := model.ArticleFilter{Slug: &slug}
	articles, err := findArticles(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	if len(articles) == 0 {
		return nil, model.ErrNotFound
	}

	return articles[0], err
}

func (as *ArticleService) UpdateArticle(ctx context.Context, article *model.Article, filter model.ArticlePatch) error {
	tx, err := as.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = updateArticle(ctx, tx, article, filter)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func updateArticle(ctx context.Context, tx *sqlx.Tx, article *model.Article, patch model.ArticlePatch) error {
	if v := patch.Body; v != nil {
		article.Body = *v
	}

	if v := patch.Title; v != nil {
		article.Title = *v
	}

	args := []interface{}{
		article.Body,
		article.Title,
		article.ID,
	}

	query := `
	UPDATE articles
	SET body = $1, title = $2, updated_at = NOW() WHERE id = $3
	RETURNING updated_at`

	if err := tx.QueryRowxContext(ctx, query, args...).Scan(&article.UpdatedAt); err != nil {
		log.Printf("error updating record: %v", err)
		return model.ErrInternal
	}

	return nil
}

func (as *ArticleService) DeleteArticle(ctx context.Context, id uint) error {
	tx, err := as.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = deleteArticle(ctx, tx, id)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func deleteArticle(ctx context.Context, tx *sqlx.Tx, id uint) error {
	query := "DELETE FROM articles WHERE id = $1"

	return execQuery(ctx, tx, query, id)
}
