package ticker

import (
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	const delta = 100 * time.Millisecond
	const count = 5

	ticker := New(delta)

	time.Sleep(2 * delta)
	select {
	case <-ticker.C:
		t.Fatal("ticker created in started state")
	default:
	}

	ticker.Start()
	for i := 0; i < count; i++ {
		<-ticker.C
	}
	ticker.Stop()

	time.Sleep(2 * delta)
	select {
	case <-ticker.C:
		t.Fatal("ticker did not stop")
	default:
	}

	ticker.Start()
	for i := 0; i < count; i++ {
		<-ticker.C
	}
	ticker.Stop()
}

func TestDuration(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("New should have panicked")
		}
	}()
	New(-1)
}
