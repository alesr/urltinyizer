package service

import "context"

// Service is an interface that defines the methods that a service should implement.
type Service interface {
	CreateShortURL(ctx context.Context, longURL string) (string, error)
	RedirectToLongURL(ctx context.Context, shortURL string) (string, error)
	GetStats(ctx context.Context, shortURL string) (int, error)
}
