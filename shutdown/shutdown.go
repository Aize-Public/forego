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

// Used to signal an operation can't be executed because the system is shutting down
var Err = errors.New("shutdown")

// channel used to broadcast a shutdown
var ch = make(chan struct{})
var ch5 = make(chan struct{})

var once sync.Once

// active services holding the shutdown to complete
var wg sync.WaitGroup

// start a global shutdown unless already started
func Begin() {
	once.Do(func() {
		close(ch)
		go func() {
			time.Sleep(5 * time.Second)
			close(ch5)
		}()
	})
}

// return a closed channel when the shutdown has started
func Started() <-chan struct{} {
	return ch
}

// return a closed channel 5 seconds after the shutdown has started
func Started5Sec() <-chan struct{} {
	return ch5
}

// returns a channel that will close when the shutdown has completed
func Done() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

// prevent the shutdown to complete until released
func Hold() ReleaseFn {
	wg.Add(1)
	return wg.Done
}

// prevents the shutdown to complete until released, and also wait for the shutdown to start
func HoldAndWait() ReleaseFn {
	wg.Add(1)
	<-ch
	return wg.Done
}

type ReleaseFn func()

func (fn ReleaseFn) Release() { fn() }

// setup signals and wait for the shutdown to complete
func WaitForSignal(c ctx.C, cf ctx.CancelFunc) {
	sigs := make(chan os.Signal, 3)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// TODO(oha): should we keep this here? seems a bit like unexpected magic to me
	sigusr := make(chan os.Signal, 1)

	for i := 0; ; {
		log.Infof(c, "waiting for signal...")
		select {
		case <-Done(): // done, we can quit
			log.Warnf(c, "shutdown complete: %v", c.Err())
			return
		case <-sigusr:
			log.Warnf(c, "SIGURS1 -- dump threads")
			_ = pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		case sig := <-sigs:
			switch i {
			case 0:
				log.Warnf(c, "got SIG %q: start a graceful shutdown...", sig.String())
				Begin()
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
