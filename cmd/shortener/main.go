package main

import (
	"fmt"
	"net/http"

	"github.com/strangerstill/url-shorter/internal/app"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		panic(err)
	}
}

func run() error {
	conf, err := app.MakeConfig()

	if err != nil {
		return err
	}
	return http.ListenAndServe(
		conf.RunAddr,
		app.MakeRouter(app.NewHandlers(conf.BaseURL)),
	)
}
