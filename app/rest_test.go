//go:build integration
// +build integration

package app

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/alesr/urltinyizer/helper"
	"github.com/alesr/urltinyizer/internal/repository"
	"github.com/alesr/urltinyizer/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func TestCreateShortURL(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupServerHelper(t, ctx)

	t.Run("create short url", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodPost,
			"http://localhost:8080/shorten",
			strings.NewReader(`{"long_url": "https://www.google.com/"}`),
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response CreateShortURLResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, "http://foo.com/595c3c", response.ShortURL)
	})

	t.Run("failed validation", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodPost,
			"http://localhost:8080/shorten",
			strings.NewReader(`{"long_url": "invalid_url"}`),
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestRedirectToLongURL(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupServerHelper(t, ctx)

	t.Run("redirect to long url", func(t *testing.T) {
		// First create a short url

		req, err := http.NewRequest(
			http.MethodPost,
			"http://localhost:8080/shorten",
			strings.NewReader(`{"long_url": "https://www.twitter.com/"}`),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response CreateShortURLResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		defer resp.Body.Close()

		givenShortURL := url.PathEscape(response.ShortURL)

		// Then use the short url to redirect to the long url

		req, err = http.NewRequest(http.MethodGet, "http://localhost:8080/"+givenShortURL, nil)
		require.NoError(t, err)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("failed validation", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/invalid_url", nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetStats(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupServerHelper(t, ctx)

	t.Run("get stats", func(t *testing.T) {
		givenShortURL := url.PathEscape("http://shorturl/foobar")

		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/"+givenShortURL+"/stats", nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response GetStatsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, 0, response.Hits)
	})

	t.Run("get more hits", func(t *testing.T) {
		// Create short url

		req, err := http.NewRequest(
			http.MethodPost,
			"http://localhost:8080/shorten",
			strings.NewReader(`{"long_url": "https://www.google.com/"}`),
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var createShortURLResp CreateShortURLResponse
		err = json.NewDecoder(resp.Body).Decode(&createShortURLResp)
		require.NoError(t, err)

		defer resp.Body.Close()

		givenShortURL := url.PathEscape(createShortURLResp.ShortURL)

		// Fetch the short url 5 times

		for i := 0; i < 5; i++ {
			req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/"+givenShortURL, nil)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, resp.StatusCode)
		}

		req, err = http.NewRequest(http.MethodGet, "http://localhost:8080/"+givenShortURL+"/stats", nil)
		require.NoError(t, err)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response GetStatsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		defer resp.Body.Close()

		// Assert that the hits are 5
		assert.Equal(t, 5, response.Hits)
	})
}

func setupServerHelper(t *testing.T, ctx context.Context) {
	db := helper.SetupDB(t, "../migrations")

	repo := repository.NewPostgreSQL(zap.NewNop(), db)
	service := service.NewServiceDefault(zap.NewNop(), "http://foo.com/", repo)

	testApp := NewREST(zap.NewNop(), chi.NewRouter(), service)
	testApp.RegisterRoutes()

	go testApp.Run(ctx)
}
