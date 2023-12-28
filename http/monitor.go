package http

import (
	"net/http"
	nprof "net/http/pprof"
	"runtime/pprof"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/utils/prom"
)

func (this *Server) SetupPrometheus(c ctx.C, path string) {
	if path == "" {
		path = "/metrics"
	}
	this.mux.Handle(path, prom.Handler())
}

// register /pprof/ and /goroutines under the given prefix
func (this *Server) SetupMonitoring(c ctx.C, prefix string) {
	this.mux.HandleFunc(prefix+"/pprof/", nprof.Index)
	this.mux.HandleFunc(prefix+"/pprof/cmdline", nprof.Cmdline)
	this.mux.HandleFunc(prefix+"/pprof/profile", nprof.Profile)
	this.mux.HandleFunc(prefix+"/pprof/symbol", nprof.Symbol)
	this.mux.HandleFunc(prefix+"/pprof/trace", nprof.Trace)
	this.mux.HandleFunc(prefix+"/goroutines", func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		p := pprof.Lookup("goroutine")
		if p == nil {
			w.WriteHeader(204)
			return
		}
		err := p.WriteTo(w, 1)
		if err != nil {
			log.Errorf(c, "can't write profile: %v", err)
		}
	})
}
