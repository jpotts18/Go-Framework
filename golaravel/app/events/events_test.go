package events

import (
	"errors"
	"testing"
)

type TestEvent struct {
	Message string
}

func TestDispatcherListen(t *testing.T) {
	d := NewDispatcher()
	called := false

	d.Listen("TestEvent", func(event Event) error {
		called = true
		return nil
	})

	d.Dispatch(&TestEvent{Message: "Hello"})

	if !called {
		t.Error("Expected listener to be called")
	}
}

func TestDispatcherMultipleListeners(t *testing.T) {
	d := NewDispatcher()
	count := 0

	d.Listen("TestEvent", func(event Event) error {
		count++
		return nil
	})
	d.Listen("TestEvent", func(event Event) error {
		count++
		return nil
	})

	d.Dispatch(&TestEvent{Message: "Hello"})

	if count != 2 {
		t.Errorf("Expected 2 listeners to be called, got %d", count)
	}
}

func TestDispatcherWildcard(t *testing.T) {
	d := NewDispatcher()
	called := false

	d.Listen("*", func(event Event) error {
		called = true
		return nil
	})

	d.Dispatch(&TestEvent{Message: "Hello"})

	if !called {
		t.Error("Expected wildcard listener to be called")
	}
}

func TestDispatcherErrorHandling(t *testing.T) {
	d := NewDispatcher()
	expectedErr := errors.New("test error")

	d.Listen("TestEvent", func(event Event) error {
		return expectedErr
	})

	err := d.Dispatch(&TestEvent{Message: "Hello"})

	if err != expectedErr {
		t.Errorf("Expected error '%v', got '%v'", expectedErr, err)
	}
}

func TestDispatcherHasListeners(t *testing.T) {
	d := NewDispatcher()

	if d.HasListeners("TestEvent") {
		t.Error("Expected no listeners initially")
	}

	d.Listen("TestEvent", func(event Event) error {
		return nil
	})

	if !d.HasListeners("TestEvent") {
		t.Error("Expected listeners after adding one")
	}
}

func TestDispatcherForget(t *testing.T) {
	d := NewDispatcher()
	called := false

	d.Listen("TestEvent", func(event Event) error {
		called = true
		return nil
	})

	d.Forget("TestEvent")
	d.Dispatch(&TestEvent{Message: "Hello"})

	if called {
		t.Error("Expected listener to not be called after forget")
	}
}

func TestDispatcherFlush(t *testing.T) {
	d := NewDispatcher()

	d.Listen("TestEvent1", func(event Event) error { return nil })
	d.Listen("TestEvent2", func(event Event) error { return nil })
	d.Listen("*", func(event Event) error { return nil })

	d.Flush()

	if d.HasListeners("TestEvent1") || d.HasListeners("TestEvent2") {
		t.Error("Expected all listeners to be flushed")
	}
}

type TestSubscriber struct {
	Called bool
}

func (s *TestSubscriber) Subscribe(dispatcher *Dispatcher) {
	dispatcher.Listen("TestEvent", func(event Event) error {
		s.Called = true
		return nil
	})
}

func TestDispatcherSubscribe(t *testing.T) {
	d := NewDispatcher()
	subscriber := &TestSubscriber{}

	d.Subscribe(subscriber)
	d.Dispatch(&TestEvent{Message: "Hello"})

	if !subscriber.Called {
		t.Error("Expected subscriber to be called")
	}
}

func TestDefaultDispatcher(t *testing.T) {
	Default().Flush()

	called := false
	Listen("TestEvent", func(event Event) error {
		called = true
		return nil
	})

	Dispatch(&TestEvent{Message: "Hello"})

	if !called {
		t.Error("Expected default listener to be called")
	}
}

func TestBuiltinEvents(t *testing.T) {
	_ = UserRegistered{UserID: 1, Email: "test@example.com"}
	_ = UserLoggedIn{UserID: 1}
	_ = UserLoggedOut{UserID: 1}
	_ = PasswordReset{UserID: 1, Email: "test@example.com"}
	_ = ModelCreated{Model: nil}
	_ = ModelUpdated{Model: nil}
	_ = ModelDeleted{Model: nil}
}
