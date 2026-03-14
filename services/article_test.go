package services_test

import (
	"context"
	"errors"
	"testing"

	"go-mongodb-api/mocks"
	"go-mongodb-api/models"
	"go-mongodb-api/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestArticleService_GetAllArticles(t *testing.T) {
	mockRepo := new(mocks.MockArticleRepository)
	svc := services.NewArticleService(mockRepo)

	expected := []models.Article{{ID: bson.NewObjectID(), Title: "Test Article"}}
	mockRepo.On("GetAll", mock.Anything, 1, 10, mock.Anything, "", "").Return(expected, int64(1), nil)

	articles, total, err := svc.GetAllArticles(context.Background(), 1, 10, map[string]string{}, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, articles, 1)
	mockRepo.AssertExpectations(t)
}

func TestArticleService_GetArticleByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockArticleRepository)
	svc := services.NewArticleService(mockRepo)

	id := bson.NewObjectID()
	expected := &models.Article{ID: id, Title: "Test Article"}
	mockRepo.On("GetByID", mock.Anything, id.Hex()).Return(expected, nil)

	article, err := svc.GetArticleByID(context.Background(), id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "Test Article", article.Title)
	mockRepo.AssertExpectations(t)
}

func TestArticleService_GetArticleByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockArticleRepository)
	svc := services.NewArticleService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "bad-id").Return(nil, errors.New("not found"))

	article, err := svc.GetArticleByID(context.Background(), "bad-id")
	assert.Error(t, err)
	assert.Nil(t, article)
	mockRepo.AssertExpectations(t)
}

func TestArticleService_CreateArticle(t *testing.T) {
	mockRepo := new(mocks.MockArticleRepository)
	svc := services.NewArticleService(mockRepo)

	article := &models.Article{Title: "New Article", Content: "Content", Slug: "new-article"}
	mockRepo.On("Create", mock.Anything, article).Return(nil)

	err := svc.CreateArticle(context.Background(), article)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestArticleService_DeleteArticle(t *testing.T) {
	mockRepo := new(mocks.MockArticleRepository)
	svc := services.NewArticleService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "article-id").Return(nil)

	err := svc.DeleteArticle(context.Background(), "article-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestArticleService_UpdateArticle_Success(t *testing.T) {
	mockRepo := new(mocks.MockArticleRepository)
	svc := services.NewArticleService(mockRepo)

	id := bson.NewObjectID()
	input := &models.Article{Title: "Updated", Content: "Content", Slug: "updated"}
	expected := &models.Article{ID: id, Title: "Updated", Content: "Content", Slug: "updated"}
	mockRepo.On("Update", mock.Anything, id.Hex(), input).Return(expected, nil)

	result, err := svc.UpdateArticle(context.Background(), id.Hex(), input)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestArticleService_UpdateArticle_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockArticleRepository)
	svc := services.NewArticleService(mockRepo)

	input := &models.Article{Title: "Updated"}
	mockRepo.On("Update", mock.Anything, "bad-id", input).Return(nil, errors.New("not found"))

	result, err := svc.UpdateArticle(context.Background(), "bad-id", input)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
