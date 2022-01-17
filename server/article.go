package server

import (
	"net/http"

	"github.com/msksgm/go-techblog-msksgm/model"
)

func articleResponse(article *model.Article) M {
	if article == nil {
		return nil
	}
	return M{
		"title": article.Title,
		"body":  article.Body,
		"slug":  article.Slug,
	}
}

func (s *Server) createArticle() http.HandlerFunc {
	type Input struct {
		Article struct {
			Title string `json:"title" validate:"required"`
			Body  string `json:"body" validate:"required"`
			Slug  string `json:"slug" validate:"required"`
		} `json:"article"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		if err := validate.Struct(input.Article); err != nil {
			validationError(w, err)
			return
		}

		article := model.Article{
			Title: input.Article.Title,
			Body:  input.Article.Body,
			Slug:  input.Article.Slug,
		}

		user := userFromContext(r.Context())
		article.Author = user
		article.AuthorID = user.ID

		if user.IsAnonymous() {
			invalidAuthTokenError(w)
			return
		}

		if err := s.articleService.CreateArticle(r.Context(), &article); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, M{"article": article})
	}
}

func (s *Server) listArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filter := model.ArticleFilter{}

		if v := query.Get("author"); v != "" {
			filter.AuthorUsername = &v
		}

		articles, err := s.articleService.Articles(r.Context(), filter)
		if err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"articles": articles})
	}
}
