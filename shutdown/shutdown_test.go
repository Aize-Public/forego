package shutdown

import (
	"testing"

	"github.com/Aize-Public/forego/test"
)

func TestShutdown(t *testing.T) {
	seq := make(chan int, 10)
	shutdowner := newShutter()
	release := shutdowner.hold()
	go func() {
		t.Logf("waiting for shutdown...")
		<-shutdowner.started()
		seq <- 1
		t.Logf("releasing hold")
		release()
	}()
	go func() {
		seq <- 0
		t.Logf("starting shutdown...")
		shutdowner.begin()
		<-shutdowner.done()
		t.Logf("shutdown finished")
		seq <- 2
	}()
	test.EqualsGo(t, 0, <-seq)
	test.EqualsGo(t, 1, <-seq)
	test.EqualsGo(t, 2, <-seq)
}
