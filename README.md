# Simple go cache - in memory cache with expiration


#### Usage
```go
import (
  "time"
  "github.com/lielamurs/simple-go-cache"
)

func main () {
  cache := sgcache.New(time.Duration(time.Second), 5000)
  cache.Set("foo", "bar", time.Minute)
  value, exists := cache.Get("foo")
  cache.Delete("foo")
}
```