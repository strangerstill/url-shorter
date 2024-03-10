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
			name:      "not empty originURL",
			originURL: "https://ieftimov.com/posts/testing-in-go-go-test/",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	parseFlags()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			originURL := strings.NewReader(test.originURL)
			postRequest := httptest.NewRequest(http.MethodPost, "/", originURL)
			postW := httptest.NewRecorder()
			setShortURL(postW, postRequest)
			postRes := postW.Result()
			resShortURL, _ := io.ReadAll(postRes.Body)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//u, _ := url.Parse(string(resShortURL))
			//shortKey := u.Path[1:]
			//originUrl, ok := storage.getURL(shortKey)

			getRequest := httptest.NewRequest(http.MethodGet, string(resShortURL), nil)
			getW := httptest.NewRecorder()

			getOriginURL(getW, getRequest)
			getRes := getW.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println(err)
					return
				}
			}(getRes.Body)

			fmt.Printf("expected code: %d, status code %d\n", test.want.code, getRes.StatusCode)
			assert.Equal(t, test.want.code, getRes.StatusCode)
			if getRes.StatusCode == http.StatusTemporaryRedirect {
				fmt.Printf("want location = %s location %s\n", test.originURL, getRes.Header.Get("Location"))
				assert.Equal(t, getRes.Header.Get("Location"), test.originURL)
			}
		})
	}

}
