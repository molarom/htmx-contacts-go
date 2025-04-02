package debug

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/arl/statsviz"
	"gitlab.com/romalor/roxi"
)

func Mux() *roxi.Mux {
	mux := roxi.New()
	mux.HandlerFunc("GET", "/debug/pprof/*idx", pprof.Index)
	mux.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
	mux.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
	mux.HandlerFunc("GET", "/debug/pprof/trace", pprof.Trace)

	viz, _ := statsviz.NewServer()
	mux.GET("/debug/statsviz/*filepath", func(ctx context.Context, r *http.Request) error {
		if r.PathValue("filepath") == "/ws" {
			viz.Ws().ServeHTTP(roxi.GetWriter(ctx), r)
		} else {
			viz.Index().ServeHTTP(roxi.GetWriter(ctx), r)
		}
		return nil
	})

	return mux
}
