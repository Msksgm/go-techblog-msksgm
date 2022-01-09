package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_healthcheck(t *testing.T) {
	srv := testServer()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	srv.router.ServeHTTP(w, req)
	expectResp := M{
		"message": "health",
		"status":  "available",
	}
	gotResp := M{}
	extractResponseBody(w.Body, &gotResp)

	if code := w.Code; code != http.StatusOK {
		t.Errorf("expected status code of 200, but got %d", code)
	}

	if !reflect.DeepEqual(expectResp, gotResp) {
		t.Errorf("expected response %v, but got %v", expectResp, gotResp)
	}
}

func extractResponseBody(body io.Reader, v interface{}) {
	mm := M{}
	_ = readJSON(body, &mm)
	byt, err := json.Marshal(mm)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byt, v)
	if err != nil {
		panic(err)
	}
}
