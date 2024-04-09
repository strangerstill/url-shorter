package app

import (
	"errors"
	"github.com/strangerstill/url-shorter/internal/models"
	"math/rand"
	"net/url"
	"time"
)

type Storager interface {
	Set(id, url string) bool
	Get(id string) (string, bool)
}

type App struct {
	store   Storager
	baseURL url.URL
}

func NewApp(baseURL url.URL) *App {
	return &App{NewStorage(), baseURL}
}

var ErrNotFound = errors.New("url was not found")

func (*App) makeRandID() string {
	const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const urlLength = 8
	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, urlLength)
	for i := range shortKey {
		shortKey[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(shortKey)
}

func (a *App) SaveURL(url string) (result models.ResultUrl, err error) {
	id := a.makeRandID()
	ok := a.store.Set(id, url)
	if ok {
		newUrl := a.baseURL
		newUrl.Path = id
		return models.ResultUrl{Result: newUrl.String()}, nil
	}
	return models.ResultUrl{Result: ""}, errors.New("failed to save url")
}

func (a *App) GetURL(id string) (string, error) {
	url, ok := a.store.Get(id)
	if !ok {
		return url, ErrNotFound
	}
	return url, nil
}
