package app

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	app     *App
	baseURL url.URL
}

func NewHandlers(baseURL url.URL) *Handlers {
	return &Handlers{NewApp(), baseURL}
}

func MakeRouter(h *Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", h.SaveURL)
	r.Get("/{url}", h.GetURL)
	return r
}

func (h *Handlers) SaveURL(w http.ResponseWriter, r *http.Request) {
	urlRaw, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url := string(urlRaw)
	if url == "" {
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}

	id, err := h.app.SaveURL(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newUrl := h.baseURL
	newUrl.Path = id

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(newUrl.String()))
	if err != nil {
		log.Println(err)
		return
	}
}

func (h *Handlers) GetURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	url, err := h.app.GetURL(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err = w.Write([]byte(url))
	if err != nil {
		log.Println(err)
	}
}
