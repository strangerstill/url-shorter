package app

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	RunAddr string
	BaseURL url.URL
}

func MakeConfig() (Config, error) {
	var conf Config
	flag.StringVar(&conf.RunAddr, "a", ":8080", "address and port to run server")
	rawBaseURL := flag.String("b", "http://localhost:8080/my-url", "base address and port to short URL")
	flag.Parse()

	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		conf.RunAddr = addr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		*rawBaseURL = envBaseURL
	}
	baseURL, err := url.Parse(*rawBaseURL)
	if err != nil {
		return conf, fmt.Errorf("base url parsing: %w", err)
	}
	conf.BaseURL = *baseURL
	return conf, nil
}
