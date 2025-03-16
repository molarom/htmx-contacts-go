package htmx

import (
	"context"
	"net/http"
	"strconv"

	"gitlab.com/romalor/roxi"
)

type HX struct {
	Boosted               bool
	CurrentURL            string
	HistoryRestoreRequest bool
	Prompt                string
	Request               bool
	TriggerName           string
	Trigger               string
}

func HTMX(next roxi.HandlerFunc) roxi.HandlerFunc {
	return func(ctx context.Context, r *http.Request) error {
		boosted, _ := strconv.ParseBool(r.Header.Get("HX-Boosted"))
		hist, _ := strconv.ParseBool(r.Header.Get("HX-History-Restore-Request"))
		req, _ := strconv.ParseBool(r.Header.Get("HX-Request"))
		ctx = set(ctx, HX{
			Boosted:               boosted,
			HistoryRestoreRequest: hist,
			Request:               req,
			CurrentURL:            r.Header.Get("HX-Current-URL"),
			Prompt:                r.Header.Get("HX-Prompt"),
			Trigger:               r.Header.Get("HX-Trigger"),
			TriggerName:           r.Header.Get("HX-Trigger-Name"),
		})
		return next(ctx, r)
	}
}
