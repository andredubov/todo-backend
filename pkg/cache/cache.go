package cache

import "time"

type Cache interface {
	Set(key, value interface{}, ttl time.Duration) error
	Get(key interface{}) (interface{}, error)
}
