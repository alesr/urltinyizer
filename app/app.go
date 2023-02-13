package app

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	// maxLongURLSize is the maximum size of a long URL (2MB)
	maxLongURLSize = 2048 * 1024
)

// App is an interface that defines the methods that an app should implement.
type App interface {
	Run() error
}

type CreateShortURLRequest struct {
	LongURL string `json:"long_url"`
}

func (r *CreateShortURLRequest) Validate() error {
	return validateURL(r.LongURL)
}

type CreateShortURLResponse struct {
	ShortURL string `json:"short_url"`
}

type RedirectToLongURLRequest string

func (r *RedirectToLongURLRequest) Validate() error {
	return validateURL(string(*r))
}

type GetStatsRequest string

func (r *GetStatsRequest) Validate() error {
	return validateURL(string(*r))
}

type GetStatsResponse struct {
	ShortURL string `json:"short_url"`
	Hits     int    `json:"hits"`
}

func validateURL(u string) error {
	if len(u) == 0 {
		return errors.New("url is required")
	}

	if _, err := url.Parse(u); err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	if u[:4] != "http" && u[:5] != "https" {
		return fmt.Errorf("invalid url: %w", errors.New("schema is required"))
	}

	if len(u) > maxLongURLSize {
		return fmt.Errorf("url is too large")
	}
	return nil
}
