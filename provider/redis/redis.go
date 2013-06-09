// The redis package offers Redis as a backing store for a cache. It
// uses the redigo library to do the talking under the hood.
//
// This package isn't intended for direct usage, it gets called only
// to register itself as a possible cache provider/
//
// Example
//   import (
//     github.com/tim-mit/cache
//     _ github.com/tim-mit/cache/provider/redis
//   )
//
//   store, err := cache.New("redis://127.0.0.1:6379?timeout=5")
//   if err != nil {
//      // handle error
//   }
//
//   err = store.Set("app.key", "Data to store")
//
//   [snip]
//
//   msg, err := store.Get("app.key")
//
package redis

import (
	"github.com/tim-mit/cache"
	"github.com/tim-mit/cache/provider"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"time"
	"strconv"
)

// validate provider.Provider interface satisfied
var _ provider.Provider = (*redisProvider)(nil)

// Container for Redis connection details passed when
// the provider is initialised.
type connParams struct {
	host string
	timeout time.Duration
}

// Base type for the Redis cache provider. Holds references to
// the connection details and any opened connections to a Redis
// instance.
type redisProvider struct {
	conn redis.Conn
	params *connParams
}

// Unpacks and validates the initialisation uri parameters when
// a new Redis Cache instance is created.
func (d *redisProvider) parseDetails(host string, params map[string][]string) (error) {

	d.params = &connParams{
		host: host,
	}
	
	t, ok := params["timeout"]
	if ok {
		timeout, err := strconv.Atoi(t[0])
		if err != nil {
			return fmt.Errorf("cache::redis: Invalid timeout specification")
		}
		d.params.timeout = time.Duration(timeout) * time.Second
	}
	
	return nil
}

// Initialisation routine for plugin. Opens a connection to the
// Redis instance provided.
func (d *redisProvider) Initialise(host string, params map[string][]string) (error) {

	err := d.parseDetails(host, params)
	
	if err != nil {
		return err
	}

	// TODO conn, read and write timeouts should all be separate items
	//      need to handle one and not the others being set though
	c, err := redis.DialTimeout("tcp", d.params.host, d.params.timeout, d.params.timeout, d.params.timeout)
	
	if err != nil {
		return fmt.Errorf("cache::redis: error during connection - %q", err)
	}
	
	d.conn = c
	return nil
}

// Stores supplied data in the Redis instance referenced by the given name. Data
// will not be returned to any subsequent Get requests after the expiry time
// has elapsed.
func (d *redisProvider) Set(name string, data interface{}, expiry time.Duration) (error) {
	_, err := d.conn.Do("set", name, data)
	if err != nil {
		return fmt.Errorf("cache::redis: Error during set %q", err)
	}
	
	return nil
}

// Retrieves any data referenced by name. If the data has expired
// it will not be returned.
func (d *redisProvider) Get(name string) (*provider.Result) {
	val, err := d.conn.Do("get", name)
	if err != nil {
		return &provider.Result{
			provider.Error{fmt.Errorf("cache::redis: Error during get %q", err)},
		}
	}
	
	return &provider.Result{val}
}

// Auto-called function to register this plugin with the base cache
// package.
func init() {
	cache.Register("redis", &redisProvider{})
}