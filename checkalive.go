// checkalive allows you send & record alive signals cyclically
package checkalive

import "sync"

// Checker collects alive signals.
type Checker interface {
	// Ping sends a alive signal.
	Ping(topic string) error

	// Report gets a copy of current topics counter.
	Report() map[string]uint64

	// Reset all topics counter.
	Reset() error
}

type checker struct {
	lock          *sync.RWMutex
	topicsCounter map[string]uint64
}

func NewChecker() Checker {
	return &checker{
		lock:          &sync.RWMutex{},
		topicsCounter: make(map[string]uint64),
	}
}

func (c *checker) Ping(topic string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	count, present := c.topicsCounter[topic]
	if !present {
		count = 0
	}
	c.topicsCounter[topic] = count + 1

	return nil
}

func (c *checker) Report() map[string]uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.topicsCounter
}

func (c *checker) Reset() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.topicsCounter = make(map[string]uint64)

	return nil
}
