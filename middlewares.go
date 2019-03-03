package tea

import (
	"net/http"
	"time"
)

// An un-customizable logging middleware
func LoggingMiddleware(h http.Handler) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(start)

		rw := w.(ResponseWriter)
		logger.Infof("%s %s %s %d %s", r.Method, r.URL.Path, r.Proto, rw.Status(), duration)
	}
	return http.HandlerFunc(hf)
}
