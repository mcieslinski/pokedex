package pkcache

import (
  "bytes"
  "testing"
  "time"
)

func TestPkCache(t *testing.T) {
  //remote := map[string][]byte {
  //  "one":   []byte("Hello"),
  //  "two":   []byte("World"),
  //  "three": []byte("I"),
  //  "four":  []byte("am"),
  //  "five":  []byte("Iron"),
  //  "six":   []byte("Man"),
  //}

  const (
    onekey = "one"
    twokey = "two"
    oneval = "Hello"
    twoval = "World"
  )

  cache := NewPkCache(5 * time.Second)
  cache.Add(onekey, []byte(oneval))
  time.Sleep(500 * time.Millisecond)
  cache.Add(twokey, []byte(twoval))

  // Test first cache value
  if v, ok := cache.Get(onekey); !ok {
    t.Errorf("Cache get failed for \"%v\", but should have succeeded.", onekey)
  } else if !bytes.Equal(v, []byte(oneval)) {
    t.Errorf("Value in cache (\"%v\") is not as expected *\"%v\")", v, oneval)
  }

  // Test second cache value
  if v, ok := cache.Get(twokey); !ok {
    t.Errorf("Cache get failed for \"%v\", but should have suceeded.", twokey)
  } else if !bytes.Equal(v, []byte(twoval)) {
    t.Errorf("Value in cache (\"%v\") is not as expected *\"%v\")", v, twoval)
  }

  // First entry should be gone (give 5 ms for leeway)
  time.Sleep(5 * time.Second)
  if _, ok := cache.Get(onekey); ok {
    t.Errorf("Found \"%v\" in cache but should have been deleted.", onekey)
  }
  if _, ok := cache.Get(twokey); !ok {
    t.Errorf("Didn't found \"%v\" in cache but it should still be there.", twokey)
  }
}
