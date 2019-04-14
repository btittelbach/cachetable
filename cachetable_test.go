package cachetable

import (
	"strconv"
	"testing"
)

func TestNewCacheTableInit(t *testing.T) {
	cases := []struct {
		ina, inb int
	}{
		{10, 10},
		{10, 2},
		{0, 0},
		{1, 0},
		{1, 1},
	}

	mytest := func(prealloc bool) {
		for _, c := range cases {
			h, err := NewCacheTable(c.ina, c.inb, prealloc)
			if c.ina == 0 || c.inb == 0 {
				if err == nil {
					t.Errorf("Expected error, didn't get it")
				}
			} else {
				if h == nil || err != nil {
					t.Errorf("NewCacheTable(%d,%d) threw unexpected error: %s", c.ina, c.inb, err)
				}
				if h.numbuckets != c.ina || h.bucketcapacity != c.inb {
					t.Errorf("NewCacheTable(%d,%d) == %d,%d", c.ina, c.inb, h.numbuckets, h.bucketcapacity)
				}
			}
		}
	}
	mytest(false) //don't preallocate
	mytest(true)  //preallocated
}

func TestNewCacheTableCapacity(t *testing.T) {
	cases := []struct {
		ina, inb int
	}{
		{10, 10},
		{10, 2},
		{0, 0},
		{1, 0},
		{1, 1},
	}

	mytest := func(prealloc bool) {
		for _, c := range cases {
			h, err := NewCacheTable(c.ina, c.inb, prealloc)
			if c.ina == 0 || c.inb == 0 {
				if err == nil {
					t.Errorf("Expected error, didn't get it")
				}
			} else {
				if h.Capacity() != c.ina*c.inb {
					t.Errorf("NewCacheTable(%q,%q).Capacity() == %q, want %q", c.ina, c.inb, h.Capacity(), c.ina*c.inb)
				}
			}
		}
	}
	mytest(false) //don't preallocate
	mytest(true)  //preallocated
}

func TestLenAndLoad(t *testing.T) {
	cases := []struct {
		ina, inb, want int
	}{
		{10, 1, 10},
		{10, 2, 20},
	}

	for _, c := range cases {
		h, _ := NewCacheTable(c.ina, c.inb*10, true)
		for i := 1; i <= c.ina*c.inb; i++ {
			key := strconv.Itoa(i)
			h.Set(key, i)
		}
		got := h.Len()
		if got != c.want {
			t.Errorf("Len(%d,%d) == %d, want %d", c.ina, c.inb, got, c.want)
		}

		load := h.Load()
		want := float32(c.ina*c.inb) / float32(c.ina*c.inb*10)
		if load != want {
			t.Errorf("Load(%d) == %f, want %f", c.ina, load, want)
		}
	}
}

func TestGetAndSet(t *testing.T) {
	h, _ := NewCacheTable(10, 10, true)
	keys := []string{"alpha5", "beta4", "charlie7", "gamma_6", "delta__8"}

	// testing primitives
	for _, key := range keys {
		h.Set(key, len(key))
	}

	for _, key := range keys {
		got, _ := h.Get(key)
		want := len(key)
		if got.Value.(int) != len(key) {
			t.Errorf("want: %q, got: %q", want, got)
		}
	}

	// testing strings
	for _, key := range keys {
		h.Set(key, key+key)
	}

	for _, key := range keys {
		got, _ := h.Get(key)
		want := key + key
		if got.Value.(string) != want {
			t.Errorf("want: %q, got: %q", want, got)
		}
	}

	// testing references to compound types
	arr := []int{2, 3, 4}
	h.Set("myArray", arr)
	a, _ := h.Get("myArray")
	k := a.Value.([]int)
	k[0] = 100
	if k[0] != arr[0] {
		t.Errorf("Reference has not been mutated")
	}
}

func TestCollisionsAndMemoryConstrain(t *testing.T) {
	// a small cachetable that is bound to have collisions
	numbuckets := 2
	bucketcap := 2
	h, _ := NewCacheTable(numbuckets, bucketcap, true)

	keys := []string{"alpha5", "beta4", "charlie7", "gamma_6", "delta__8"}

	if h.Capacity() != numbuckets*bucketcap {
		t.Errorf("Wrong Capacity")
	}

	for _, key := range keys {
		h.Set(key, len(key))
		if h.Len() > numbuckets*bucketcap {
			t.Errorf("Constraint did not work, buckets magically increased")
		}
	}

	number_of_items_recovered := 0
	for _, key := range keys {
		got, inmap := h.Get(key)
		want := len(key)
		if inmap && got.Value.(int) == want {
			number_of_items_recovered++
		}
	}
	if number_of_items_recovered != numbuckets*bucketcap {
		t.Errorf("expected to recover %q elements but got %q", numbuckets*bucketcap, number_of_items_recovered)
	}
}

func TestOverwrite(t *testing.T) {
	// a cachetable with just one elem
	h, _ := NewCacheTable(1, 1, true)
	keya := "alpha"
	keyb := "beta"
	vala := 10
	valb := 20
	var status bool

	h.Set(keya, vala)

	// verify alpha is there
	got, inmap := h.Get(keya)
	if inmap == false || got.Value.(int) != vala {
		t.Errorf("Element just added not found!")
	}

	// add beta now
	status = h.Set(keyb, valb)
	if !status {
		t.Errorf("Unable to add element")
	}

	// verify beta is there
	got, inmap = h.Get(keyb)
	if inmap == false || got.Value.(int) != valb {
		t.Errorf("Element just added not found!")
	}

	// verify alpha is gone
	got, inmap = h.Get(keya)
	if inmap == false || got != nil {
		t.Errorf("Found element we should have overwritten!")
	}
}

func TestDelete(t *testing.T) {
	// a cachetable with just one elem
	h, _ := NewCacheTable(1, 1, true)
	var status bool
	keya := "alpha"
	keyb := "beta"
	vala := 10
	valb := 20
	h.Set(keya, vala)

	// verify it's there
	got, inmap := h.Get(keya)
	if inmap == false || got.Value.(int) != vala {
		t.Errorf("Element just added not found!")
	}

	// lets delete it
	_, status = h.Delete(keya)
	if !status {
		t.Errorf("Unable to delete")
	}

	// verify it's gone
	got, inmap = h.Get(keya)
	if inmap == false || got != nil {
		t.Errorf("Found element we just deleted!")
	}

	// add beta now
	status = h.Set(keyb, valb)
	if !status {
		t.Errorf("Unable to add element")
	}

	// lastly, lets delete a non-existent element
	_, status = h.Delete("gamma")
	if status {
		t.Errorf("Deleted a missing key")
	}
}
