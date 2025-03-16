package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/htmx-contacts/app/routes/contacts"
	"gitlab.com/romalor/htmx-contacts/pkg/debug"
	"gitlab.com/romalor/htmx-contacts/pkg/middleware/errs"
	"gitlab.com/romalor/htmx-contacts/pkg/middleware/htmx"
	"gitlab.com/romalor/htmx-contacts/pkg/middleware/logging"
	"gitlab.com/romalor/htmx-contacts/pkg/stores/contacts"
	"gitlab.com/romalor/htmx-contacts/pkg/tpl"
)

func main() {
	// create the logger.
	log := slog.New(slog.Default().Handler())

	// add logging to 404 responses.
	notfound := func(w http.ResponseWriter, r *http.Request) {
		log.Warn("no route registered", "method", r.Method, "path", r.URL.Path)
		roxi.HandlerFunc(roxi.NotFound).ServeHTTP(w, r)
	}

	// create the mux.
	mux := roxi.NewWithDefaults(
		roxi.WithOptionsHandler(roxi.DefaultCORS),
		roxi.WithNotFoundHandler(http.HandlerFunc(notfound)),
		roxi.WithMiddleware(
			logging.Logging(log.Info),
			htmx.HTMX,
			errs.Errors(log.Error),
		),
	)

	// serve static content.
	mux.FileServer("/static/*file", http.FS(os.DirFS("static")))

	// read config and register routes.
	app.Routes(mux, appConfig())
	mux.PrintRoutes()

	// optional pprof handlers.
	go runServer(debug.Mux(), "9000")

	// run the app.
	runServer(mux, "8080")
}

func appConfig() app.Config {
	s, err := contacts.NewStore("contacts.json")
	handleErr(err)

	return app.Config{
		TplBundle: tpl.NewBundle("base",
			"templates/*.html",
			"templates/layouts/*.html",
			"templates/components/*.html",
		),
		Store: s,
	}
}

func runServer(h http.Handler, port string) {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h,
	}

	handleErr(srv.ListenAndServe())
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
