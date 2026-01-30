package config

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

var (
	instance *Config
	once     sync.Once
)

func Instance() *Config {
	once.Do(func() {
		instance = &Config{
			data: make(map[string]interface{}),
		}
	})
	return instance
}

func New() *Config {
	return &Config{
		data: make(map[string]interface{}),
	}
}

func (c *Config) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	parts := strings.Split(key, ".")
	c.setNested(c.data, parts, value)
}

func (c *Config) setNested(data map[string]interface{}, parts []string, value interface{}) {
	if len(parts) == 1 {
		data[parts[0]] = value
		return
	}

	key := parts[0]
	if _, exists := data[key]; !exists {
		data[key] = make(map[string]interface{})
	}

	if nested, ok := data[key].(map[string]interface{}); ok {
		c.setNested(nested, parts[1:], value)
	}
}

func (c *Config) Get(key string, defaultValue ...interface{}) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	parts := strings.Split(key, ".")
	value := c.getNested(c.data, parts)

	if value == nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return value
}

func (c *Config) getNested(data map[string]interface{}, parts []string) interface{} {
	if len(parts) == 0 {
		return nil
	}

	value, exists := data[parts[0]]
	if !exists {
		return nil
	}

	if len(parts) == 1 {
		return value
	}

	if nested, ok := value.(map[string]interface{}); ok {
		return c.getNested(nested, parts[1:])
	}

	return nil
}

func (c *Config) GetString(key string, defaultValue ...string) string {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}

	if str, ok := value.(string); ok {
		return str
	}

	return ""
}

func (c *Config) GetInt(key string, defaultValue ...int) int {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	return 0
}

func (c *Config) GetBool(key string, defaultValue ...bool) bool {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true" || v == "1"
	case int:
		return v != 0
	}

	return false
}

func (c *Config) GetSlice(key string) []interface{} {
	value := c.Get(key)
	if value == nil {
		return nil
	}

	if slice, ok := value.([]interface{}); ok {
		return slice
	}

	return nil
}

func (c *Config) GetMap(key string) map[string]interface{} {
	value := c.Get(key)
	if value == nil {
		return nil
	}

	if m, ok := value.(map[string]interface{}); ok {
		return m
	}

	return nil
}

func (c *Config) Has(key string) bool {
	return c.Get(key) != nil
}

func (c *Config) All() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]interface{})
	for k, v := range c.data {
		result[k] = v
	}
	return result
}

func (c *Config) LoadJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range jsonData {
		c.data[k] = v
	}

	return nil
}

func (c *Config) LoadEnv(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		if prefix == "" || strings.HasPrefix(key, prefix) {
			configKey := strings.ToLower(strings.ReplaceAll(key, "_", "."))
			if prefix != "" {
				configKey = strings.TrimPrefix(configKey, strings.ToLower(prefix)+".")
			}
			c.data[configKey] = value
		}
	}
}

func Env(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func EnvInt(key string, defaultValue ...int) int {
	value := os.Getenv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	if i, err := strconv.Atoi(value); err == nil {
		return i
	}

	return 0
}

func EnvBool(key string, defaultValue ...bool) bool {
	value := os.Getenv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	return strings.ToLower(value) == "true" || value == "1"
}
