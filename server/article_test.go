package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/msksgm/go-techblog-msksgm/mock"
	"github.com/msksgm/go-techblog-msksgm/model"
)

func Test_createArticle(t *testing.T) {
	articleStore := &mock.ArticleService{}
	userStore := &mock.UserService{}
	srv := testServer()
	srv.articleService = articleStore
	srv.userService = userStore

	token, err := generateUserToken(
		&model.User{
			ID:       1,
			Username: "username",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	input := `{
		"article": {
			"title": "title",
			"body": "body",
			"slug": "slug"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", strings.NewReader(input))
	req.Header.Add("Authorization", strings.Join([]string{"Bearer", token}, " "))
	w := httptest.NewRecorder()

	currentUser := &model.User{
		Username: "username",
		Token:    token,
	}
	userStore.GetCurrentUserFn = func() *model.User {
		return currentUser
	}

	var article model.Article
	articleStore.CreateArticleFn = func(a *model.Article) error {
		article = *a
		return nil
	}
	srv.router.ServeHTTP(w, req)
	expectedResp := articleResponse(&article)

	gotResp := M{}
	err = extractResponseArticleBody(w.Body, &gotResp)
	if err != nil {
		t.Fatal(err)
	}

	if code := w.Code; code != http.StatusCreated {
		t.Errorf("expected status code of 201, but got %d", code)
	}

	if !reflect.DeepEqual(expectedResp, gotResp) {
		t.Errorf("expected response %v, but got %v", expectedResp, gotResp)
	}
}

// func Test_listArticles(t *testing.T) {
// 	articleStore := &mock.ArticleService{}
// 	userStore := &mock.UserService{}
// 	srv := testServer()
// 	srv.articleService = articleStore
// 	srv.userService = userStore

// 	token, err := generateUserToken(
// 		&model.User{
// 			ID:       1,
// 			Username: "username",
// 		},
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	req := httptest.NewRequest(http.MethodGet, "/api/v1/articles", nil)
// 	req.Header.Add("Authorization", strings.Join([]string{"Bearer", token}, " "))
// 	w := httptest.NewRecorder()

// 	currentUser := &model.User{
// 		Username: "username",
// 		Token:    token,
// 	}
// 	userStore.GetCurrentUserFn = func() *model.User {
// 		return currentUser
// 	}

// 	articles := []*model.Article{
// 		{
// 			Title: "title1",
// 			Body:  "body1",
// 			Slug:  "slug1",
// 		},
// 		{
// 			Title: "title2",
// 			Body:  "body2",
// 			Slug:  "slug2",
// 		},
// 		{
// 			Title: "title3",
// 			Body:  "body3",
// 			Slug:  "slug3",
// 		},
// 	}
// 	articleStore.ArticlesFn = func() ([]*model.Article, error) {
// 		return articles, nil
// 	}
// 	srv.router.ServeHTTP(w, req)
// 	// expectedResp := articleResponse(&articles)

// 	// gotResp := M{}
// 	// extractResponseArticleBody(w.Body, &gotResp)

// 	if code := w.Code; code != http.StatusOK {
// 		t.Errorf("expected status code of 200, but got %d", code)
// 	}

// 	// if !reflect.DeepEqual(expectedResp, gotResp) {
// 	// 	t.Errorf("expected response %v, but got %v", expectedResp, gotResp)
// 	// }
// }

func Test_getArticle(t *testing.T) {
	articleStore := &mock.ArticleService{}
	userStore := &mock.UserService{}
	srv := testServer()
	srv.articleService = articleStore
	srv.userService = userStore

	token, err := generateUserToken(
		&model.User{
			ID:       1,
			Username: "username",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/{string}", nil)
	req.Header.Add("Authorization", strings.Join([]string{"Bearer", token}, " "))
	w := httptest.NewRecorder()

	currentUser := &model.User{
		Username: "username",
		Token:    token,
	}
	userStore.GetCurrentUserFn = func() *model.User {
		return currentUser
	}

	articles := []*model.Article{
		{
			Title: "title1",
			Body:  "body1",
			Slug:  "slug1",
		},
	}
	articleStore.ArticlesFn = func() ([]*model.Article, error) {
		return articles, nil
	}
	srv.router.ServeHTTP(w, req)
	expectedResp := articleResponse(articles[0])

	gotResp := M{}
	err = extractResponseArticleBody(w.Body, &gotResp)
	if err != nil {
		t.Fatal(err)
	}

	if code := w.Code; code != http.StatusOK {
		t.Errorf("expected status code of 200, but got %d", code)
	}

	if !reflect.DeepEqual(expectedResp, gotResp) {
		t.Errorf("expected response %v, but got %v", expectedResp, gotResp)
	}
}

// func Test_updateArticle(t *testing.T) {
// 	articleStore := &mock.ArticleService{}
// 	userStore := &mock.UserService{}
// 	srv := testServer()
// 	srv.articleService = articleStore
// 	srv.userService = userStore

// 	token, err := generateUserToken(
// 		&model.User{
// 			ID:       1,
// 			Username: "username",
// 		},
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	input := `{
// 		"article": {
// 			"title": "title_updated",
// 			"body": "body_updated",
// 			"slug": "slug_updated"
// 		}
// 	}`

// 	req := httptest.NewRequest(http.MethodPatch, "/api/v1/articles/{slug}", strings.NewReader(input))
// 	req.Header.Add("Authorization", strings.Join([]string{"Bearer", token}, " "))
// 	w := httptest.NewRecorder()

// 	currentUser := &model.User{
// 		Username: "username",
// 		Token:    token,
// 	}
// 	userStore.GetCurrentUserFn = func() *model.User {
// 		return currentUser
// 	}

// 	article := &model.Article{
// 		Title: "title",
// 		Body:  "body",
// 		Slug:  "slug",
// 	}
// 	articleStore.ArticleBySlugFn = func() (*model.Article, error) {
// 		return article, nil
// 	}
// 	var updateArticle model.Article
// 	articleStore.UpdateArticleFn = func(a *model.Article) error {
// 		updateArticle = *a
// 		return nil
// 	}
// 	srv.router.ServeHTTP(w, req)
// 	expectedResp := articleResponse(&updateArticle)

// 	gotResp := M{}
// 	extractResponseArticleBody(w.Body, &gotResp)

// 	fmt.Println(expectedResp)
// 	fmt.Println(gotResp)
// 	if code := w.Code; code != http.StatusOK {
// 		t.Errorf("expected status code of 200, but got %d", code)
// 	}

// 	if !reflect.DeepEqual(expectedResp, gotResp) {
// 		t.Errorf("expected response %v, but got %v", expectedResp, gotResp)
// 	}
// }

func Test_deleteArticle(t *testing.T) {
	articleStore := &mock.ArticleService{}
	userStore := &mock.UserService{}
	srv := testServer()
	srv.articleService = articleStore
	srv.userService = userStore

	token, err := generateUserToken(
		&model.User{
			ID:       1,
			Username: "username",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/articles/{slug}", nil)
	req.Header.Add("Authorization", strings.Join([]string{"Bearer", token}, " "))
	w := httptest.NewRecorder()

	currentUser := &model.User{
		Username: "username",
		Token:    token,
	}
	userStore.GetCurrentUserFn = func() *model.User {
		return currentUser
	}

	article := &model.Article{
		Title: "title",
		Body:  "body",
		Slug:  "slug",
	}
	articleStore.ArticleBySlugFn = func() (*model.Article, error) {
		return article, nil
	}

	articleStore.DeleteArticleFn = func() error {
		return nil
	}
	srv.router.ServeHTTP(w, req)

	if code := w.Code; code != http.StatusNoContent {
		t.Errorf("expected status code of 204, but got %d", code)
	}
}

func extractResponseArticleBody(body io.Reader, v interface{}) error {
	mm := M{}
	_ = readJSON(body, &mm)
	byt, err := json.Marshal(mm["article"])
	if err != nil {
		return fmt.Errorf("err: %v is occuered when json.Marshal()", err)
	}
	err = json.Unmarshal(byt, v)
	if err != nil {
		return fmt.Errorf("err: %v is occuered when json.Unmarshal", err)
	}
	return nil
}
