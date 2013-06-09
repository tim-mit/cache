package provider

import (
	"time"
	"fmt"
	"reflect"
)

// Container for any errors thrown by the plugin interfaces.
type Error struct {
	Value error
}

// Container for the data stored in a Cache item.
type Result struct {
	Value interface{}
}

// Interface specification each plugin must satisfy.
type Provider interface {
	Initialise(host string, params map[string][]string) (error)
	Set(name string, data interface{}, expiry time.Duration) (error)
	Get(name string) (*Result)
}

// Utility function which given a Result will attempt to convert the
// Value stored within to a string.
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

// Utility function which given a Result will attempt to convert the
// Value stored within to a byte slice.
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