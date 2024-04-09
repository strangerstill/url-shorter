package app

import (
	"errors"
	"math/rand"
	"time"
)

type App struct {
	store *Storage
}

func NewApp() *App {
	return &App{NewStorage()}
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

func (a *App) SaveURL(url string) (id string, err error) {
	id = a.makeRandID()
	ok := a.store.Set(id, url)
	if ok {
		return id, nil
	}
	return id, errors.New("failed to save url")
}

func (a *App) GetURL(id string) (string, error) {
	url, ok := a.store.Get(id)
	if !ok {
		return url, ErrNotFound
	}
	return url, nil
}
