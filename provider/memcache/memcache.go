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

type memcacheProvider struct {
	client *memcache.Client
	params *connParams
}

type connParams struct {
	host string
	timeout time.Duration
}

func (d *memcacheProvider) parseDetails(host string, params map[string][]string) (error) {

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

func (d *memcacheProvider) Initialise(host string, params map[string][]string) (error) {
	err := d.parseDetails(host, params)
	
	if err != nil {
		return err
	}
	
	d.client = memcache.New(d.params.host)
	d.client.Timeout = d.params.timeout
	
	return nil
}

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

func (d *memcacheProvider) Get(name string) (*provider.Result) {
	val, err := d.client.Get(name)
	if err != nil {
		return &provider.Result{
			provider.Error{fmt.Errorf("cache::memcache: error during get -", err)},
		}
	}
	
	return &provider.Result{val.Value}
}

func init() {
	cache.Register("memcache", &memcacheProvider{})
}