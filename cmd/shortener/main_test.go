package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetShortURL(t *testing.T) {
	urls = make(map[string]string)
	type want struct {
		code int
	}
	tests := []struct {
		name      string
		originURL string
		want      want
	}{
		{
			name:      "empty originURL",
			originURL: "",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:      "valid originURL",
			originURL: "https://pkg.go.dev/testing",
			want: want{
				code: http.StatusCreated,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\nrun test: %v body: %v\n", test.name, test.originURL)
			originURL := strings.NewReader(test.originURL)
			request := httptest.NewRequest(http.MethodPost, "/", originURL)
			w := httptest.NewRecorder()
			setShortURL(w, request)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println(err)
					return
				}
			}(res.Body)
			fmt.Printf("expected code: %d, status code: %d\n", test.want.code, res.StatusCode)
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}

}

func TestGetOriginURL(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name      string
		originURL string
		want      want
	}{
		{
			name:      "empty originURL",
			originURL: "",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:      "urls include originURL",
			originURL: "https://pkg.go.dev/testing",
			want: want{
				code: http.StatusTemporaryRedirect,
			},
		},
		{
			name:      "urls don't include originURL",
			originURL: "https://ieftimov.com/posts/testing-in-go-go-test/",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var shortKet string
			if test.originURL == "" {
				shortKet = ""
			} else {
				for key, value := range urls {
					if value == test.originURL {
						shortKet = key
					}
				}
			}
			fmt.Printf("\nrun test: %v, origin: %v, short: %v\n", test.name, test.originURL, shortKet)
			request := httptest.NewRequest(http.MethodGet, "/"+shortKet, nil)
			w := httptest.NewRecorder()
			getOriginURL(w, request)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println(err)
					return
				}
			}(res.Body)
			fmt.Printf("expected code: %d, status code %d\n", test.want.code, res.StatusCode)
			assert.Equal(t, test.want.code, res.StatusCode)
			if res.StatusCode == http.StatusTemporaryRedirect {
				fmt.Printf("want location = %s location %s\n", test.originURL, res.Header.Get("Location"))
				assert.Equal(t, res.Header.Get("Location"), test.originURL)
			}
		})
	}

}
