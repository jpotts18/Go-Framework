package container

import (
	"fmt"
	"reflect"
	"sync"
)

type BindingFunc func(c *Container) interface{}

type Container struct {
	bindings   map[string]BindingFunc
	instances  map[string]interface{}
	singletons map[string]bool
	mu         sync.RWMutex
}

func New() *Container {
	return &Container{
		bindings:   make(map[string]BindingFunc),
		instances:  make(map[string]interface{}),
		singletons: make(map[string]bool),
	}
}

func (c *Container) Bind(abstract string, concrete BindingFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bindings[abstract] = concrete
	c.singletons[abstract] = false
}

func (c *Container) Singleton(abstract string, concrete BindingFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bindings[abstract] = concrete
	c.singletons[abstract] = true
}

func (c *Container) Instance(abstract string, instance interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.instances[abstract] = instance
}

func (c *Container) Make(abstract string) (interface{}, error) {
	c.mu.RLock()
	if instance, exists := c.instances[abstract]; exists {
		c.mu.RUnlock()
		return instance, nil
	}

	binding, exists := c.bindings[abstract]
	isSingleton := c.singletons[abstract]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no binding found for: %s", abstract)
	}

	instance := binding(c)

	if isSingleton {
		c.mu.Lock()
		c.instances[abstract] = instance
		c.mu.Unlock()
	}

	return instance, nil
}

func (c *Container) MustMake(abstract string) interface{} {
	instance, err := c.Make(abstract)
	if err != nil {
		panic(err)
	}
	return instance
}

func (c *Container) Bound(abstract string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.bindings[abstract]
	if !exists {
		_, exists = c.instances[abstract]
	}
	return exists
}

func (c *Container) Resolve(target interface{}) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}

	elem := v.Elem()
	t := elem.Type()

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("target must point to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("inject")
		if tag == "" {
			continue
		}

		instance, err := c.Make(tag)
		if err != nil {
			return fmt.Errorf("failed to resolve %s: %w", tag, err)
		}

		fieldVal := elem.Field(i)
		if !fieldVal.CanSet() {
			return fmt.Errorf("cannot set field %s", field.Name)
		}

		instanceVal := reflect.ValueOf(instance)
		if !instanceVal.Type().AssignableTo(field.Type) {
			return fmt.Errorf("type mismatch for field %s", field.Name)
		}

		fieldVal.Set(instanceVal)
	}

	return nil
}

func (c *Container) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bindings = make(map[string]BindingFunc)
	c.instances = make(map[string]interface{})
	c.singletons = make(map[string]bool)
}
