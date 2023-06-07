package shutdown

import (
	"errors"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

type shutter struct {
	// channel used to broadcast a shutdown
	ch   chan struct{}
	ch5  chan struct{}
	once sync.Once
	// active services holding the shutdown to complete
	wg sync.WaitGroup
}

var shutdowner = newShutter()

func newShutter() *shutter {
	return &shutter{
		ch:  make(chan struct{}),
		ch5: make(chan struct{}),
	}
}

// start a global shutdown unless already started
func (this *shutter) begin() {
	this.once.Do(func() {
		close(this.ch)
		go func() {
			time.Sleep(5 * time.Second)
			close(this.ch5)
		}()
	})
}

// return a closed channel when the shutdown has started
func (this *shutter) started() <-chan struct{} {
	return this.ch
}

// return a closed channel 5 seconds after the shutdown has started
func (this *shutter) started5Sec() <-chan struct{} {
	return this.ch5
}

// returns a channel that will close when the shutdown has completed
func (this *shutter) done() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		this.wg.Wait()
		close(ch)
	}()
	return ch
}

// prevent the shutdown to complete until released
func (this *shutter) hold() ReleaseFn {
	this.wg.Add(1)
	return this.wg.Done
}

// prevents the shutdown to complete until released, and also wait for the shutdown to start
func (this *shutter) holdAndWait() ReleaseFn {
	this.wg.Add(1)
	<-this.ch
	return this.wg.Done
}

// setup signals and wait for the shutdown to complete
func (this *shutter) waitForSignal(c ctx.C, cf ctx.CancelFunc) {
	sigs := make(chan os.Signal, 3)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// TODO(oha): should we keep this here? seems a bit like unexpected magic to me
	sigusr := make(chan os.Signal, 1)

	for i := 0; ; {
		log.Infof(c, "waiting for signal...")
		select {
		case <-this.done(): // done, we can quit
			log.Warnf(c, "shutdown complete: %v", c.Err())
			return
		case <-sigusr:
			log.Warnf(c, "SIGURS1 -- dump threads")
			_ = pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		case sig := <-sigs:
			switch i {
			case 0:
				log.Warnf(c, "got SIG %q: start a graceful shutdown...", sig.String())
				this.begin()
				i++
			case 1:
				log.Warnf(c, "got SIG %q: canceling root context...", sig.String())
				cf(errors.New("SIGINT"))
				i++
			default:
				log.Warnf(c, "got SIG %q: os.Exit(-1)...", sig.String())
				os.Exit(-1)
			}
		}
	}
}
