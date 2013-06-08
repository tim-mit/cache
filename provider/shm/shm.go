package shm

import (
	"github.com/tim-mit/cache"
	"github.com/tim-mit/cache/provider"
	"time"
	"sync"
)

// validate provider.Provider interface satisfied
var _ provider.Provider = (*shmProvider)(nil)

type item struct {
	data interface{}
	expiry *time.Time
}

type backingStore struct {
	sync.Mutex
	store map[string]*item
}

type connParams struct {
	host string
	timeout time.Duration
}

type shmProvider struct {
	data *backingStore
	params *connParams
}

func (d *shmProvider) cleaner() {
	// walk the map and remove any items older than their expiry
	
	// call GC to actually have it flushed?
}

func (d *shmProvider) Initialise(_ string, _ map[string][]string) (error) {
	d.data = &backingStore{
		store: map[string]*item{},
	}
	
	return nil
}

func (d *shmProvider) Set(name string, data interface{}, expiry time.Duration) (error) {
	var exp *time.Time
	if expiry > -1 {
		t := time.Now().Add(expiry)
		exp = &t
	}
	
	d.data.Lock()
	 
	d.data.store[name] = &item{
		data: data,
		expiry: exp,
	}
	
	d.data.Unlock()
	
	return nil
}

func (d *shmProvider) Get(name string) (*provider.Result) {
	val := d.data.store[name]
	
	if val.expiry.Before(time.Now()) {
		return &provider.Result{nil}
	}
	
	return &provider.Result{val.data}
}

func init() {
	cache.Register("shm", &shmProvider{})
}