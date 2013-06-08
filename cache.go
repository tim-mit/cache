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

type Cache struct {
	provider provider.Provider
}

func New(connUri string) (*Cache, error) {

	url, err := url.Parse(connUri)
	if err != nil {
		return nil, fmt.Errorf("cache: invalid connection uri", connUri)
	}

	d, ok := providers[url.Scheme]
	if !ok {
		return nil, fmt.Errorf("cache: unknown provider %q (forgotten import?)", url.Scheme)
	}
	
	cache := &Cache{
		provider: d,
	}
	
	cache.provider.Initialise(url.Host, url.Query())
	
	return cache, nil	
}

func (c *Cache) Set(name string, data interface{}, expiry time.Duration) (error) {
	err := c.provider.Set(name, data, expiry)
	if err != nil {
		return fmt.Errorf("cache: failure storing data %q", err)
	}
	
	return nil
} 

func (c *Cache) Get(name string) (*provider.Result) {
	return c.provider.Get(name)
}