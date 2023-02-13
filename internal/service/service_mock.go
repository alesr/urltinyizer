package service

import "context"

var _ Service = (*Mock)(nil)

type Mock struct {
	CreateShortURLFunc    func(ctx context.Context, longURL string) (string, error)
	RedirectToLongURLFunc func(ctx context.Context, shortURL string) (string, error)
	GetStatsFunc          func(ctx context.Context, shortURL string) (int, error)
}

func (m *Mock) CreateShortURL(ctx context.Context, longURL string) (string, error) {
	return m.CreateShortURLFunc(ctx, longURL)
}

func (m *Mock) RedirectToLongURL(ctx context.Context, shortURL string) (string, error) {
	return m.RedirectToLongURLFunc(ctx, shortURL)
}

func (m *Mock) GetStats(ctx context.Context, shortURL string) (int, error) {
	return m.GetStatsFunc(ctx, shortURL)
}
