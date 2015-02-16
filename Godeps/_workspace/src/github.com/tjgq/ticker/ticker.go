// Package ticker implements a ticker that can be stopped and re-started.
package ticker

import (
	"sync"
	"time"
)

// A Ticker holds a channel that delivers ticks at intervals.
type Ticker struct {
	C       chan time.Time // The channel on which ticks are delivered.
	d       time.Duration
	mu      sync.Mutex
	running bool
	stop    chan bool
}

// New returns a new ticker that ticks every d seconds. It adjusts the
// intervals or drops ticks to make up for slow receivers. The ticker
// is initially in the stopped state.
func New(d time.Duration) *Ticker {
	if d <= 0 {
		panic("ticker: non-positive duration")
	}
	return &Ticker{
		C:       make(chan time.Time),
		d:       d,
		running: false,
		stop:    make(chan bool),
	}
}

// Start (re-)starts the ticker. Ticks will be delivered on the ticker's
// channel until Stop is called.
func (t *Ticker) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.running {
		go t.loop()
		t.running = true
	}
}

// Stop stops the ticker. No ticks will be delivered on the ticker's channel
// after Stop returns and before Start is called again.
func (t *Ticker) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.running {
		t.stop <- true
		t.running = false
	}
}

// Stopped returns whether the ticker is stopped.
func (t *Ticker) Stopped() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return !t.running
}

func (t *Ticker) loop() {
	tk := time.NewTicker(t.d)
	for {
		select {
		case tm := <-tk.C:
			select {
			case t.C <- tm:
			default:
			}
		case <-t.stop:
			tk.Stop()
			return
		}
	}
}
