package tea

import (
	"net/http"
	"time"
)

type Logger interface {
	Infof(string, ...interface{})
}

// An un-customizable logging middleware
func NewLoggingMiddleware(loggers ...Logger) Middleware {
	var logger Logger
	if len(loggers) > 0 {
		logger = loggers[0]
	}

	return func(h http.Handler) http.Handler {
		hf := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			duration := time.Since(start)

			rw := w.(ResponseWriter)
			logger.Infof("%s %s %s %d %s", r.Method, r.URL.Path, r.Proto, rw.Status(), duration)
		}
		return http.HandlerFunc(hf)
	}
}
