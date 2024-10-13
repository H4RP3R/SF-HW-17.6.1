package main

import (
	"fmt"
	"sync"
)

type counter struct {
	mutex sync.RWMutex
	data  map[string]int
}

func (c *counter) Read(key string) int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.data[key]
}

func (c *counter) Increment(key string) {
	c.mutex.Lock()
	c.data[key] += 1
	c.mutex.Unlock()
}

func (c *counter) PrintStats() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var total int
	for _, value := range c.data {
		total += value
	}

	fmt.Println()
	for key, value := range c.data {
		percent := float64(value) / float64(total) * 100
		fmt.Printf("%s: %d (%.2f%%)\n", key, value, percent)
	}
}

func NewCounter(keys ...string) *counter {
	c := counter{}
	c.data = make(map[string]int, len(keys))
	for _, k := range keys {
		c.data[k] = 0
	}

	return &c
}
