package cache

import (
	"sync"
	"time"
)

type Store interface {
	Get(key string) (interface{}, bool)
	Put(key string, value interface{}, duration time.Duration) error
	Forever(key string, value interface{}) error
	Forget(key string) error
	Flush() error
	Has(key string) bool
	Increment(key string, value int64) (int64, error)
	Decrement(key string, value int64) (int64, error)
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
	forever   bool
}

type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		items: make(map[string]*cacheItem),
	}
	go store.cleanup()
	return store
}

func (s *MemoryStore) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, item := range s.items {
			if !item.forever && now.After(item.expiresAt) {
				delete(s.items, key)
			}
		}
		s.mu.Unlock()
	}
}

func (s *MemoryStore) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[key]
	if !exists {
		return nil, false
	}

	if !item.forever && time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.value, true
}

func (s *MemoryStore) Put(key string, value interface{}, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = &cacheItem{
		value:     value,
		expiresAt: time.Now().Add(duration),
		forever:   false,
	}
	return nil
}

func (s *MemoryStore) Forever(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = &cacheItem{
		value:   value,
		forever: true,
	}
	return nil
}

func (s *MemoryStore) Forget(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
	return nil
}

func (s *MemoryStore) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = make(map[string]*cacheItem)
	return nil
}

func (s *MemoryStore) Has(key string) bool {
	_, exists := s.Get(key)
	return exists
}

func (s *MemoryStore) Increment(key string, value int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[key]
	if !exists {
		s.items[key] = &cacheItem{
			value:   value,
			forever: true,
		}
		return value, nil
	}

	current, ok := item.value.(int64)
	if !ok {
		return 0, nil
	}

	newValue := current + value
	item.value = newValue
	return newValue, nil
}

func (s *MemoryStore) Decrement(key string, value int64) (int64, error) {
	return s.Increment(key, -value)
}

func (s *MemoryStore) Remember(key string, duration time.Duration, callback func() interface{}) interface{} {
	if value, exists := s.Get(key); exists {
		return value
	}

	value := callback()
	s.Put(key, value, duration)
	return value
}

func (s *MemoryStore) RememberForever(key string, callback func() interface{}) interface{} {
	if value, exists := s.Get(key); exists {
		return value
	}

	value := callback()
	s.Forever(key, value)
	return value
}

func (s *MemoryStore) Pull(key string) (interface{}, bool) {
	value, exists := s.Get(key)
	if exists {
		s.Forget(key)
	}
	return value, exists
}

var defaultStore Store = NewMemoryStore()

func Default() Store {
	return defaultStore
}

func SetDefault(store Store) {
	defaultStore = store
}

func Get(key string) (interface{}, bool) {
	return defaultStore.Get(key)
}

func Put(key string, value interface{}, duration time.Duration) error {
	return defaultStore.Put(key, value, duration)
}

func Forever(key string, value interface{}) error {
	return defaultStore.Forever(key, value)
}

func Forget(key string) error {
	return defaultStore.Forget(key)
}

func Flush() error {
	return defaultStore.Flush()
}

func Has(key string) bool {
	return defaultStore.Has(key)
}
