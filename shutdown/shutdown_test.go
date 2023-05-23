package shutdown_test

import (
	"testing"
	"time"

	"github.com/Aize-Public/forego/shutdown"
	"github.com/Aize-Public/forego/test"
)

func TestShutdown(t *testing.T) {
	seq := make(chan int, 10)
	go func() {
		t.Logf("waiting for shutdown...")
		defer shutdown.HoldAndWait().Release()
		seq <- 1
		t.Logf("releasing hold")
	}()
	go func() {
		time.Sleep(time.Millisecond)
		seq <- 0
		t.Logf("starting shutdown...")
		shutdown.Begin()
		<-shutdown.Done()
		t.Logf("shutdown finished")
		seq <- 2
	}()
	test.EqualsGo(t, 0, <-seq)
	test.EqualsGo(t, 1, <-seq)
	test.EqualsGo(t, 2, <-seq)
}
