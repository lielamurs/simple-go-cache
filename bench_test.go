package sgcache

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkCache(b *testing.B) {
	cache := New(time.Second, 50, 5000000)
	defer cache.Close()

	for n := 0; n < b.N; n++ {
		cache.Set(fmt.Sprint(n%1000000), "value", time.Second*3)
		cache.Get(fmt.Sprint(n % 1000000))
	}
}
