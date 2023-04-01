package shortener

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
)

type fileShortener struct {
	rootDir string
}

func (f *fileShortener) Shorten(ctx context.Context, logger *log.Logger, url, key string) error {
	// Not sure, but thinking about moving exists check into this function rather than it living at the service level...
	// might be over thinking it...
	keyPath := path.Join(f.rootDir, key)
	// might want to re do this with a file handler...
	err := os.WriteFile(keyPath, []byte(url), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", keyPath, err)
	}
	return nil
}

func (f *fileShortener) Embiggen(ctx context.Context, logger *log.Logger, key string) (string, error) {
	keyPath := path.Join(f.rootDir, key)
	url, err := os.ReadFile(keyPath)
	if err != nil {
		logger.Printf("failed to read file: %s - %v", keyPath, err)
		return "", err
	}
	return string(url), nil
}

func (f *fileShortener) DoesKeyExist(ctx context.Context, logger *log.Logger, key string) (bool, error) {
	keyPath := path.Join(f.rootDir, key)
	_, err := os.Stat(keyPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		logger.Printf("failed to get file info: %s - %v", keyPath, err)
		return false, err
	}
	return true, nil
}

func NewFileShortener(rootDir string) (Shortener, error) {
	s, err := os.Stat(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to check stat on root dir: %v", err)
	}
	if !s.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", rootDir)
	}
	return &fileShortener{
		rootDir: rootDir,
	}, nil
}
