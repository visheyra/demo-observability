package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
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
var yes int64
var no int64
var nbr = stats.Int64("http/requests", "requests count", "")
var nbrYes = stats.Int64("demo_observability/yes", "yes count", "")
var nbrNo = stats.Int64("demo_observability/no", "no count", "")

//
// Telemetry and tracing code
//

func prepareTracing(addressjaeger string) {
	if addressjaeger == "" {
		log.Fatal().Msg("no endpoint for jaeger defined")
	}
	exporter, err := jaeger.NewExporter(jaeger.Options{
		Endpoint:    addressjaeger,
		ServiceName: "demo-observability"},
	)
	if err != nil {
		log.Fatal().Msg("can't create jaeger exporter")
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}

func randomMessage(ctx context.Context) []byte {
	_, span := trace.StartSpan(ctx, "demo.server.handler")
	fmt.Println(ctx)
	defer span.End()

	log.Debug().Msg("Random Message Gen Reached")
	x := rand.Int()

	if x%2 == 0 {

		span.Annotate(
			[]trace.Attribute{
				trace.Int64Attribute("value", int64(x)),
				trace.StringAttribute("response", "no"),
			},
			"process no",
		)

		no += 1
		stats.Record(ctx, nbrNo.M(no))
		return []byte("no")
	}
	yes += 1
	stats.Record(ctx, nbrYes.M(yes))

	span.Annotate(
		[]trace.Attribute{
			trace.Int64Attribute("value", int64(x)),
			trace.StringAttribute("response", "yes"),
		},
		"process yes",
	)
	return []byte("yes")
}

func prepareTelemetry(pe *prometheus.Exporter) {
	handled = 0

	viewCount := &view.View{
		Name:        "demo_observ_http_count",
		Description: "number of http requests made",
		TagKeys:     nil,
		Measure:     nbr,
		Aggregation: view.LastValue(),
	}

	yesCount := &view.View{
		Name:        "demo_observ_yes_count",
		Description: "number of yes response made",
		TagKeys:     nil,
		Measure:     nbrYes,
		Aggregation: view.LastValue(),
	}

	noCount := &view.View{
		Name:        "demo_observ_no_count",
		Description: "number of no response made",
		TagKeys:     nil,
		Measure:     nbrNo,
		Aggregation: view.LastValue(),
	}

	view.RegisterExporter(pe)
	view.Register(viewCount, yesCount, noCount)
	view.SetReportingPeriod(10 * time.Second)
}

//
// HTTP Server code
//

func handler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Request handled")
	ctx, span := trace.StartSpan(r.Context(), "demo.server.handler")
	defer span.End()

	handled += 1
	w.WriteHeader(http.StatusOK)
	w.Write(randomMessage(ctx))

	stats.Record(ctx, nbr.M(handled))
}

func Serve(port int, jaegerurl string) {
	yes = 0
	no = 0
	handled = 0

	rand.Seed(time.Now().UTC().UnixNano())

	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "demo",
	})
	if err != nil {
		log.Fatal().Msg("fail to create exporter")
	}

	prepareTracing(jaegerurl)
	prepareTelemetry(pe)

	mux := http.NewServeMux()
	mux.Handle("/metrics", pe)
	mux.HandleFunc("/", handler)
	h := &ochttp.Handler{Handler: mux}
	log.Info().Msg("starting http handler")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), h); err != nil {
		log.Fatal().Msg("fail to start http server")
	}
}
