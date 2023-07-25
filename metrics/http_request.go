package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Aize-Public/forego/metrics/prom"
)

/*
	Usage

	func x(path) {
		defer metrics.HttpRequest{
			Path: path,
			Code: 200,
		}.Observe(time.Now())
		// ...
	}

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		m := metrics.HttpRequest{Path: path} // create metrics for the path
		defer m.ObserveSince(time.Now()) // defer sample the elapsed time
		h(m.AsResponseWriter(w), r) // use a custom response writer to track the status code
	})
*/

type HttpRequest struct {
	Path string
	Code int
}

var httpRequest = prom.Register(&prom.Histogram{
	Name:    "http_request",
	Buckets: []float64{.001, .0025, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30, 60, 60 * 5}, // 1ms .. 5'
	Labels:  []string{"path", "code"},
})

func (m *HttpRequest) ObserveSince(start time.Time) {
	// note: we use a pointer receiver to allow AsResponseWriter to modify HttpRequest
	httpRequest.Observe(time.Since(start).Seconds(), m.Path, fmt.Sprint(m.Code))
}

// wraps the given response write with one that updates the metric
// this MUST be used, or you will miss the either explicit or implicit WriteHeader()
func (m *HttpRequest) AsResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, m}
}

type responseWriter struct {
	http.ResponseWriter
	m *HttpRequest
}

var _ http.ResponseWriter = &responseWriter{}

func (w *responseWriter) WriteHeader(code int) {
	w.m.Code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.m.Code == 0 {
		w.m.Code = 200
	}
	return w.ResponseWriter.Write(b)
}
