package storage

import (
	"testing"
	"time"

	"github.com/drerr0r/url-shortener/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMockStorage_SaveAndGetURL(t *testing.T) {
	storage := NewMockStorage()

	url := &models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "test123",
		CreatedAt:   time.Now(),
	}

	err := storage.SaveURL(url)
	assert.NoError(t, err)

	count := storage.GetURLCount()
	assert.Equal(t, 1, count)

	retrievedURL, err := storage.GetURL("test123")
	assert.NoError(t, err)
	assert.Equal(t, url.OriginalURL, retrievedURL.OriginalURL)
	assert.Equal(t, url.ShortCode, retrievedURL.ShortCode)

	_, err = storage.GetURL("nonexistent")
	assert.Equal(t, ErrNotFound, err)
}

func TestMockStorage_URLExists(t *testing.T) {
	storage := NewMockStorage()

	url := &models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "test123",
	}

	err := storage.SaveURL(url)
	assert.NoError(t, err)

	count := storage.GetURLCount()
	assert.Equal(t, 1, count)

	exists, err := storage.URLExists("test123")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = storage.URLExists("nonexistent")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestMockStorage_DeleteURL(t *testing.T) {
	storage := NewMockStorage()

	url := &models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "test123",
	}

	err := storage.SaveURL(url)
	assert.NoError(t, err)

	count := storage.GetURLCount()
	assert.Equal(t, 1, count)

	err = storage.DeleteURL("test123")
	assert.NoError(t, err)

	_, err = storage.GetURL("test123")
	assert.Equal(t, ErrNotFound, err)

	count = storage.GetURLCount()
	assert.Equal(t, 0, count)
}

func TestMockStorage_GetURLByOriginal(t *testing.T) {
	storage := NewMockStorage()

	url := &models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "test123",
	}

	err := storage.SaveURL(url)
	assert.NoError(t, err)

	retrievedURL, err := storage.GetURLByOriginal("https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedURL)
	assert.Equal(t, url.ShortCode, retrievedURL.ShortCode)

	retrievedURL, err = storage.GetURLByOriginal("https://nonexistent.com")
	assert.NoError(t, err)
	assert.Nil(t, retrievedURL)
}

func TestMockStorage_GetURLs(t *testing.T) {
	storage := NewMockStorage()

	urls := []*models.URL{
		{OriginalURL: "https://example1.com", ShortCode: "test1"},
		{OriginalURL: "https://example2.com", ShortCode: "test2"},
		{OriginalURL: "https://example3.com", ShortCode: "test3"},
	}

	for _, url := range urls {
		err := storage.SaveURL(url)
		assert.NoError(t, err)
	}

	retrievedURLs, err := storage.GetURLs(2, 0)
	assert.NoError(t, err)
	assert.Len(t, retrievedURLs, 2)

	count, err := storage.GetURLsCount()
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}
