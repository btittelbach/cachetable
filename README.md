### CacheTable

A quick and naive golang implementation of a CacheTable as described in Fabian "ryg" Giesen's (blog)[https://fgiesen.wordpress.com/2019/02/11/cache-tables/]:

Based on the naive hashmap code of Prakhar Srivastav.

### Usage

```go
package main

import (
    "fmt"
    "github.com/btittelbach/cachetable"
)

func main() {
    h, _ := cachetable.NewCacheTable(100,20) // create the cachetable
    keys := []string{"alpha", "beta", "charlie", "gamma", "delta"}

    // add the keys
    for _, key := range keys {
        h.Set(key, len(key))
    }

    fmt.Println("The load factor is:", h.Load())

    // retrieve the keys
    for _, key := range keys {
        val, present := h.Get(key)
        if present {
            fmt.Println("Key:", key, "->", "Value:", val.Value.(int))
        } else {
            fmt.Println(key, "is not present")
        }
    }

    // delete a key
    fmt.Println(h.Delete("alpha"))
    _, present := h.Get("alpha")
    if present {
        fmt.Println("The key's still there")
    } else {
        fmt.Println("Value associated with alpha deleted")
    }
}
```
