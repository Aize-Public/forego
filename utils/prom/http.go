package prom

import (
	"log"
	"net/http"
)

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := Range(func(name string, m Metric) error {
			err := m.Print(name, w)
			return err
		})
		if err != nil {
			w.WriteHeader(500)
			log.Printf("can't print: %v", err)
		}
	}
}
