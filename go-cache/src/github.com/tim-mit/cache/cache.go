package cache

import (
	"github.com/tim-mit/cache/driver"
	"time"
	"fmt"
	"net/url"
)

var drivers = make(map[string]driver.Driver)

// Register makes a cache driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver driver.Driver) {
	if driver == nil {
		panic("cache: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("cache: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

type Cache struct {
	driver driver.Driver
}

func New(connUri string) (*Cache, error) {

	url, err := url.Parse(connUri)
	if err != nil {
		return nil, fmt.Errorf("cache: invalid connection uri", connUri)
	}

	d, ok := drivers[url.Scheme]
	if !ok {
		return nil, fmt.Errorf("cache: unknown driver %q (forgotten import?)", url.Scheme)
	}
	
	cache := &Cache{
		driver: d,
	}
	
	cache.driver.Initialise(url.Host, url.Query())
	
	return cache, nil	
}

func (c *Cache) Set(name string, data interface{}, expiry time.Duration) (error) {
	err := c.driver.Set(name, data, expiry)
	if err != nil {
		return fmt.Errorf("cache: failure storing data %q", err)
	}
	
	return nil
} 

func (c *Cache) Get(name string) (*driver.Result) {
	return c.driver.Get(name)
}
