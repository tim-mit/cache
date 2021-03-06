package shm

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

func TestShm(t *testing.T) {

	store, err := cache.New("shm://")
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
	
	val3, err := store.Get("test-string-not-set").String()
	chkerr(t, err)
	
	if val3 != "" {
		t.Log("non-nil response to request for key with no content")
		t.Fail()
	}
}