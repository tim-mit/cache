// cache is an attempt to create and build a unified API to many
// different cache providers. Application code can use this library
// to insulate itself from being tied to one cache implementation.
//
// A modular provider system (a la the core database/sql package)
// is used so any backing store that offers cache semantics can be
// utilised.
//
// A 'cache' in this implementation is taken to mean (and only mean)
// storing arbitrary blobs of data for subsequent retrieval.
//
// This is the base package and doesn't actually implement any cache
// of its own. This package needs to be used in conjunction with a
// provider package.
//
package cache

import (
	"github.com/tim-mit/cache/provider"
	"time"
	"fmt"
	"net/url"
)

var providers = make(map[string]provider.Provider)

// Register makes a cache provider available by the provided name.
// If Register is called twice with the same name or if provider is nil,
// it panics.
func Register(name string, provider provider.Provider) {
	if provider == nil {
		panic("cache: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("cache: Register called twice for provider " + name)
	}
	providers[name] = provider
}

// A Cache is the base type representing a store which can have items
// added to or read from.
type Cache struct {
	provider provider.Provider
}

// Creates and initialises a new Cache store of the type specified in the uri.
// The returned Cache pointer can then be used to interact with the store.
func New(cacheUri string) (*Cache, error) {

	url, err := url.Parse(cacheUri)
	if err != nil {
		return nil, fmt.Errorf("cache: Invalid cache initialisation uri", cacheUri)
	}

	d, ok := providers[url.Scheme]
	if !ok {
		return nil, fmt.Errorf("cache: Unknown provider %q (forgotten import?)", url.Scheme)
	}
	
	cache := &Cache{
		provider: d,
	}
	
	cache.provider.Initialise(url.Host, url.Query())
	
	return cache, nil	
}

// Store data in the Cache item referenced by the given name. The data should
// not be returned after the expiry time has elapsed.
func (c *Cache) Set(name string, data interface{}, expiry time.Duration) (error) {
	err := c.provider.Set(name, data, expiry)
	if err != nil {
		return fmt.Errorf("cache: Failure storing data %q", err)
	}
	
	return nil
} 

// Retrieve any data stored in the Cache and referenced by the given name.
func (c *Cache) Get(name string) (*provider.Result) {
	return c.provider.Get(name)
}