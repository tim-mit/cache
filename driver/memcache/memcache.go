package memcache

import (
	"github.com/tim-mit/cache"
	"github.com/tim-mit/cache/driver"
	"github.com/bradfitz/gomemcache/memcache"
	"fmt"
	"time"
	"log"
	"strconv"
)

// validate driver.Driver interface satisfied
var _ driver.Driver = (*memcacheDriver)(nil)

type memcacheDriver struct {
	client *memcache.Client
	params *connParams
}

type connParams struct {
	host string
	timeout time.Duration
}

func (d *memcacheDriver) parseDetails(host string, params map[string][]string) (error) {

	d.params = &connParams{
		host: host,
	}
	
	t, ok := params["timeout"]
	if ok {
		timeout, err := strconv.Atoi(t[0])
		if err != nil {
			return fmt.Errorf("cache::redis: invalid timeout specification")
		}
		d.params.timeout = time.Duration(timeout) * time.Second
	}
	
	return nil
}

func (d *memcacheDriver) Initialise(host string, params map[string][]string) (error) {
	err := d.parseDetails(host, params)
	
	if err != nil {
		return err
	}
	
	d.client = memcache.New(d.params.host)
	d.client.Timeout = d.params.timeout
	
	return nil
}

func (d *memcacheDriver) Set(name string, data interface{}, expiry time.Duration) (error) {

	byteData, err := (&driver.Result{data}).Bytes()

	err = d.client.Set(&memcache.Item{
		Key: name,
		Value: byteData,
	})
	
	if err != nil {
		return fmt.Errorf("cache::memcache: error during set -", err)
	}
	
	return nil
}

func (d *memcacheDriver) Get(name string) (*driver.Result) {
	log.Println("getting")
	val, err := d.client.Get(name)
	if err != nil {
		return &driver.Result{
			driver.Error{fmt.Errorf("cache::memcache: error during get -", err)},
		}
	}
	
	log.Println("returning data")
	return &driver.Result{val.Value}
}

func init() {
	cache.Register("memcache", &memcacheDriver{})
}