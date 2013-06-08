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

type connParams struct {
	host string
	timeout time.Duration
}

type redisProvider struct {
	conn redis.Conn
	params *connParams
}

func (d *redisProvider) parseDetails(host string, params map[string][]string) (error) {

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

func (d *redisProvider) Set(name string, data interface{}, expiry time.Duration) (error) {
	_, err := d.conn.Do("set", name, data)
	if err != nil {
		return fmt.Errorf("cache::redis: error during set -", err)
	}
	
	return nil
}

func (d *redisProvider) Get(name string) (*provider.Result) {
	val, err := d.conn.Do("get", name)
	if err != nil {
		return &provider.Result{
			provider.Error{fmt.Errorf("cache::redis: error during get -", err)},
		}
	}
	
	return &provider.Result{val}
}

func init() {
	cache.Register("redis", &redisProvider{})
}