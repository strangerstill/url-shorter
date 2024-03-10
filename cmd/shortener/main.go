package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var urls map[string]string

func makeShortUrl() string {
	const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const urlLength = 8
	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, urlLength)
	for i := range shortKey {
		shortKey[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(shortKey)
}

func setShortURL(w http.ResponseWriter, r *http.Request) {
	originURL, _ := io.ReadAll(r.Body)
	if string(originURL) != "" {
		shortUrl := makeShortUrl()
		urls[shortUrl] = string(originURL)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte("http://localhost" + flagBaseShortAddr + "/" + shortUrl))
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("empty request"))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func getOriginURL(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.String()[1:]
	url, ok := urls[shortURL]
	if !ok {
		http.Error(w, "origin url not found", http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	parseFlags()

	if err := run(); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		panic(err)
	}
}

func run() error {
	urls = make(map[string]string)

	r := chi.NewRouter()

	r.Post("/", setShortURL)
	r.Get("/{url}", getOriginURL)

	fmt.Println("flagRunAddr = ", flagRunAddr)

	fmt.Println("Running server on", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, r)
}
