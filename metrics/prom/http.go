package prom

import (
	"log"
	"net/http"
)

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range Registry {
			err := m.Print(w)
			if err != nil {
				w.WriteHeader(500)
				log.Printf("can't print: %v", err)
				return
			}
		}
	}
}
