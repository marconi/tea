package tea

import "net/http"

// strip down version of https://godoc.org/github.com/codegangsta/negroni#ResponseWriter
// at the moment it only supports recording of status
type ResponseWriter interface {
	http.ResponseWriter
	Status() int
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{ResponseWriter: rw}
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}
