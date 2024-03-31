package inmemdb

import (
	"sync"
	"time"
)

// Cache represents a generic caching interface for key-value pairs.
// Different cache implementations can be created by implementing this interface.
type Cache[K comparable, V any] interface {
	// Get retrieves the value associated with the given key from the cache.
	// It returns the value and true if the key exists in the cache, otherwise, it returns the zero value of the value type and false.
	Get(K) (V, bool)

	// Set adds or updates the key-value pair in the cache without setting an expiration time.
	// If the key already exists, its value will be overwritten with the new value.
	Set(K, V)

	// SetWithTimeout adds or updates the key-value pair in the cache with an expiration time.
	// If the timeout duration is zero or negative, the key-value pair will not have an expiration time.
	SetWithTimeout(K, V, time.Duration)

	// NotFoundSet adds a key-value pair to the cache if the key does not already exist and returns true.
	// Otherwise, it does nothing and returns false.
	NotFoundSet(K, V) bool

	// NotFoundSetWithTimeout adds a key-value pair to the cache with an expiration time if the key does not already exist and returns true.
	// If the timeout is zero or negative, the key-value pair will not have an expiration time.
	// If expiry is disabled, it behaves like NotFoundSet.
	NotFoundSetWithTimeout(K, V, time.Duration) bool

	// Delete removes the key-value pair associated with the given key from the cache.
	Delete(K)

	// TransferTo transfers all key-value pairs from the source cache to the provided destination cache.
	TransferTo(Cache[K, V])

	// CopyTo copies all key-value pairs from the source cache to the provided destination cache.
	CopyTo(Cache[K, V])

	// Keys returns a slice containing the keys of the cache in arbitrary order.
	Keys() []K

	// Count returns the number of key-value pairs in the cache.
	Count() int

	setValueWithTimeout(K, valueWithTimeout[V])
}

type CacheBuilder[K comparable, V any] struct {
	et    EvictType
	size  uint
	tmIvl time.Duration
}

func New[K comparable, V any](size uint) CacheBuilder[K, V] {
	return CacheBuilder[K, V]{
		size: size,
		et:   Manual,
	}
}

func (cb CacheBuilder[K, V]) TimeInterval(t time.Duration) CacheBuilder[K, V] {
	cb.tmIvl = t
	return cb
}

func (b *CacheBuilder[K, V]) EvictType(evictType EvictType) {
	b.et = evictType
}

func (b *CacheBuilder[K, V]) Build() Cache[K, V] {
	switch b.et {
	case Manual:
		return newManual(WithTimeInterval[K, V](b.tmIvl))
	default:
		panic("in-memdb: unknown evict-type")
	}
}

type baseCache[K comparable, V any] struct {
	m  map[K]valueWithTimeout[V]
	mu sync.RWMutex
}

type EvictType string

const (
	Manual EvictType = "manual"
)
