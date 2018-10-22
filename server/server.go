package server

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

var handled int64
var nbr = stats.Int64("http/requests", "requests count", "")

//
// Telemetry and tracing code
//

func prepareTracing(addressjaeger string) {
	if addressjaeger == "" {
		log.Fatal().Msg("no endpoint for jaeger defined")
	}
	exporter, err := jaeger.NewExporter(jaeger.Options{
		Endpoint:    addressjaeger,
		ServiceName: "demo"},
	)
	if err != nil {
		log.Fatal().Msg("can't create jaeger exporter")
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}

func prepareTelemetry(pe *prometheus.Exporter) {
	handled = 0

	viewCount := &view.View{
		Name:        "http_count",
		Description: "number of http requests made",
		TagKeys:     nil,
		Measure:     nbr,
		Aggregation: view.LastValue(),
	}

	view.RegisterExporter(pe)
	view.Register(viewCount)
	view.SetReportingPeriod(10 * time.Second)
}

//
// HTTP Server code
//

func handler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Request handled")
	ctx, span := trace.StartSpan(context.Background(), "demo.server.handler")
	span.Annotate(nil, "user")
	defer span.End()

	handled += 1
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("I've been asked to serve SO MANY TIMES ... at least %d\n", int(handled))))

	stats.Record(ctx, nbr.M(handled))
}

func Serve(port int, jaegerurl string) {
	handled = 0
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "demo",
	})
	if err != nil {
		log.Error().Msg("fail to create exporter")
	}

	prepareTracing(jaegerurl)
	prepareTelemetry(pe)

	mux := http.NewServeMux()
	mux.Handle("/metrics", pe)
	mux.HandleFunc("/", handler)
	h := &ochttp.Handler{Handler: mux}
	log.Info().Msg("starting http handler")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), h); err != nil {
		log.Error().Msg("fail to start http server")
		runtime.Goexit()
	}
}
