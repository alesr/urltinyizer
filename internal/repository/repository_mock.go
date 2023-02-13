package repository

import "context"

var _ Repository = (*Mock)(nil)

type Mock struct {
	GetShortURLFunc  func(ctx context.Context, longURL string) (string, error)
	GetLongURLFunc   func(ctx context.Context, shortURL string) (string, error)
	GetStatsFunc     func(ctx context.Context, shortURL string) (int, error)
	SaveShortURLFunc func(ctx context.Context, longURL, shortURL string) error
}

func (m *Mock) GetShortURL(ctx context.Context, longURL string) (string, error) {
	return m.GetShortURLFunc(ctx, longURL)
}

func (m *Mock) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	return m.GetLongURLFunc(ctx, shortURL)
}

func (m *Mock) GetStats(ctx context.Context, shortURL string) (int, error) {
	return m.GetStatsFunc(ctx, shortURL)
}

func (m *Mock) SaveShortURL(ctx context.Context, longURL, shortURL string) error {
	return m.SaveShortURLFunc(ctx, longURL, shortURL)
}
