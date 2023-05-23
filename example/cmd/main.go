package main

import (
	"os"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/example"
	"github.com/Aize-Public/forego/http"
	"github.com/Aize-Public/forego/shutdown"
)

func main() {
	c, cf := ctx.Background()
	log.Warnf(c, "init")
	defer log.Warnf(c, "exit")

	s := http.NewServer(c)
	store := example.Store{
		Data: map[string]any{
			"true": true,
		},
	}
	http.RegisterAPI(c, s, &example.Get{Store: &store})
	http.RegisterAPI(c, s, &example.Set{Store: &store})

	addr, err := s.Listen(c, "127.0.0.1:0")
	if err != nil {
		log.Errorf(c, "err")
		os.Exit(-1)
	}

	log.Infof(c, "listening %v", addr.String())

	shutdown.WaitForSignal(c, cf)
}
