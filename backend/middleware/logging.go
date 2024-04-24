package middleware

import (
	"log"
	"net/http"
	"time"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *wrappedResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &wrappedResponseWriter{w, http.StatusOK}
		next.ServeHTTP(ww, r)

		log.Println(ww.status, r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	}))
}
