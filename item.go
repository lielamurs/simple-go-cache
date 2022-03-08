package sgcache

import (
	"time"
)

type Item struct {
	data interface{}
	ttl  time.Time
}

func (item *Item) expired() bool {
	return item.ttl.Before(time.Now())
}
