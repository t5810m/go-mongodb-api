package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-mongodb-api/handlers"
	"go-mongodb-api/mocks"
	"go-mongodb-api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestArticleHandler_GetAllArticles_Success(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	articles := []models.Article{{ID: bson.NewObjectID(), Title: "Test Article", Content: "Content", Slug: "test"}}
	mockSvc.On("GetAllArticles", mock.Anything, 1, 10, mock.Anything, "", "").Return(articles, int64(1), nil)

	r := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	h.GetAllArticles(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestArticleHandler_GetAllArticles_Error(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	mockSvc.On("GetAllArticles", mock.Anything, 1, 10, mock.Anything, "", "").Return([]models.Article{}, int64(0), errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	h.GetAllArticles(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestArticleHandler_GetArticleByID_Success(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	id := bson.NewObjectID()
	article := &models.Article{ID: id, Title: "Test Article"}
	mockSvc.On("GetArticleByID", mock.Anything, id.Hex()).Return(article, nil)

	r := httptest.NewRequest(http.MethodGet, "/articles/"+id.Hex(), nil)
	r = addChiURLParam(r, "id", id.Hex())
	w := httptest.NewRecorder()

	h.GetArticleByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestArticleHandler_GetArticleByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	mockSvc.On("GetArticleByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/articles/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.GetArticleByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestArticleHandler_CreateArticle_Success(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	mockSvc.On("CreateArticle", mock.Anything, mock.AnythingOfType("*models.Article")).Return(nil)

	body := `{"title":"New Article","content":"Some content here","slug":"new-article"}`
	r := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateArticle(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestArticleHandler_CreateArticle_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreateArticle(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestArticleHandler_UpdateArticle_Success(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	id := bson.NewObjectID()
	updated := &models.Article{ID: id, Title: "Updated", Content: "Content", Slug: "updated"}
	mockSvc.On("UpdateArticle", mock.Anything, id.Hex(), mock.AnythingOfType("*models.Article")).Return(updated, nil)

	body := `{"title":"Updated","content":"Content","slug":"updated"}`
	r := httptest.NewRequest(http.MethodPut, "/articles/"+id.Hex(), bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", id.Hex())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateArticle(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestArticleHandler_UpdateArticle_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	mockSvc.On("UpdateArticle", mock.Anything, "bad-id", mock.AnythingOfType("*models.Article")).Return(nil, errors.New("not found"))

	body := `{"title":"Updated"}`
	r := httptest.NewRequest(http.MethodPut, "/articles/bad-id", bytes.NewBufferString(body))
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.UpdateArticle(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestArticleHandler_CreateArticle_ValidationError(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	body := `{"title":"A"}`
	r := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateArticle(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestArticleHandler_UpdateArticle_InvalidBody(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPut, "/articles/some-id", bytes.NewBufferString("not-json"))
	r = addChiURLParam(r, "id", "some-id")
	w := httptest.NewRecorder()

	h.UpdateArticle(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestArticleHandler_DeleteArticle_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	mockSvc.On("DeleteArticle", mock.Anything, "bad-id").Return(errors.New("not found"))

	r := httptest.NewRequest(http.MethodDelete, "/articles/bad-id", nil)
	r = addChiURLParam(r, "id", "bad-id")
	w := httptest.NewRecorder()

	h.DeleteArticle(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestArticleHandler_DeleteArticle_Success(t *testing.T) {
	mockSvc := new(mocks.MockArticleService)
	h := handlers.NewArticleHandler(mockSvc)

	mockSvc.On("DeleteArticle", mock.Anything, "article-id").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/articles/article-id", nil)
	r = addChiURLParam(r, "id", "article-id")
	w := httptest.NewRecorder()

	h.DeleteArticle(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}
