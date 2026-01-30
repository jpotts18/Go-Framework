package session

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

type Store interface {
	Get(sessionID, key string) (interface{}, bool)
	Put(sessionID, key string, value interface{}) error
	Forget(sessionID, key string) error
	Flush(sessionID string) error
	All(sessionID string) map[string]interface{}
	Regenerate(oldID string) (string, error)
	Destroy(sessionID string) error
}

type sessionData struct {
	data      map[string]interface{}
	flash     map[string]interface{}
	expiresAt time.Time
}

type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]*sessionData
	lifetime time.Duration
}

func NewMemoryStore(lifetime time.Duration) *MemoryStore {
	store := &MemoryStore{
		sessions: make(map[string]*sessionData),
		lifetime: lifetime,
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
		for id, sess := range s.sessions {
			if now.After(sess.expiresAt) {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}

func (s *MemoryStore) getOrCreate(sessionID string) *sessionData {
	if sess, exists := s.sessions[sessionID]; exists {
		sess.expiresAt = time.Now().Add(s.lifetime)
		return sess
	}

	sess := &sessionData{
		data:      make(map[string]interface{}),
		flash:     make(map[string]interface{}),
		expiresAt: time.Now().Add(s.lifetime),
	}
	s.sessions[sessionID] = sess
	return sess
}

func (s *MemoryStore) Get(sessionID, key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, exists := s.sessions[sessionID]
	if !exists {
		return nil, false
	}

	if value, ok := sess.flash[key]; ok {
		return value, true
	}

	value, ok := sess.data[key]
	return value, ok
}

func (s *MemoryStore) Put(sessionID, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess := s.getOrCreate(sessionID)
	sess.data[key] = value
	return nil
}

func (s *MemoryStore) Flash(sessionID, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess := s.getOrCreate(sessionID)
	sess.flash[key] = value
	return nil
}

func (s *MemoryStore) Forget(sessionID, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sess, exists := s.sessions[sessionID]; exists {
		delete(sess.data, key)
		delete(sess.flash, key)
	}
	return nil
}

func (s *MemoryStore) Flush(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sess, exists := s.sessions[sessionID]; exists {
		sess.data = make(map[string]interface{})
		sess.flash = make(map[string]interface{})
	}
	return nil
}

func (s *MemoryStore) All(sessionID string) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, exists := s.sessions[sessionID]
	if !exists {
		return make(map[string]interface{})
	}

	result := make(map[string]interface{})
	for k, v := range sess.data {
		result[k] = v
	}
	for k, v := range sess.flash {
		result[k] = v
	}
	return result
}

func (s *MemoryStore) Regenerate(oldID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newID, err := GenerateID()
	if err != nil {
		return "", err
	}

	if sess, exists := s.sessions[oldID]; exists {
		s.sessions[newID] = sess
		delete(s.sessions, oldID)
	} else {
		s.sessions[newID] = &sessionData{
			data:      make(map[string]interface{}),
			flash:     make(map[string]interface{}),
			expiresAt: time.Now().Add(s.lifetime),
		}
	}

	return newID, nil
}

func (s *MemoryStore) Destroy(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sessionID)
	return nil
}

func (s *MemoryStore) AgeFlashData(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sess, exists := s.sessions[sessionID]; exists {
		sess.flash = make(map[string]interface{})
	}
}

func GenerateID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

type Session struct {
	id    string
	store Store
}

func New(store Store) (*Session, error) {
	id, err := GenerateID()
	if err != nil {
		return nil, err
	}
	return &Session{id: id, store: store}, nil
}

func Load(store Store, id string) *Session {
	return &Session{id: id, store: store}
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Get(key string, defaultValue ...interface{}) interface{} {
	value, exists := s.store.Get(s.id, key)
	if !exists && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func (s *Session) Put(key string, value interface{}) error {
	return s.store.Put(s.id, key, value)
}

func (s *Session) Forget(key string) error {
	return s.store.Forget(s.id, key)
}

func (s *Session) Flush() error {
	return s.store.Flush(s.id)
}

func (s *Session) All() map[string]interface{} {
	return s.store.All(s.id)
}

func (s *Session) Has(key string) bool {
	_, exists := s.store.Get(s.id, key)
	return exists
}

func (s *Session) Pull(key string, defaultValue ...interface{}) interface{} {
	value := s.Get(key, defaultValue...)
	s.Forget(key)
	return value
}

func (s *Session) Regenerate() (string, error) {
	newID, err := s.store.Regenerate(s.id)
	if err != nil {
		return "", err
	}
	s.id = newID
	return newID, nil
}

func (s *Session) Destroy() error {
	return s.store.Destroy(s.id)
}
