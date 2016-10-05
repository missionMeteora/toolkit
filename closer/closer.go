package closer

import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// New returns a new instance of Closer
func New() (c *Closer) {
	c = &Closer{
		// Create buffered chan with a lenght of one for our message channel
		mc: make(chan error, 1),
	}

	go c.listen()

	// Give goroutines 3 milliseconds to spin up
	time.Sleep(time.Millisecond * 3)
	return
}

// Closer manages the closing of services
type Closer struct {
	// Message channel
	mc chan error
	// Closed state
	cs int32
}

// listen will listen for closing signals (interrupt, terminate, abort, quit) and call close
func (c *Closer) listen() {
	// sc represents the signal channel
	sc := make(chan os.Signal, 1)
	// Listen for signal notifications
	// Discussion topic: Should we include SIGQUIT? If we catch the signal, we won't get to see the unwind
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	// Signal received
	<-sc
	c.Close(nil)
}

// Wait will wait until it receives a close notification
// If an error is the reason for closure - said error will be returned
func (c *Closer) Wait() (err error) {
	err = <-c.mc
	return
}

// Close will close selected instance of Closer
func (c *Closer) Close(err error) (ok bool) {
	if ok = atomic.CompareAndSwapInt32(&c.cs, 0, 1); !ok {
		return
	}

	c.mc <- err
	return
}
