package main

import (
	"context"
	"fmt"
	"log"

	"github.com/calvine/urshrt/keygenerator"
	"github.com/calvine/urshrt/shortener"
)

type URLShortenerService interface {
	Shorten(ctx context.Context, logger *log.Logger, url string, keyLength int) (string, error)
	Embiggen(ctx context.Context, logger *log.Logger, key string) (string, error)
}

type urlShortenerService struct {
	maxTries     int
	shortener    shortener.Shortener
	keyGenerator keygenerator.KeyGenerator
}

func (u *urlShortenerService) Shorten(ctx context.Context, logger *log.Logger, url string, keyLength int) (string, error) {
	var exists = true
	var tries = 0
	var err error
	var shrtKey string
	for exists && tries < u.maxTries {
		tries++
		shrtKey, err = u.keyGenerator.GenerateKey(keyLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate short key for url of length %d: %w", keyLength, err)
		}
		exists, err = u.shortener.DoesKeyExist(ctx, logger, shrtKey)
		if err != nil {
			return "", fmt.Errorf("failed to check if short key for url exists %s: %w", shrtKey, err)
		}
	}
	if exists {
		return "", fmt.Errorf("failed to generate unique key of length %d after %d tries", keyLength, tries)
	}
	// at this point we should have a unique short key for the URL
	err = u.shortener.Shorten(ctx, logger, url, shrtKey)
	if err != nil {
		return "", fmt.Errorf("failed to save shortened url key: %w", err)
	}
	return shrtKey, nil
}

func (u *urlShortenerService) Embiggen(ctx context.Context, logger *log.Logger, key string) (string, error) {
	principal, err := u.shortener.Embiggen(ctx, logger, key)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve principal for key %s: %w", key, err)
	}
	return principal, nil
}

func NewURLShortenerService(shortener shortener.Shortener, keyGenerator keygenerator.KeyGenerator, maxTries int) (URLShortenerService, error) {
	return &urlShortenerService{
		keyGenerator: keyGenerator,
		maxTries:     maxTries,
		shortener:    shortener,
	}, nil
}
