package tea

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

type Middleware func(http.Handler) http.Handler

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
		mux: httptreemux.NewContextMux(),
		middlewares: alice.New(
			alice.Constructor(NewLoggingMiddleware(logger)),
		),
		Logger: logger,
	}
}

func (t *Tea) Get(path string, h http.HandlerFunc) {
	t.mux.GET(path, t.applyMiddlewares(h))
}

func (t *Tea) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.mux.ServeHTTP(NewResponseWriter(w), r)
}

func (t *Tea) Use(h Middleware) {
	// push new middleware at the top
	t.middlewares = alice.New(alice.Constructor(h)).Extend(t.middlewares)
}

func (t *Tea) Apply(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	var wh http.Handler = http.HandlerFunc(h)
	for _, m := range middlewares {
		wh = m(wh)
	}
	return wh.ServeHTTP
}

func (t *Tea) applyMiddlewares(h http.HandlerFunc) http.HandlerFunc {
	return t.middlewares.ThenFunc(h).ServeHTTP
}

func (t *Tea) Start(addr string) error {
	t.Logger.Printf("🍵  is served on %s\n", addr)
	return http.ListenAndServe(addr, t)
}
