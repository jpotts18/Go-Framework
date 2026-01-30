package session

import (
	"testing"
	"time"
)

func TestGenerateID(t *testing.T) {
	id1, err := GenerateID()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(id1) == 0 {
		t.Error("Expected non-empty ID")
	}

	id2, _ := GenerateID()
	if id1 == id2 {
		t.Error("Expected different IDs")
	}
}

func TestMemoryStorePutAndGet(t *testing.T) {
	store := NewMemoryStore(time.Hour)

	store.Put("session1", "key1", "value1")
	value, exists := store.Get("session1", "key1")

	if !exists {
		t.Error("Expected key to exist")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
}

func TestMemoryStoreGetNotExists(t *testing.T) {
	store := NewMemoryStore(time.Hour)

	_, exists := store.Get("session1", "nonexistent")
	if exists {
		t.Error("Expected key to not exist")
	}
}

func TestMemoryStoreForget(t *testing.T) {
	store := NewMemoryStore(time.Hour)

	store.Put("session1", "key1", "value1")
	store.Forget("session1", "key1")

	_, exists := store.Get("session1", "key1")
	if exists {
		t.Error("Expected key to not exist after forget")
	}
}

func TestMemoryStoreFlush(t *testing.T) {
	store := NewMemoryStore(time.Hour)

	store.Put("session1", "key1", "value1")
	store.Put("session1", "key2", "value2")
	store.Flush("session1")

	all := store.All("session1")
	if len(all) != 0 {
		t.Error("Expected all keys to be flushed")
	}
}

func TestMemoryStoreRegenerate(t *testing.T) {
	store := NewMemoryStore(time.Hour)

	store.Put("old-session", "key1", "value1")
	newID, err := store.Regenerate("old-session")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if newID == "old-session" {
		t.Error("Expected new session ID")
	}

	value, exists := store.Get(newID, "key1")
	if !exists || value != "value1" {
		t.Error("Expected data to be transferred to new session")
	}

	_, exists = store.Get("old-session", "key1")
	if exists {
		t.Error("Expected old session to be destroyed")
	}
}

func TestMemoryStoreDestroy(t *testing.T) {
	store := NewMemoryStore(time.Hour)

	store.Put("session1", "key1", "value1")
	store.Destroy("session1")

	_, exists := store.Get("session1", "key1")
	if exists {
		t.Error("Expected session to be destroyed")
	}
}

func TestMemoryStoreFlash(t *testing.T) {
	store := NewMemoryStore(time.Hour)

	store.Flash("session1", "message", "Hello!")
	value, exists := store.Get("session1", "message")

	if !exists {
		t.Error("Expected flash data to exist")
	}
	if value != "Hello!" {
		t.Errorf("Expected 'Hello!', got %v", value)
	}

	store.AgeFlashData("session1")

	_, exists = store.Get("session1", "message")
	if exists {
		t.Error("Expected flash data to be aged")
	}
}

func TestSession(t *testing.T) {
	store := NewMemoryStore(time.Hour)
	sess, err := New(store)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if sess.ID() == "" {
		t.Error("Expected non-empty session ID")
	}

	sess.Put("name", "John")
	name := sess.Get("name")
	if name != "John" {
		t.Errorf("Expected 'John', got %v", name)
	}

	if !sess.Has("name") {
		t.Error("Expected Has to return true")
	}

	all := sess.All()
	if len(all) != 1 {
		t.Errorf("Expected 1 item in All, got %d", len(all))
	}

	sess.Forget("name")
	if sess.Has("name") {
		t.Error("Expected key to be forgotten")
	}
}

func TestSessionPull(t *testing.T) {
	store := NewMemoryStore(time.Hour)
	sess, _ := New(store)

	sess.Put("key1", "value1")
	value := sess.Pull("key1")

	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
	if sess.Has("key1") {
		t.Error("Expected key to be removed after pull")
	}
}

func TestSessionGetWithDefault(t *testing.T) {
	store := NewMemoryStore(time.Hour)
	sess, _ := New(store)

	value := sess.Get("nonexistent", "default")
	if value != "default" {
		t.Errorf("Expected 'default', got %v", value)
	}
}

func TestSessionRegenerate(t *testing.T) {
	store := NewMemoryStore(time.Hour)
	sess, _ := New(store)
	oldID := sess.ID()

	sess.Put("key1", "value1")
	newID, err := sess.Regenerate()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if newID == oldID {
		t.Error("Expected new session ID")
	}
	if sess.ID() != newID {
		t.Error("Expected session to use new ID")
	}
}

func TestLoadSession(t *testing.T) {
	store := NewMemoryStore(time.Hour)
	store.Put("existing-session", "key1", "value1")

	sess := Load(store, "existing-session")
	value := sess.Get("key1")

	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
}
