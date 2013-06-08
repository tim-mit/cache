package shm

import (
	"github.com/tim-mit/cache"
	"github.com/tim-mit/cache/driver"
	"time"
	"sync"
)

// validate driver.Driver interface satisfied
var _ driver.Driver = (*shmDriver)(nil)

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

type shmDriver struct {
	data *backingStore
	params *connParams
}

func (d *shmDriver) cleaner() {
	// walk the map and remove any items older than their expiry
	
	// call GC to actually have it flushed?
}

func (d *shmDriver) Initialise(_ string, _ map[string][]string) (error) {
	d.data = &backingStore{
		store: map[string]*item{},
	}
	
	return nil
}

func (d *shmDriver) Set(name string, data interface{}, expiry time.Duration) (error) {
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

func (d *shmDriver) Get(name string) (*driver.Result) {
	val := d.data.store[name]
	
	if val.expiry.Before(time.Now()) {
		return &driver.Result{nil}
	}
	
	return &driver.Result{val.data}
}

func init() {
	cache.Register("shm", &shmDriver{})
}