package loggger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type responseData struct {
	size   int
	status int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	var sugar = *logger.Sugar()
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					return
				}
			}()
			responseData := &responseData{
				size:   0,
				status: 0,
			}
			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}
			start := time.Now()
			next.ServeHTTP(&lw, req)
			duration := time.Since(start)
			sugar.Infoln(
				"uri", req.RequestURI,
				"method", req.Method,
				"status", responseData.status,
				"duration", duration,
				"size", responseData.size,
			)
		}

		return http.HandlerFunc(fn)
	}
}
