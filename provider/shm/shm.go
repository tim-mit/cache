// The shm package is a very simplistic cache provider intended for
// applications that aren't going to rely too heavily on cache. The
// data is stored in the application heap space as a string indexed
// hash map.
//
// This package isn't intended for direct usage, it gets called only
// to register itself as a possible cache provider/
//
// Example
//   import (
//     github.com/tim-mit/cache
//     _ github.com/tim-mit/cache/provider/shm
//   )
//
//   store, err := cache.New("shm://")
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
package shm

import (
	"github.com/tim-mit/cache"
	"github.com/tim-mit/cache/provider"
	"time"
	"sync"
)

// validate provider.Provider interface satisfied
var _ provider.Provider = (*shmProvider)(nil)

// Container to hold any cached items within. Allows the item
// expiry to be stored alongside for validation and cleanup.
type item struct {
	data interface{}
	expiry *time.Time
}

// Container to hold the cache item map and a mutex to marshall
// write access by callees.
type backingStore struct {
	sync.Mutex
	store map[string]*item
}
 
// The base type for this plugin. Holds a pointer to
// any initialised backingStore.
type shmProvider struct {
	data *backingStore
}

// Internal functionality to ensure expired items are
// purged from the backingStore holding them.
// TODO - implement
func (d *shmProvider) cleaner() {
	// walk the map and remove any items older than their expiry
	
	// call GC to actually have it flushed?
}

// Initialisation routine for plugin. Creates a backingStore
// to hold any subsequent items added.
func (d *shmProvider) Initialise(_ string, _ map[string][]string) (error) {
	d.data = &backingStore{
		store: map[string]*item{},
	}
	
	return nil
}

// Stores supplied data in the backingStore referenced by the given name. Data
// will not be returned to any subsequent Get requests after the expiry time
// has elapsed.
func (d *shmProvider) Set(name string, data interface{}, expiry time.Duration) (error) {
	var exp *time.Time
	if expiry > -1 {
		t := time.Now().Add(expiry)
		exp = &t
	}
	
	// The map data type does not have any concurrent access
	// guarantees so we lock here while making changes.
	d.data.Lock()
	 
	d.data.store[name] = &item{
		data: data,
		expiry: exp,
	}
	
	d.data.Unlock()
	
	return nil
}

// Retrieves any data referenced by name. If the data has expired
// it will not be returned.
func (d *shmProvider) Get(name string) (*provider.Result) {
	val := d.data.store[name]
	
	if val.expiry.Before(time.Now()) {
		return &provider.Result{nil}
	}
	
	return &provider.Result{val.data}
}

// Auto-called function to register this plugin with the base cache
// package.
func init() {
	cache.Register("shm", &shmProvider{})
}