package sgcache

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkCache(b *testing.B) {
	cache := New(time.Second, 900000000)
	defer cache.Close()

	for n := 0; n < b.N; n++ {
		cache.Set(fmt.Sprint(n%1000000), "value", time.Second*3)
		cache.Get(fmt.Sprint(n % 1000000))
	}
}

func BenchmarkConcurrentGet(b *testing.B) {
	b.StopTimer()
	n := 10000
	c := New(time.Duration(time.Second), 900000000)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "foo" + strconv.Itoa(i)
		keys[i] = k
		c.Set(k, "bar", time.Second*3)
	}
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(k string) {
			for j := 0; j < each; j++ {
				c.Get("foo")
			}
			wg.Done()
		}(v)
	}
	b.StartTimer()
	wg.Wait()
}
