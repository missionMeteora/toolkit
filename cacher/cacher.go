package cacher

import (
	"sync"
	"time"
)

type (
	CacheFiller func() (data interface{}, err error)

	Cacher struct {
		data map[string]*cacheItem
		mux  sync.Mutex
	}
)

// New returns a new cacher
// purgeEvery sets the timer for purging old cache, set to 0 to never purge.
func New(purgeEvery time.Duration) *Cacher {
	c := &Cacher{
		data: map[string]*cacheItem{},
	}
	if purgeEvery > 0 {
		go c.purge(purgeEvery)
	}
	return c
}

func (c *Cacher) Get(key string, fn CacheFiller, ttl time.Duration) (data interface{}, err error) {
	var ci *cacheItem
	c.mux.Lock()
	if ci = c.data[key]; ci == nil {
		ci = &cacheItem{
			fn:  fn,
			ttl: int64(ttl / time.Second),
		}
		c.data[key] = ci
	}
	c.mux.Unlock()

	return ci.Value()
}

func (c *Cacher) Delete(key string) {
	c.mux.Lock()
	delete(c.data, key)
	c.mux.Unlock()
}

func (c *Cacher) purge(ttl time.Duration) {
	for {
		time.Sleep(ttl)
		ts := time.Now().Unix()
		c.mux.Lock()
		for key, ci := range c.data {
			if ci.expiresAt < ts {
				delete(c.data, key)
			}
		}
		c.mux.Unlock()
	}
}

type cacheItem struct {
	fn        CacheFiller
	expiresAt int64
	ttl       int64

	data interface{}
	err  error

	sync.Mutex
}

func (ci *cacheItem) Value() (data interface{}, err error) {
	ts := time.Now().Unix()
	ci.Lock()
	if ci.expiresAt == 0 || ci.expiresAt < ts {
		ci.data, ci.err = ci.fn()
		ci.expiresAt = ts + ci.ttl
	}
	data, err = ci.data, ci.err
	ci.Unlock()
	return
}
