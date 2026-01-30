package container

import (
        "testing"
)

func TestContainerNew(t *testing.T) {
        c := New()
        if c == nil {
                t.Fatal("Expected container to be created")
        }
        if c.bindings == nil {
                t.Error("Expected bindings map to be initialized")
        }
        if c.instances == nil {
                t.Error("Expected instances map to be initialized")
        }
        if c.singletons == nil {
                t.Error("Expected singletons map to be initialized")
        }
}

func TestContainerBind(t *testing.T) {
        c := New()

        c.Bind("test", func(c *Container) interface{} {
                return "test-value"
        })

        if !c.Bound("test") {
                t.Error("Expected 'test' to be bound")
        }

        result, err := c.Make("test")
        if err != nil {
                t.Fatalf("Expected no error, got: %v", err)
        }

        if result != "test-value" {
                t.Errorf("Expected 'test-value', got: %v", result)
        }
}

func TestContainerBindCreatesNewInstanceEachTime(t *testing.T) {
        c := New()
        callCount := 0

        c.Bind("counter", func(c *Container) interface{} {
                callCount++
                return callCount
        })

        result1, _ := c.Make("counter")
        result2, _ := c.Make("counter")

        if result1 == result2 {
                t.Error("Expected different instances for each Make() call with Bind")
        }
        if callCount != 2 {
                t.Errorf("Expected binding to be called twice, got: %d", callCount)
        }
}

func TestContainerSingleton(t *testing.T) {
        c := New()
        callCount := 0

        c.Singleton("singleton", func(c *Container) interface{} {
                callCount++
                return &TestService{Name: "Singleton"}
        })

        result1, _ := c.Make("singleton")
        result2, _ := c.Make("singleton")

        service1 := result1.(*TestService)
        service2 := result2.(*TestService)

        if service1 != service2 {
                t.Error("Expected same instance for Singleton")
        }
        if callCount != 1 {
                t.Errorf("Expected singleton to be called once, got: %d", callCount)
        }
}

func TestContainerInstance(t *testing.T) {
        c := New()

        instance := "my-value"
        c.Instance("my-instance", instance)

        result, err := c.Make("my-instance")
        if err != nil {
                t.Fatalf("Expected no error, got: %v", err)
        }

        if result != instance {
                t.Error("Expected same instance")
        }
}

func TestContainerMakeNotFound(t *testing.T) {
        c := New()

        _, err := c.Make("not-found")
        if err == nil {
                t.Error("Expected error for non-existent binding")
        }
}

func TestContainerMustMake(t *testing.T) {
        c := New()
        c.Bind("test", func(c *Container) interface{} {
                return "value"
        })

        result := c.MustMake("test")
        if result != "value" {
                t.Errorf("Expected 'value', got: %v", result)
        }
}

func TestContainerMustMakePanics(t *testing.T) {
        c := New()

        defer func() {
                if r := recover(); r == nil {
                        t.Error("Expected MustMake to panic for non-existent binding")
                }
        }()

        c.MustMake("not-found")
}

func TestContainerBound(t *testing.T) {
        c := New()

        if c.Bound("test") {
                t.Error("Expected 'test' to not be bound")
        }

        c.Bind("test", func(c *Container) interface{} {
                return nil
        })

        if !c.Bound("test") {
                t.Error("Expected 'test' to be bound")
        }
}

func TestContainerBoundWithInstance(t *testing.T) {
        c := New()

        if c.Bound("instance-test") {
                t.Error("Expected 'instance-test' to not be bound")
        }

        c.Instance("instance-test", "some value")

        if !c.Bound("instance-test") {
                t.Error("Expected 'instance-test' to be bound after Instance()")
        }
}

type TestService struct {
        Name string
}

type TestController struct {
        Service *TestService `inject:"test-service"`
}

func TestContainerResolve(t *testing.T) {
        c := New()

        service := &TestService{Name: "MyService"}
        c.Instance("test-service", service)

        controller := &TestController{}
        err := c.Resolve(controller)
        if err != nil {
                t.Fatalf("Expected no error, got: %v", err)
        }

        if controller.Service != service {
                t.Error("Expected service to be injected")
        }
        if controller.Service.Name != "MyService" {
                t.Errorf("Expected service name to be 'MyService', got: %s", controller.Service.Name)
        }
}

func TestContainerResolveNonPointer(t *testing.T) {
        c := New()

        err := c.Resolve("not a pointer")
        if err == nil {
                t.Error("Expected error when resolving non-pointer")
        }
}

func TestContainerFlush(t *testing.T) {
        c := New()

        c.Bind("test1", func(c *Container) interface{} { return "1" })
        c.Instance("test2", "2")
        c.Singleton("test3", func(c *Container) interface{} { return "3" })

        c.Flush()

        if c.Bound("test1") || c.Bound("test2") || c.Bound("test3") {
                t.Error("Expected all bindings to be flushed")
        }
}

func TestContainerNestedDependencies(t *testing.T) {
        c := New()

        c.Singleton("database", func(c *Container) interface{} {
                return map[string]string{"connection": "postgres"}
        })

        c.Bind("repository", func(c *Container) interface{} {
                db := c.MustMake("database").(map[string]string)
                return map[string]interface{}{
                        "db": db,
                }
        })

        repo, err := c.Make("repository")
        if err != nil {
                t.Fatalf("Expected no error, got: %v", err)
        }

        repoMap := repo.(map[string]interface{})
        dbMap := repoMap["db"].(map[string]string)

        if dbMap["connection"] != "postgres" {
                t.Errorf("Expected 'postgres', got: %s", dbMap["connection"])
        }
}
