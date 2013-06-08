#cache

A generic cache client interface for [Go](http://golang.org).

##Introduction

[cache](https://github.com/tim-mit/cache) is an attempt to create and build a unified API to many different cache
providers. Application code can use this library to insulate itself from being tied to one cache implementation.

A modular provider system (a la the core [database/sql](http://golang.org/pkg/database/sql/) package) is used so
any backing store that offers cache semantics can be utilised.

##Providers
Cache provider interfaces are currently available for
* [Memcached](http://memcached.org/)
* [Redis](http://redis.io)
* SHM (shared application memory, concurrent access safe)

##Example
Storing a string in [Redis](http://redis.io):
    store, err := cache.New("redis://127.0.0.1:6379?timeout=2")
    if err != nil {
      // handle error
    }
    store.Set("msg", "hello world!")

and later on retrieving it:
    store, err := cache.New("redis://127.0.0.1:6379?timeout=2")
    if err != nil {
      // handle error
    }
    val, err = store.Get("msg").String()
  
##Install
The core package can be installed using `go get`:

  go get github.com/tim-mit/cache
  
but to be useful you need to also install at least one of the following providers;

    go get github.com/tim-mit/cache/providers/memcache
    go get github.com/tim-mit/cache/providers/redis
    go get github.com/tim-mit/cache/providers/shm
  
The [Memcache](http://github.com/tim-mit/cache/providers/memcache) and [Redis](http://github.com/tim-mit/cache/providers/redis) providers have dependencies on the following two packages respectively

    go get github.com/bradfitz/gomemcache/memcache
    go get github.com/garyburd/redigo

##License
This package is available for use under the [MIT license](http://opensource.org/licenses/MIT).
