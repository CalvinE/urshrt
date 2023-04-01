package main

import (
	"log"
	"os"
	"path"

	"github.com/calvine/urshrt/keygenerator"
	"github.com/calvine/urshrt/shortener"
)

func main() {
	logger := log.New(os.Stdout, "urshrt", log.LstdFlags)
	// init service
	keygen, err := keygenerator.NewRandomKeyGenerator()
	if err != nil {
		log.Panicf("failed to make key generator: %v", err)
	}
	rootPath := getEnvOrDefault("SHRTROOTPATH", path.Join(".", "shrt_test"))
	shortener, err := shortener.NewFileShortener(rootPath)
	if err != nil {
		log.Panicf("failed to make url shortener: %v", err)
	}
	// max tries is hardcoded for now...
	service, err := NewURLShortenerService(shortener, keygen, 10)
	if err != nil {
		log.Panicf("failed to make url shortener service: %v", err)
	}
	server := NewServer(service)
	// init server
	addr := getEnvOrDefault("ADDR", ":8080")
	if err := server.InitServer(logger, addr); err != nil {
		log.Printf("server stopped: %v", err)
	}
}

func getEnvOrDefault(name string, defaultValue string) string {
	value, found := os.LookupEnv(name)
	if !found {
		return defaultValue
	}
	return value
}
