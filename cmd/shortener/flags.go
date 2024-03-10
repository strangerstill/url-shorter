package main

import (
	"flag"
	"os"
)

var flagRunAddr string
var flagBaseURL string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8888", "address and port to run server")
	flag.StringVar(&flagBaseURL, "b", ":8000", "base address and port to short URL")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		flagBaseURL = envBaseURL
	}
}
