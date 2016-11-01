package checkalive

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCheckerSync(t *testing.T) {
	c := NewChecker()

	topic1 := "foo"
	topic2 := "bar"

	if err := c.Ping(topic1); err != nil {
		t.Errorf("unexpected err: %+v", err)
	}
	if c := c.Report()[topic1]; c != 1 {
		t.Errorf("expect 1, got: %d", c)
	}

	if err := c.Ping(topic1); err != nil {
		t.Errorf("unexpected err: %+v", err)
	}
	if c := c.Report()[topic1]; c != 2 {
		t.Errorf("expect 2, got: %d", c)
	}

	if err := c.Ping(topic2); err != nil {
		t.Errorf("unexpected err: %+v", err)
	}
	if c := c.Report()[topic2]; c != 1 {
		t.Errorf("expect 1, got: %d", c)
	}

	counter, err := c.Reset()
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}
	if l := len(counter); l != 2 {
		t.Errorf("map size should be 2, got: %d", l)
	}

	counter, err = c.Reset()
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}
	if l := len(counter); l != 0 {
		t.Errorf("map size should be 0, got: %d", l)
	}
}

func TestCheckerAsync(t *testing.T) {
	c := NewChecker()

	topic1 := "foo"
	topic2 := "bar"
	topic3 := "baz"
	var done sync.WaitGroup

	pingTopic := func(topic string) {
		done.Add(1)
		go func() {
			defer done.Done()

			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			if err := c.Ping(topic); err != nil {
				t.Errorf("unexpected err: %+v", err)
			}
		}()
	}

	for i := 0; i < 10; i++ {
		pingTopic(topic1)
		pingTopic(topic2)
		pingTopic(topic3)
	}

	done.Wait()

	r := c.Report()

	if l := len(r); l != 3 {
		t.Errorf("map size should be 3, got: %d", l)
	}

	if c := r[topic1]; c != 10 {
		t.Errorf("should ping 10 times, got: %d", c)
	}
	if c := r[topic2]; c != 10 {
		t.Errorf("should ping 10 times, got: %d", c)
	}
	if c := r[topic3]; c != 10 {
		t.Errorf("should ping 10 times, got: %d", c)
	}

	counter, err := c.Reset()
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}
	if l := len(counter); l != 3 {
		t.Errorf("map size should be 3, got: %d", l)
	}

	counter, err = c.Reset()
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}
	if l := len(counter); l != 0 {
		t.Errorf("map size should be 0, got: %d", l)
	}
}
