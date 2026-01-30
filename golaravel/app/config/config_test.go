package config

import (
        "os"
        "testing"
)

func TestConfigNew(t *testing.T) {
        c := New()
        if c == nil {
                t.Fatal("Expected config to be created")
        }
        if c.data == nil {
                t.Error("Expected data map to be initialized")
        }
}

func TestConfigSetAndGet(t *testing.T) {
        c := New()

        c.Set("app.name", "TestApp")
        value := c.Get("app.name")

        if value != "TestApp" {
                t.Errorf("Expected 'TestApp', got: %v", value)
        }
}

func TestConfigNestedSet(t *testing.T) {
        c := New()

        c.Set("database.connection.host", "localhost")
        c.Set("database.connection.port", 5432)

        host := c.Get("database.connection.host")
        port := c.Get("database.connection.port")

        if host != "localhost" {
                t.Errorf("Expected 'localhost', got: %v", host)
        }
        if port != 5432 {
                t.Errorf("Expected 5432, got: %v", port)
        }
}

func TestConfigGetDefault(t *testing.T) {
        c := New()

        value := c.Get("nonexistent", "default")
        if value != "default" {
                t.Errorf("Expected 'default', got: %v", value)
        }
}

func TestConfigGetNil(t *testing.T) {
        c := New()

        value := c.Get("nonexistent")
        if value != nil {
                t.Errorf("Expected nil, got: %v", value)
        }
}

func TestConfigGetString(t *testing.T) {
        c := New()
        c.Set("name", "John")

        value := c.GetString("name")
        if value != "John" {
                t.Errorf("Expected 'John', got: %s", value)
        }
}

func TestConfigGetStringDefault(t *testing.T) {
        c := New()

        value := c.GetString("nonexistent", "default")
        if value != "default" {
                t.Errorf("Expected 'default', got: %s", value)
        }
}

func TestConfigGetInt(t *testing.T) {
        c := New()
        c.Set("port", 8080)

        value := c.GetInt("port")
        if value != 8080 {
                t.Errorf("Expected 8080, got: %d", value)
        }
}

func TestConfigGetIntFromFloat(t *testing.T) {
        c := New()
        c.Set("value", 42.0)

        value := c.GetInt("value")
        if value != 42 {
                t.Errorf("Expected 42, got: %d", value)
        }
}

func TestConfigGetIntFromString(t *testing.T) {
        c := New()
        c.Set("port", "3000")

        value := c.GetInt("port")
        if value != 3000 {
                t.Errorf("Expected 3000, got: %d", value)
        }
}

func TestConfigGetIntDefault(t *testing.T) {
        c := New()

        value := c.GetInt("nonexistent", 1234)
        if value != 1234 {
                t.Errorf("Expected 1234, got: %d", value)
        }
}

func TestConfigGetBool(t *testing.T) {
        c := New()
        c.Set("debug", true)

        value := c.GetBool("debug")
        if !value {
                t.Error("Expected true")
        }
}

func TestConfigGetBoolFromString(t *testing.T) {
        c := New()
        c.Set("enabled", "true")

        value := c.GetBool("enabled")
        if !value {
                t.Error("Expected true from string 'true'")
        }
}

func TestConfigGetBoolFromInt(t *testing.T) {
        c := New()
        c.Set("flag", 1)

        value := c.GetBool("flag")
        if !value {
                t.Error("Expected true from int 1")
        }
}

func TestConfigGetBoolDefault(t *testing.T) {
        c := New()

        value := c.GetBool("nonexistent", true)
        if !value {
                t.Error("Expected default true")
        }
}

func TestConfigGetSlice(t *testing.T) {
        c := New()
        c.Set("items", []interface{}{"a", "b", "c"})

        slice := c.GetSlice("items")
        if len(slice) != 3 {
                t.Errorf("Expected 3 items, got: %d", len(slice))
        }
}

func TestConfigGetSliceNil(t *testing.T) {
        c := New()

        slice := c.GetSlice("nonexistent")
        if slice != nil {
                t.Error("Expected nil for nonexistent slice")
        }
}

func TestConfigGetMap(t *testing.T) {
        c := New()
        c.Set("settings", map[string]interface{}{"key": "value"})

        m := c.GetMap("settings")
        if m == nil {
                t.Fatal("Expected map to be returned")
        }
        if m["key"] != "value" {
                t.Errorf("Expected 'value', got: %v", m["key"])
        }
}

func TestConfigHas(t *testing.T) {
        c := New()
        c.Set("exists", "value")

        if !c.Has("exists") {
                t.Error("Expected Has to return true for existing key")
        }
        if c.Has("nonexistent") {
                t.Error("Expected Has to return false for nonexistent key")
        }
}

func TestConfigAll(t *testing.T) {
        c := New()
        c.Set("key1", "value1")
        c.Set("key2", "value2")

        all := c.All()
        if len(all) != 2 {
                t.Errorf("Expected 2 items, got: %d", len(all))
        }
}

func TestEnv(t *testing.T) {
        os.Setenv("TEST_ENV_VAR", "test_value")
        defer os.Unsetenv("TEST_ENV_VAR")

        value := Env("TEST_ENV_VAR")
        if value != "test_value" {
                t.Errorf("Expected 'test_value', got: %s", value)
        }
}

func TestEnvDefault(t *testing.T) {
        value := Env("NONEXISTENT_VAR", "default")
        if value != "default" {
                t.Errorf("Expected 'default', got: %s", value)
        }
}

func TestEnvInt(t *testing.T) {
        os.Setenv("TEST_PORT", "8080")
        defer os.Unsetenv("TEST_PORT")

        value := EnvInt("TEST_PORT")
        if value != 8080 {
                t.Errorf("Expected 8080, got: %d", value)
        }
}

func TestEnvIntDefault(t *testing.T) {
        value := EnvInt("NONEXISTENT_PORT", 3000)
        if value != 3000 {
                t.Errorf("Expected 3000, got: %d", value)
        }
}

func TestEnvBool(t *testing.T) {
        os.Setenv("TEST_DEBUG", "true")
        defer os.Unsetenv("TEST_DEBUG")

        value := EnvBool("TEST_DEBUG")
        if !value {
                t.Error("Expected true")
        }
}

func TestEnvBoolDefault(t *testing.T) {
        value := EnvBool("NONEXISTENT_DEBUG", true)
        if !value {
                t.Error("Expected default true")
        }
}

func TestConfigInstance(t *testing.T) {
        instance1 := Instance()
        instance2 := Instance()

        if instance1 != instance2 {
                t.Error("Expected Instance() to return same instance (singleton)")
        }
}

func TestConfigLoadEnv(t *testing.T) {
        os.Setenv("MYAPP_TEST_VALUE", "test123")
        defer os.Unsetenv("MYAPP_TEST_VALUE")

        c := New()
        c.LoadEnv("MYAPP")

        if c.Get("test.value") != "test123" {
                t.Errorf("Expected 'test123', got: %v (all: %v)", c.Get("test.value"), c.All())
        }
}
