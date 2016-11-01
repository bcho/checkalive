// checkalive allows you send & record alive signals cyclically
package checkalive

import "sync"

// Checker collects alive signals.
type Checker interface {
	// Ping sends a alive signal.
	Ping(topic string) error

	// Report gets a copy of current topics counter.
	Report() map[string]uint64

	// Reset all topics counter. Returns current topics counter copy.
	Reset() (map[string]uint64, error)
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

func (c *checker) copyCounter() map[string]uint64 {
	r := make(map[string]uint64)

	for k, v := range c.topicsCounter {
		r[k] = v
	}

	return r
}

func (c *checker) Report() map[string]uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.copyCounter()
}

func (c *checker) Reset() (map[string]uint64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	rv := c.copyCounter()
	c.topicsCounter = make(map[string]uint64)

	return rv, nil
}
