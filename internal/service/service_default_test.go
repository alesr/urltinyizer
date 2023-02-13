package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/alesr/urltinyizer/internal/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateShortURL(t *testing.T) {
	t.Parallel()

	t.Run("create short url", func(t *testing.T) {
		t.Parallel()

		given := "https://www.foo.com"
		expect := "http://bar/7633a1"

		repoMock := &repository.Mock{
			GetShortURLFunc: func(ctx context.Context, longURL string) (string, error) {
				return "", nil
			},
			SaveShortURLFunc: func(ctx context.Context, longURL, shortURL string) error {
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), "http://bar/", repoMock)

		observed, err := svc.CreateShortURL(context.Background(), given)
		require.NoError(t, err)

		require.Equal(t, expect, observed)
	})

	t.Run("get existing short url", func(t *testing.T) {
		t.Parallel()

		given := "https://www.foo.com"
		expect := "http://bar/7633a1"

		repoMock := &repository.Mock{
			GetShortURLFunc: func(ctx context.Context, longURL string) (string, error) {
				return expect, nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), "http://bar/", repoMock)

		observed, err := svc.CreateShortURL(context.Background(), given)
		require.NoError(t, err)

		require.Equal(t, expect, observed)
	})

	t.Run("error getting short url", func(t *testing.T) {
		t.Parallel()

		given := "https://www.foo.com"

		repoMock := &repository.Mock{
			GetShortURLFunc: func(ctx context.Context, longURL string) (string, error) {
				return "", fmt.Errorf("error getting short url")
			},
		}

		svc := NewServiceDefault(zap.NewNop(), "http://bar/", repoMock)

		_, err := svc.CreateShortURL(context.Background(), given)
		require.Error(t, err)
	})

	t.Run("error saving short url", func(t *testing.T) {
		t.Parallel()

		given := "https://www.foo.com"

		repoMock := &repository.Mock{
			GetShortURLFunc: func(ctx context.Context, longURL string) (string, error) {
				return "", nil
			},
			SaveShortURLFunc: func(ctx context.Context, longURL, shortURL string) error {
				return fmt.Errorf("error saving short url")
			},
		}

		svc := NewServiceDefault(zap.NewNop(), "http://bar/", repoMock)

		_, err := svc.CreateShortURL(context.Background(), given)
		require.Error(t, err)
	})
}

func TestRedirectToLongURL(t *testing.T) {
	t.Parallel()

	t.Run("redirect to long url", func(t *testing.T) {
		t.Parallel()

		given := "http://bar/7633a1"
		expect := "https://www.foo.com"

		repoMock := &repository.Mock{
			GetLongURLFunc: func(ctx context.Context, shortURL string) (string, error) {
				return expect, nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), "http://bar/", repoMock)

		observed, err := svc.RedirectToLongURL(context.Background(), given)
		require.NoError(t, err)

		require.Equal(t, expect, observed)
	})

	t.Run("error getting long url", func(t *testing.T) {
		t.Parallel()

		given := "http://bar/7633a1"

		repoMock := &repository.Mock{
			GetLongURLFunc: func(ctx context.Context, shortURL string) (string, error) {
				return "", fmt.Errorf("error getting long url")
			},
		}

		svc := NewServiceDefault(zap.NewNop(), "http://bar/", repoMock)

		_, err := svc.RedirectToLongURL(context.Background(), given)
		require.Error(t, err)
	})

	t.Run("error short url not found", func(t *testing.T) {
		t.Parallel()

		given := "http://bar/7633a1"

		repoMock := &repository.Mock{
			GetLongURLFunc: func(ctx context.Context, shortURL string) (string, error) {
				return "", nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), "http://bar/", repoMock)

		_, err := svc.RedirectToLongURL(context.Background(), given)
		require.Error(t, err)
	})
}

func TestGetStats(t *testing.T) {
	t.Parallel()

	t.Run("get stats", func(t *testing.T) {
		t.Parallel()

		given := "http://bar/7633a1"
		expect := 10

		repoMock := &repository.Mock{
			GetStatsFunc: func(ctx context.Context, shortURL string) (int, error) {
				return 10, nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), "http://bar/", repoMock)

		observed, err := svc.GetStats(context.Background(), given)
		require.NoError(t, err)

		require.Equal(t, expect, observed)
	})

	t.Run("error getting stats", func(t *testing.T) {
		t.Parallel()

		given := "http://bar/7633a1"

		repoMock := &repository.Mock{
			GetStatsFunc: func(ctx context.Context, shortURL string) (int, error) {
				return 0, fmt.Errorf("error getting stats")
			},
		}

		svc := NewServiceDefault(zap.NewNop(), "http://bar/", repoMock)

		_, err := svc.GetStats(context.Background(), given)
		require.Error(t, err)
	})

	t.Run("error short url not found", func(t *testing.T) {
		t.Parallel()

		given := "http://bar/7633a1"

		repoMock := &repository.Mock{
			GetStatsFunc: func(ctx context.Context, shortURL string) (int, error) {
				return 0, errors.New("some error")
			},
		}

		svc := NewServiceDefault(zap.NewNop(), "http://bar/", repoMock)

		_, err := svc.GetStats(context.Background(), given)
		require.Error(t, err)
	})
}
