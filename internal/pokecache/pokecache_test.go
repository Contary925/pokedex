package pokecache

import (
	"testing"
	"time"
)

func TestReap(t *testing.T) {
	c := NewCache(50 * time.Millisecond)
	c.Add("key1", []byte("val1"))
	if res, _ := c.Get("key1"); string(res) != "val1" {
		t.Errorf("Value %v written incorrectly as %v", "val1", string(res))
	}
	time.Sleep(30 * time.Millisecond)
	c.Add("key2", []byte("val2"))
	if res, _ := c.Get("key2"); string(res) != "val2" {
		t.Errorf("Value %v written incorrectly as %v", "val2", string(res))
	}
	time.Sleep(30 * time.Millisecond)
	if _, ok := c.Get("key1"); ok {
		t.Errorf("Entry key1, val1 still in cache despite expiring")
	}
	if _, ok := c.Get("key2"); !ok {
		t.Errorf("Entry key2, val2 gone from cache despite not expiring")
	}
	time.Sleep(50 * time.Millisecond)
	if _, ok := c.Get("key1"); ok {
		t.Errorf("Entry key1, val1 still in cache despite expiring")
	}
	if _, ok := c.Get("key2"); ok {
		t.Errorf("Entry key2, val2 still in cache despite expiring")
	}
}