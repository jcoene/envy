package envy

import (
	"fmt"
	"path"
	"time"

	"github.com/nu7hatch/gouuid"
)

// A handler is a method that consumes a lock. It is responsible
// for watching for loss of the lock and releasing it when done.
type Handler func(*Lock)

type Lock struct {
	Lost chan int // Handler should read to indicate the lock was lost

	name     string
	key      string
	ttl      uint64
	id       string
	released chan int
	c        *Client
}

func (c *Client) NewLock(name string, ttl int) *Lock {
	id, err := uuid.NewV4()
	if err != nil {
		panic("unable to obtain uuid")
	}

	return &Lock{
		Lost: make(chan int),

		name:     name,
		key:      path.Join(c.ns, "locks", name),
		ttl:      uint64(ttl),
		id:       id.String(),
		released: make(chan int),
		c:        c,
	}
}

// Obtain a lock and run the given method
func (l *Lock) With(fn Handler) {
	l.wait()
	go l.keep()
	fn(l)
}

func (l *Lock) Release() {
	l.released <- 1
	l.c.c.CompareAndDelete(l.key, l.id, 0)
}

func (l *Lock) try() bool {
	_, err := l.c.c.Create(l.key, l.id, l.ttl)
	return err == nil
}

func (l *Lock) wait() {
	t := time.NewTicker(time.Duration(1) * time.Second)

	for !l.try() {
		<-t.C
	}
}

func (l *Lock) keep() {
	t := time.NewTicker(time.Duration(l.ttl/2) * time.Second)
	for {
		select {
		case <-l.released:
			return
		case <-t.C:
			if _, err := l.renew(); err != nil {
				fmt.Printf("unable to keep lock: %s\n", err)
				l.Lost <- 1
				return
			}
		}
	}
}

func (l *Lock) renew() (bool, error) {
	var err error

	for r := 0; r < 3; r++ {
		if _, err = l.c.c.CompareAndSwap(l.key, l.id, l.ttl, l.id, 0); err == nil {
			return true, nil
		}
		time.Sleep(1 * time.Second)
	}

	return false, err
}
