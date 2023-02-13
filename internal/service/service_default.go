package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/alesr/urltinyizer/internal/repository"
	"go.uber.org/zap"
)

var _ Service = (*ServiceDefault)(nil)

type ServiceDefault struct {
	logger  *zap.Logger
	appHost string
	repo    repository.Repository
}

func NewServiceDefault(logger *zap.Logger, appHost string, repo repository.Repository) *ServiceDefault {
	return &ServiceDefault{
		logger:  logger,
		appHost: appHost,
		repo:    repo,
	}
}

func (s *ServiceDefault) CreateShortURL(ctx context.Context, longURL string) (string, error) {
	existingShortURL, err := s.repo.GetShortURL(ctx, longURL)
	if err != nil {
		return "", fmt.Errorf("could not get short url: %w", err)
	}

	if existingShortURL != "" {
		return existingShortURL, nil
	}

	shortURL, err := s.generateShortURL(longURL)
	if err != nil {
		return "", fmt.Errorf("could not generate short url: %w", err)
	}

	s.logger.Info("generated short url", zap.String("short_url", shortURL))

	if err := s.repo.SaveShortURL(ctx, shortURL, longURL); err != nil {
		return "", fmt.Errorf("could not save short url: %w", err)
	}
	return shortURL, nil
}

func (s *ServiceDefault) RedirectToLongURL(ctx context.Context, shortURL string) (string, error) {
	longURL, err := s.repo.GetLongURL(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("could not get long url: %w", err)
	}

	if longURL == "" {
		return "", fmt.Errorf("could not find long url for short url %s", shortURL)
	}
	return longURL, nil
}

func (s *ServiceDefault) GetStats(ctx context.Context, shortURL string) (int, error) {
	stats, err := s.repo.GetStats(ctx, shortURL)
	if err != nil {
		return 0, fmt.Errorf("could not get stats: %w", err)
	}
	return stats, nil
}

func (s *ServiceDefault) generateShortURL(longURL string) (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, longURL); err != nil {
		return "", fmt.Errorf("could not generate short url: %w", err)
	}
	return s.appHost + fmt.Sprintf("%x", h.Sum(nil))[:6], nil
}
