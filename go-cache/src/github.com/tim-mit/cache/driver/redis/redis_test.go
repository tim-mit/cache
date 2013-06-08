package redis

import (
    "testing"
    "github.com/tim-mit/cache"
	"time"
	"bytes"
)

func chkerr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedis(t *testing.T) {

	store, err := cache.New("redis://127.0.0.1:6379?timeout=5")
	chkerr(t, err)
	t.Log("creating cache store passes")
	
	testString := "TEST STRING"
	err = store.Set("test-string", testString, time.Duration(10) * time.Second)
	chkerr(t, err)
	t.Log("store of test string passes")
	
	val, err := store.Get("test-string").String()
	
	chkerr(t, err)
	if val != testString {
		t.Log("retrieving test string fails due to mismatch ( [orig]", testString, " != [retrieved]", val, ")")
		t.Fail()
	} else {
		t.Log("retrieving test string passes")
	}
	
	testByteArray := []byte("test byte array")
	err = store.Set("test-byte-array", testByteArray, time.Duration(10) * time.Second)
	chkerr(t, err)
	t.Log("store of test byte array passes")
	
	val2, err := store.Get("test-byte-array").Bytes()
	chkerr(t, err)
	
	if !bytes.Equal(testByteArray, val2) {
		t.Log("retrieving test byte array fails due to mismatch ( [orig]", val2, " != [retrieved]", testByteArray, ")")
		t.Fail()
	} else {
		t.Log("retrieving test byte array passes")
	}
}

