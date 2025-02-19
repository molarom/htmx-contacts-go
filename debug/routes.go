package debug

import (
	"net/http"
	"net/http/pprof"

	"gitlab.com/romalor/roxi"
)

func Mux() *roxi.Mux {
	mux := roxi.NewWithDefaults()
	mux.Handler("GET", "/debug/pprof/*idx", http.HandlerFunc(pprof.Index))
	mux.Handler("GET", "/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handler("GET", "/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handler("GET", "/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handler("GET", "/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	return mux
}
