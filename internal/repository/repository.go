package repository

import "context"

// Repository is an interface that defines the methods that a repository should implement.
type Repository interface {
	GetShortURL(ctx context.Context, longURL string) (string, error)
	GetLongURL(ctx context.Context, shortURL string) (string, error)
	GetStats(ctx context.Context, shortURL string) (int, error)
	SaveShortURL(ctx context.Context, longURL, shortURL string) error
}
