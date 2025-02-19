package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/roxy/app"
	"gitlab.com/romalor/roxy/stdlib/tpl"
)

func main() {
	RunRoxiServer()
}

func appConfig() app.Config {
	db, err := os.ReadFile("contacts.json")
	handleErr(err)

	var contacts app.Contacts
	handleErr(json.Unmarshal(db, &contacts))

	return app.Config{
		TplBundle:    tpl.NewBundle("base", "templates/layouts/*.html", "templates/*.html"),
		ContactStore: contacts,
	}
}

func RunRoxiServer() {
	mux := roxi.New(
		roxi.WithLogger(slog.New(slog.Default().Handler()).Info),
		roxi.WithOptionsHandler(roxi.HandlerFunc(roxi.DefaultCORS)),
	)

	// mux.Handler("GET", "/static", http.StripPrefix("/static", http.FileServerFS(os.DirFS("static"))))

	app.Routes(mux, appConfig())
	mux.PrintRoutes()

	runServer(mux)
}

func RunStdServer() {
	mux := http.NewServeMux()

	// app.StdRoutes(mux, appConfig())

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})

	runServer(mux)
}

func RunHTTPRouter() {
	mux := httprouter.New()
	mux.RedirectTrailingSlash = false
	mux.RedirectFixedPath = false

	mux.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(204)
	})

	runServer(mux)
}

func runServer(h http.Handler) {
	srv := &http.Server{
		Addr:    ":8080",
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
