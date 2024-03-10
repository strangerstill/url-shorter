package main

import (
	"flag"
)

var flagRunAddr string
var flagBaseShortAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8888", "address and port to run server")
	flag.StringVar(&flagBaseShortAddr, "b", ":8000", "base address and port to short URL")
	flag.Parse()
}
