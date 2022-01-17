package mock

import (
	"context"

	"github.com/msksgm/go-techblog-msksgm/model"
)

type ArticleService struct {
	CreateArticleFn func(*model.Article) error
}

func (m *ArticleService) CreateArticle(_ context.Context, article *model.Article) error {
	return m.CreateArticleFn(article)
}
