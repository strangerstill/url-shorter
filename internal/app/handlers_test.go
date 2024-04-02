package app

import (
	"encoding/json"
	"fmt"
	"github.com/strangerstill/url-shorter/internal/models"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const baseURL = "https://localhost:8765"

func makeHandler() http.Handler {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	loggerItem := zap.Must(zap.NewDevelopment())
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(loggerItem)
	return MakeRouter(NewHandlers(*baseURL), loggerItem)
}

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

	h := makeHandler()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\nrun test: %v body: %v\n", test.name, test.originURL)
			originURL := strings.NewReader(test.originURL)
			request := httptest.NewRequest(http.MethodPost, "/", originURL)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)

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
				code: http.StatusTemporaryRedirect,
			},
		},
	}

	h := makeHandler()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			originURL := strings.NewReader(test.originURL)
			postRequest := httptest.NewRequest(http.MethodPost, "/", originURL)
			postW := httptest.NewRecorder()
			h.ServeHTTP(postW, postRequest)
			postRes := postW.Result()

			if postRes.StatusCode == http.StatusCreated {
				resShortURL, _ := io.ReadAll(postRes.Body)

				getRequest := httptest.NewRequest(http.MethodGet, string(resShortURL), nil)
				getW := httptest.NewRecorder()
				h.ServeHTTP(getW, getRequest)
				getRes := getW.Result()

				fmt.Printf("expected code: %d, status code %d\n", test.want.code, getRes.StatusCode)
				assert.Equal(t, test.want.code, getRes.StatusCode)
				if getRes.StatusCode == http.StatusTemporaryRedirect {
					fmt.Printf("want location = %s location %s\n", test.originURL, getRes.Header.Get("Location"))
					assert.Equal(t, getRes.Header.Get("Location"), test.originURL)
				}
			}
		})
	}
}

func TestSetShortURLJSON(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name      string
		payload   string
		originURL string
		want      want
	}{
		{
			name:      "valid originURL JSON",
			originURL: "https://pkg.go.dev/testing",
			want: want{
				code: http.StatusTemporaryRedirect,
			},
		},
	}
	h := makeHandler()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\nrun test: %v\n", test.name)
			originURL, err := json.Marshal(models.PayloadUrl{Url: test.originURL})
			if err != nil {
				log.Println(err)
				return
			}
			postRequest := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(string(originURL)))
			postW := httptest.NewRecorder()
			h.ServeHTTP(postW, postRequest)
			postRes := postW.Result()
			assert.Equal(t, postRes.StatusCode, http.StatusCreated)
			if postRes.StatusCode == http.StatusCreated {
				resShortURL, _ := io.ReadAll(postRes.Body)

				var data models.ResultUrl
				err := json.Unmarshal(resShortURL, &data)
				if err != nil {
					log.Println(err)
					return
				}
				getRequest := httptest.NewRequest(http.MethodGet, data.Result, nil)
				getW := httptest.NewRecorder()
				h.ServeHTTP(getW, getRequest)
				getRes := getW.Result()

				fmt.Printf("expected code: %d, status code %d\n", test.want.code, getRes.StatusCode)
				assert.Equal(t, test.want.code, getRes.StatusCode)
				if getRes.StatusCode == http.StatusTemporaryRedirect {
					fmt.Printf("want location = %s location %s\n", test.originURL, getRes.Header.Get("Location"))
					assert.Equal(t, getRes.Header.Get("Location"), test.originURL)
				}
			}
		})
	}
}
