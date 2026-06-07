package server

import (
	"time"
)

// Mock consumer structure representing the fix in consumer.go
type consumer struct {
	mu       sync.Mutex
	isLeader bool
	active   bool
	stream   *stream
	// other fields...
}

func (c *consumer) setLeader(isLeader bool) {
	c.mu.Lock()
	wasLeader := c.isLeader
	c.isLeader = isLeader
	active := c.active
	c.mu.Unlock()

	if isLeader && !wasLeader {
		go c.enableLoop(active)
	} else if !isLeader && wasLeader {
		c.disableLoop()
	}
}

func (c *consumer) enableLoop(active bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isLeader {
		return
	}
	// Wait for stream to be ready and subscription interest to propagate
	go func() {
		for i := 0; i < 10; i++ {
			c.mu.Lock()
			isLeader := c.isLeader
			streamReady := c.stream != nil && c.stream.isReady()
			c.mu.Unlock()

			if !isLeader {
				return
			}
			if streamReady && c.hasSubInterest() {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		c.mu.Lock()
		if c.isLeader {
			c.playLoop()
		}
		c.mu.Unlock()
	}()
}

func (c *consumer) hasSubInterest() bool {
	// Check if the delivery subject has active subscription interest in the cluster
	return true
}

func (c *consumer) playLoop() {
	// Delivery loop logic...
}

func (c *consumer) disableLoop() {
	// Stop delivery loop...
}
