package provider

import (
	"time"
	"fmt"
	"reflect"
)

type Error struct {
	Value error
}

type Result struct {
	Value interface{}
}

type Provider interface {
	Initialise(host string, params map[string][]string) (error)
	Set(name string, data interface{}, expiry time.Duration) (error)
	Get(name string) (*Result)
}

func (i *Result) String() (string, error) {

	val := i.Value
	
	switch val := val.(type) {
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	case Error:
		return "", val.Value
	}
	return "", fmt.Errorf("cache::provider: unable to convert %v to string", reflect.TypeOf(val))
}

func (i *Result) Bytes() ([]byte, error) {

	val := i.Value
	
	switch val := val.(type) {
	case string:
		return []byte(val), nil
	case []byte:
		return val, nil
	case Error:
		return nil, val.Value
	}
	return nil, fmt.Errorf("cache::provider: unable to convert %v to []byte", reflect.TypeOf(val))
}