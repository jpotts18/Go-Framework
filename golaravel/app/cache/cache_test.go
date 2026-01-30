package cache

import (
	"testing"
	"time"
)

func TestMemoryStorePutAndGet(t *testing.T) {
	store := NewMemoryStore()

	store.Put("key1", "value1", time.Hour)
	value, exists := store.Get("key1")

	if !exists {
		t.Error("Expected key to exist")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
}

func TestMemoryStoreGetNotExists(t *testing.T) {
	store := NewMemoryStore()

	_, exists := store.Get("nonexistent")
	if exists {
		t.Error("Expected key to not exist")
	}
}

func TestMemoryStoreForever(t *testing.T) {
	store := NewMemoryStore()

	store.Forever("key1", "value1")
	value, exists := store.Get("key1")

	if !exists {
		t.Error("Expected key to exist")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
}

func TestMemoryStoreForget(t *testing.T) {
	store := NewMemoryStore()

	store.Put("key1", "value1", time.Hour)
	store.Forget("key1")

	_, exists := store.Get("key1")
	if exists {
		t.Error("Expected key to not exist after forget")
	}
}

func TestMemoryStoreFlush(t *testing.T) {
	store := NewMemoryStore()

	store.Put("key1", "value1", time.Hour)
	store.Put("key2", "value2", time.Hour)
	store.Flush()

	if store.Has("key1") || store.Has("key2") {
		t.Error("Expected all keys to be flushed")
	}
}

func TestMemoryStoreHas(t *testing.T) {
	store := NewMemoryStore()

	if store.Has("key1") {
		t.Error("Expected Has to return false for nonexistent key")
	}

	store.Put("key1", "value1", time.Hour)

	if !store.Has("key1") {
		t.Error("Expected Has to return true for existing key")
	}
}

func TestMemoryStoreIncrement(t *testing.T) {
	store := NewMemoryStore()

	val, err := store.Increment("counter", 1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}

	val, _ = store.Increment("counter", 5)
	if val != 6 {
		t.Errorf("Expected 6, got %d", val)
	}
}

func TestMemoryStoreDecrement(t *testing.T) {
	store := NewMemoryStore()

	store.Increment("counter", 10)
	val, _ := store.Decrement("counter", 3)

	if val != 7 {
		t.Errorf("Expected 7, got %d", val)
	}
}

func TestMemoryStoreExpiration(t *testing.T) {
	store := NewMemoryStore()

	store.Put("key1", "value1", 10*time.Millisecond)

	value, exists := store.Get("key1")
	if !exists {
		t.Error("Expected key to exist initially")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}

	time.Sleep(20 * time.Millisecond)

	_, exists = store.Get("key1")
	if exists {
		t.Error("Expected key to be expired")
	}
}

func TestMemoryStoreRemember(t *testing.T) {
	store := NewMemoryStore()
	callCount := 0

	callback := func() interface{} {
		callCount++
		return "computed-value"
	}

	val1 := store.Remember("key1", time.Hour, callback)
	val2 := store.Remember("key1", time.Hour, callback)

	if val1 != "computed-value" || val2 != "computed-value" {
		t.Error("Expected 'computed-value'")
	}
	if callCount != 1 {
		t.Errorf("Expected callback to be called once, got %d", callCount)
	}
}

func TestMemoryStorePull(t *testing.T) {
	store := NewMemoryStore()

	store.Put("key1", "value1", time.Hour)

	value, exists := store.Pull("key1")
	if !exists {
		t.Error("Expected key to exist")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}

	_, exists = store.Get("key1")
	if exists {
		t.Error("Expected key to be removed after pull")
	}
}

func TestDefaultFunctions(t *testing.T) {
	Flush()

	Put("testkey", "testvalue", time.Hour)

	if !Has("testkey") {
		t.Error("Expected key to exist")
	}

	val, exists := Get("testkey")
	if !exists || val != "testvalue" {
		t.Error("Expected to get 'testvalue'")
	}

	Forget("testkey")
	if Has("testkey") {
		t.Error("Expected key to be forgotten")
	}
}
