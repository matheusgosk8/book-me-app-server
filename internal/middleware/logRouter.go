package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func LogRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		elapsed := float64(time.Since(start).Microseconds()) / 1000.0
		log.Printf("[%s] %s - finished in %.3f ms", r.Method, r.URL.Path, elapsed)
	})
}
