package genericcache

import (
	"time"

	"github.com/alphadose/haxmap"
)

func NewCache[K hashable, V any](lifetime time.Duration) *Cache[K, V] {
	cache := &Cache[K, V]{
		cache:    haxmap.New[K, wrapper[V]](),
		lifetime: lifetime,
	}

	if lifetime != NoExpiration {
		go func() {
			wipe := time.NewTicker(lifetime)
			for {
				select {
				case <-wipe.C:
					cache.ForEach(func(k K, v wrapper[V]) bool {
						if v.ExpiresAt <= time.Now().Unix() {
							cache.Del(k)
						}

						return true
					})
				}
			}
		}()
	}

	return cache
}

func (c *Cache[K, V]) Get(key K) (value V, err error) {
	v, ok := c.cache.Get(key)
	if !ok {
		err = ErrorKey
		return
	}
	if v.ExpiresAt <= time.Now().Unix() {
		err = ErrorExpired
		return
	}
	return v.Value, nil
}

func (c *Cache[K, V]) Set(key K, value V) {
	if c.lifetime != NoExpiration {
		c.cache.Set(key, wrapper[V]{
			Value:     value,
			ExpiresAt: time.Now().Add(c.lifetime).Unix(),
		})
	} else {
		c.cache.Set(key, wrapper[V]{
			Value:     value,
			ExpiresAt: NoExpiration,
		})
	}
}

func (c *Cache[K, V]) Swap(key K, newValue V) (oldValue V, swapped bool) {
	var expiresAt int64

	if c.lifetime != NoExpiration {
		expiresAt = time.Now().Add(c.lifetime).Unix()
	} else {
		expiresAt = NoExpiration
	}

	o, s := c.cache.Swap(key, wrapper[V]{
		Value:     newValue,
		ExpiresAt: expiresAt,
	})

	return o.Value, s
}

func (c *Cache[K, V]) Del(key K) {
	c.cache.Del(key)
}

func (c *Cache[K, V]) ForEach(f func(K, wrapper[V]) bool) {
	c.cache.ForEach(f)
}
