package cachemap

import "errors"

// Node which is stored at each level
type Node struct {
	key         string
	Value       interface{}
	create_time uint64
}

// HashMap implemented with a fixed bucketsize.
// // NOPE: Uses chaining to resolve collisions.
// removes oldest element in bucket
type HashMap struct {
	current_time uint64
	bucketsize   int
	numbuckets   int
	count        int
	buckets      [][]Node
}

/** PRIVATE METHODS **/

// returns the index at which the key needs to go
func (h *HashMap) getIndex(key string) int {
	return int(hash(key)) % h.numbuckets
}

// Implements the Jenkins hash function
func hash(key string) uint32 {
	var h uint32
	for _, c := range key {
		h += uint32(c)
		h += (h << 10)
		h ^= (h >> 6)
	}
	h += (h << 3)
	h ^= (h >> 11)
	h += (h << 15)
	return h
}

/** PUBLIC METHODS **/

// Len returns the count of the elements in the hashmap
func (h *HashMap) Len() int {
	return h.count
}

// Size returns the bucket size of the hashamp
func (h *HashMap) BucketSize() int {
	return h.bucketsize
}

// NewCacheMap is the constuctor that returns a new HashMap of a fixed size
// returns an error when a size of 0 is provided
func NewCacheMap(numbuckets, bucketsize int) (*HashMap, error) {
	h := new(HashMap)
	if bucketsize < 1 {
		return h, errors.New("bucketsize of hashmap has to be > 1")
	}
	h.bucketsize = bucketsize
	h.numbuckets = numbuckets
	h.count = 0
	h.current_time = 0
	h.buckets = make([][]Node, numbuckets)
	for i := range h.buckets {
		h.buckets[i] = make([]Node, 0)
	}
	return h, nil
}

// Get returns the value associated with a key in the hashmap,
// and an error indicating whether the value exists
func (h *HashMap) Get(key string) (*Node, bool) {
	index := h.getIndex(key)
	chain := h.buckets[index]
	for _, node := range chain {
		if node.key == key {
			return &node, true
		}
	}
	return nil, false
}

// Set the value for an associated key in the hashmap
func (h *HashMap) Set(key string, value interface{}) bool {
	index := h.getIndex(key)
	chain := h.buckets[index]

	// first see if the key already exists
	// also find oldest element in case it does not
	oldest_time := uint64(1<<64 - 1)
	oldest_index := 0
	for i := range chain {
		// if found, update the value
		node := &chain[i]
		if node.key == key {
			node.Value = value
			return true
		}
		if node.create_time < oldest_time {
			oldest_index = i
			oldest_time = node.create_time
		}
	}

	// if key doesn't exist, add it to the hashmap
	newnode := Node{key: key, Value: value, create_time: h.current_time}
	h.current_time++ //increment cachemap insert time

	//before the roll-over, compress time
	if h.current_time == 1<<64-1 {
		h.CompressTime()
		newnode.create_time = h.current_time
	}

	// it bucket is full, overwrite oldest element
	if len(chain) >= h.bucketsize {
		chain[oldest_index] = newnode
	} else {
		// there's enough space, let's append the node
		chain = append(chain, newnode)
		h.buckets[index] = chain
		h.count++
	}
	return true
}

// Delete the value associated with key in the hashmap
func (h *HashMap) Delete(key string) (*Node, bool) {
	index := h.getIndex(key)
	chain := h.buckets[index]

	found := false
	var location int
	var mapNode *Node

	// start a search for the key
	for loc, node := range chain {
		if node.key == key {
			found = true
			location = loc
			mapNode = &node
		}
	}

	// if found delete the elem from the slice
	if found {
		h.count--
		N := len(chain) - 1
		chain[location], chain[N] = chain[N], chain[location]
		chain = chain[:N]
		h.buckets[index] = chain
		return mapNode, true
	}

	// if not found return false
	return nil, false
}

// Load returns the load factor of the hashmap
func (h *HashMap) Load() float32 {
	return float32(h.count) / float32(h.bucketsize*h.numbuckets)
}

func (h *HashMap) CompressTime() {
	//TODO
	//put pointers to all nodes into a big list of size h.size
	//sort list by create-time
	//give each node a new incremental number
	//set h.current_time to next free number
}
