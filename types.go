package genericcache

import (
	"time"
	"unsafe"

	"github.com/alphadose/haxmap"
	"golang.org/x/exp/constraints"
)

type Cache[K hashable, V any] struct {
	cache    *haxmap.Map[K, wrapper[V]]
	lifetime time.Duration
}

const NoExpiration = 0

type Error string

func (e Error) Error() string { return string(e) }

const ErrorKey = Error("key does not exist in cache")
const ErrorExpired = Error("key exists but it's expired")

type hashable interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string | uintptr | ~unsafe.Pointer
}

type wrapper[V any] struct {
	Value     V
	ExpiresAt int64
}
