package shortener

import (
	"context"
	"log"
)

// Shortener is designed to be the interface to the backing store of shortened URLs
// not the generator of the short url keys.
type Shortener interface {
	// DoesKeyExist checks to see if the url short key already exists
	DoesKeyExist(ctx context.Context, logger *log.Logger, key string) (bool, error)
	// Shorten takes a url and returns the short key to retrieve it.
	Shorten(ctx context.Context, logger *log.Logger, url, key string) error
	// Embiggen takes the short url key and returns the original url
	Embiggen(ctx context.Context, logger *log.Logger, key string) (string, error)
}
