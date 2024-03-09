package main

import (
	"encoding/json"
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
		_, err := w.Write([]byte("http://localhost:8080/" + shortUrl))
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

func chooseHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		setShortURL(w, r)
	} else if r.Method == http.MethodGet {
		getOriginURL(w, r)
	}
}

func main() {
	urls = make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc("/", chooseHandle)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
