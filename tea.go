package tea

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	logger = log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

type Tea struct {
	mux         *httptreemux.ContextMux
	middlewares alice.Chain
	Logger      *log.Logger
}

func New() *Tea {
	return &Tea{
		mux:         httptreemux.NewContextMux(),
		middlewares: alice.New(LoggingMiddleware),
		Logger:      logger,
	}
}

func (t *Tea) Get(path string, h http.HandlerFunc) {
	t.mux.GET(path, t.middlewares.ThenFunc(h).ServeHTTP)
}

func (t *Tea) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.mux.ServeHTTP(NewResponseWriter(w), r)
}

func (t *Tea) Start(addr string) error {
	t.Logger.Printf("🍵  is served on %s\n", addr)
	return http.ListenAndServe(addr, t)
}
