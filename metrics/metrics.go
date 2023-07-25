package metrics

import (
	"net/http"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/metrics/prom"
)

// blocks
func Listen(c ctx.C, addr string) error {
	s := http.Server{
		Addr: addr,
		//Handler: promhttp.Handler(),
		Handler: prom.Handler(),
	}
	log.Infof(c, "promhttp listening to %q", s.Addr)
	err := s.ListenAndServe()
	switch err {
	case http.ErrServerClosed:
	default:
		log.Warnf(c, "monitoring: %v", err)
	}
	return err
}
