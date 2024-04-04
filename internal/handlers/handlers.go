package handlers

import (
	"encoding/json"
	"errors"
	"github.com/strangerstill/url-shorter/internal/app"
	"github.com/strangerstill/url-shorter/internal/loggger"
	"github.com/strangerstill/url-shorter/internal/models"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	myApp *app.App
}

func NewHandlers(baseURL url.URL) *Handlers {
	return &Handlers{app.NewApp(baseURL)}
}

func MakeRouter(h *Handlers, loggerItem *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(loggger.LoggingMiddleware(loggerItem))
	r.Use(ZipMiddleware)
	r.Post("/", h.SaveURL)
	r.Get("/{url}", h.GetURL)
	r.Post("/api/shorten", h.SaveURLJSON)
	return r
}

func (h *Handlers) SaveURLJSON(w http.ResponseWriter, r *http.Request) {
	var payload models.PayloadUrl
	data := json.NewDecoder(r.Body)
	if err := data.Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := payload.Url
	if url == "" {
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}
	newUrl, err := h.myApp.SaveURL(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(newUrl); err != nil {
		return
	}
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

	newUrl, err := h.myApp.SaveURL(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = io.WriteString(w, newUrl.Result)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h *Handlers) GetURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	url, err := h.myApp.GetURL(id)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
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
