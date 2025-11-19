package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"gitlab.com/romalor/rika"

	"gitlab.com/romalor/htmx-contacts/app/routes/contacts"
	"gitlab.com/romalor/htmx-contacts/pkg/debug"
	"gitlab.com/romalor/htmx-contacts/pkg/stores/contacts"
	"gitlab.com/romalor/htmx-contacts/pkg/tpl"
)

func main() {
	// create the logger.
	log := slog.New(slog.Default().Handler())

	srv := rika.New()

	// serve static content.
	srv.Mux.FileServer("/static/*file", http.FS(os.DirFS("static")))

	// read config and register routes.
	app.Routes(srv.Mux, appConfig())

	// optional pprof handlers.
	go runServer(debug.Mux(), "9000")
	log.Info("debug mux started", "port", 9000)

	_ = srv.Start(":8080")	
}

func appConfig() app.Config {
	s, err := contacts.NewStore("contacts.json")
	handleErr(err)

	tpl, err := tpl.NewBundle(
		"templates/*.html",
		[]string{
			"templates/layouts/*.html",
			"templates/components/*.html",
		},
	)
	handleErr(err)

	return app.Config{
		TplBundle:tpl, 		
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
