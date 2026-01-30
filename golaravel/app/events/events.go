package events

import (
	"reflect"
	"sync"
)

type Event interface{}

type Listener func(event Event) error

type Dispatcher struct {
	mu        sync.RWMutex
	listeners map[string][]Listener
	wildcards []Listener
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		listeners: make(map[string][]Listener),
		wildcards: make([]Listener, 0),
	}
}

func (d *Dispatcher) Listen(eventName string, listener Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if eventName == "*" {
		d.wildcards = append(d.wildcards, listener)
		return
	}

	d.listeners[eventName] = append(d.listeners[eventName], listener)
}

func (d *Dispatcher) Dispatch(event Event) error {
	eventName := d.getEventName(event)

	d.mu.RLock()
	listeners := make([]Listener, 0)
	listeners = append(listeners, d.listeners[eventName]...)
	listeners = append(listeners, d.wildcards...)
	d.mu.RUnlock()

	for _, listener := range listeners {
		if err := listener(event); err != nil {
			return err
		}
	}

	return nil
}

func (d *Dispatcher) DispatchAsync(event Event) {
	go func() {
		_ = d.Dispatch(event)
	}()
}

func (d *Dispatcher) getEventName(event Event) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func (d *Dispatcher) HasListeners(eventName string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return len(d.listeners[eventName]) > 0 || len(d.wildcards) > 0
}

func (d *Dispatcher) Forget(eventName string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.listeners, eventName)
}

func (d *Dispatcher) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.listeners = make(map[string][]Listener)
	d.wildcards = make([]Listener, 0)
}

func (d *Dispatcher) Subscribe(subscriber Subscriber) {
	subscriber.Subscribe(d)
}

type Subscriber interface {
	Subscribe(dispatcher *Dispatcher)
}

var defaultDispatcher = NewDispatcher()

func Default() *Dispatcher {
	return defaultDispatcher
}

func Listen(eventName string, listener Listener) {
	defaultDispatcher.Listen(eventName, listener)
}

func Dispatch(event Event) error {
	return defaultDispatcher.Dispatch(event)
}

func DispatchAsync(event Event) {
	defaultDispatcher.DispatchAsync(event)
}

type UserRegistered struct {
	UserID int64
	Email  string
}

type UserLoggedIn struct {
	UserID int64
}

type UserLoggedOut struct {
	UserID int64
}

type PasswordReset struct {
	UserID int64
	Email  string
}

type ModelCreated struct {
	Model interface{}
}

type ModelUpdated struct {
	Model interface{}
}

type ModelDeleted struct {
	Model interface{}
}
