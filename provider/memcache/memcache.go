// The memcache package offers Memcached as a backing store for a cache. It
// uses the gomemcache library to do the talking under the hood.
//
// This package isn't intended for direct usage, it gets called only
// to register itself as a possible cache provider/
//
// Example
//   import (
//     github.com/tim-mit/cache
//     _ github.com/tim-mit/cache/provider/memcache
//   )
//
//   store, err := cache.New("memcache://127.0.0.1:11211?timeout=5")
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
package memcache

import (
	"github.com/tim-mit/cache"
	"github.com/tim-mit/cache/provider"
	"github.com/bradfitz/gomemcache/memcache"
	"fmt"
	"time"
	"strconv"
)

// validate provider.Provider interface satisfied
var _ provider.Provider = (*memcacheProvider)(nil)

// Container for Memcached connection details passed when
// the provider is initialised.
type connParams struct {
	host string
	timeout time.Duration
}

// Base type for the Memcached cache provider. Holds references to
// the connection details and any opened connections to a Memcached
// instance.
type memcacheProvider struct {
	client *memcache.Client
	params *connParams
}

// Unpacks and validates the initialisation uri parameters when
// a new Memcached Cache instance is created.
func (d *memcacheProvider) parseDetails(host string, params map[string][]string) (error) {

	d.params = &connParams{
		host: host,
	}
	
	t, ok := params["timeout"]
	if ok {
		timeout, err := strconv.Atoi(t[0])
		if err != nil {
			return fmt.Errorf("cache::memcache: Invalid timeout specification")
		}
		d.params.timeout = time.Duration(timeout) * time.Second
	}
	
	return nil
}

// Initialisation routine for plugin. Opens a connection to the
// Memcached instance provided.
func (d *memcacheProvider) Initialise(host string, params map[string][]string) (error) {
	err := d.parseDetails(host, params)
	
	if err != nil {
		return err
	}
	
	d.client = memcache.New(d.params.host)
	d.client.Timeout = d.params.timeout
	
	return nil
}

// Stores supplied data in the Memcached instance referenced by the given name. Data
// will not be returned to any subsequent Get requests after the expiry time
// has elapsed.
func (d *memcacheProvider) Set(name string, data interface{}, expiry time.Duration) (error) {

	byteData, err := (&provider.Result{data}).Bytes()

	err = d.client.Set(&memcache.Item{
		Key: name,
		Value: byteData,
	})
	
	if err != nil {
		return fmt.Errorf("cache::memcache: error during set -", err)
	}
	
	return nil
}

// Retrieves any data referenced by name. If the data has expired
// it will not be returned.
func (d *memcacheProvider) Get(name string) (*provider.Result) {
	val, err := d.client.Get(name)
	if err != nil {
		return &provider.Result{
			provider.Error{fmt.Errorf("cache::memcache: error during get -", err)},
		}
	}
	
	return &provider.Result{val.Value}
}

// Auto-called function to register this plugin with the base cache
// package.
func init() {
	cache.Register("memcache", &memcacheProvider{})
}