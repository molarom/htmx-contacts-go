package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/htmx-contacts/app"
	"gitlab.com/romalor/htmx-contacts/debug"
	"gitlab.com/romalor/htmx-contacts/stores/contacts"
	"gitlab.com/romalor/htmx-contacts/tpl"
)

func main() {
	go RunDebugServer()
	RunRoxiServer()
}

func appConfig() app.Config {
	s, err := contacts.NewStore("contacts.json")
	handleErr(err)

	return app.Config{
		TplBundle: tpl.NewBundle("base", "templates/layouts/*.html", "templates/*.html"),
		Store:     s,
	}
}

func RunRoxiServer() {
	mux := roxi.NewWithDefaults(
		roxi.WithLogger(slog.New(slog.Default().Handler()).Info),
		roxi.WithOptionsHandler(roxi.DefaultCORS),
	)

	mux.FileServer("/static/*file", http.FS(os.DirFS("static")))

	app.Routes(mux, appConfig())
	mux.PrintTree()

	runServer(mux, "8080")
}

func RunDebugServer() {
	runServer(debug.Mux(), "9000")
}

func runServer(h http.Handler, port string) {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h,
	}

	if err := srv.ListenAndServe(); err != nil {
		handleErr(err)
	}
}

func handleErr(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
