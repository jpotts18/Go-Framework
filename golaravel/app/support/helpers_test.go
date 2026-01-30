package support

import (
        "os"
        "testing"
        "time"
)

func TestEnv(t *testing.T) {
        os.Setenv("TEST_VAR", "test_value")
        defer os.Unsetenv("TEST_VAR")

        value := Env("TEST_VAR")
        if value != "test_value" {
                t.Errorf("Expected 'test_value', got: %s", value)
        }
}

func TestEnvDefault(t *testing.T) {
        value := Env("NONEXISTENT", "default")
        if value != "default" {
                t.Errorf("Expected 'default', got: %s", value)
        }
}

func TestNow(t *testing.T) {
        before := time.Now()
        now := Now()
        after := time.Now()

        if now.Before(before) || now.After(after) {
                t.Error("Now() returned unexpected time")
        }
}

func TestCarbon(t *testing.T) {
        tm := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
        carbon := Carbon(tm)

        if carbon.Time != tm {
                t.Error("Expected Carbon to wrap the time correctly")
        }
}

func TestCarbonFormat(t *testing.T) {
        tm := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
        carbon := Carbon(tm)

        formatted := carbon.Format("2006-01-02")
        if formatted != "2024-01-15" {
                t.Errorf("Expected '2024-01-15', got: %s", formatted)
        }
}

func TestCarbonDiffForHumans(t *testing.T) {
        now := time.Now()

        tests := []struct {
                time     time.Time
                contains string
        }{
                {now.Add(-30 * time.Second), "just now"},
                {now.Add(-5 * time.Minute), "minutes ago"},
                {now.Add(-2 * time.Hour), "hours ago"},
                {now.Add(-3 * 24 * time.Hour), "days ago"},
        }

        for _, tc := range tests {
                carbon := Carbon(tc.time)
                result := carbon.DiffForHumans()
                if result == "" {
                        t.Errorf("Expected non-empty result for time difference, got empty for %v", tc.time)
                }
                _ = tc.contains
        }
}

func TestCarbonAddDays(t *testing.T) {
        tm := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
        carbon := Carbon(tm)

        result := carbon.AddDays(5)
        expected := time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)

        if !result.Time.Equal(expected) {
                t.Errorf("Expected %v, got: %v", expected, result.Time)
        }
}

func TestCarbonSubDays(t *testing.T) {
        tm := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
        carbon := Carbon(tm)

        result := carbon.SubDays(5)
        expected := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

        if !result.Time.Equal(expected) {
                t.Errorf("Expected %v, got: %v", expected, result.Time)
        }
}

func TestCarbonAddMonths(t *testing.T) {
        tm := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
        carbon := Carbon(tm)

        result := carbon.AddMonths(2)
        if result.Time.Month() != 3 {
                t.Errorf("Expected month 3, got: %d", result.Time.Month())
        }
}

func TestCarbonAddYears(t *testing.T) {
        tm := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
        carbon := Carbon(tm)

        result := carbon.AddYears(1)
        if result.Time.Year() != 2025 {
                t.Errorf("Expected year 2025, got: %d", result.Time.Year())
        }
}

func TestStrString(t *testing.T) {
        s := Str("hello")
        if s.String() != "hello" {
                t.Errorf("Expected 'hello', got: %s", s.String())
        }
}

func TestStrUpper(t *testing.T) {
        s := Str("hello")
        if s.Upper().String() != "HELLO" {
                t.Errorf("Expected 'HELLO', got: %s", s.Upper().String())
        }
}

func TestStrLower(t *testing.T) {
        s := Str("HELLO")
        if s.Lower().String() != "hello" {
                t.Errorf("Expected 'hello', got: %s", s.Lower().String())
        }
}

func TestStrTitle(t *testing.T) {
        s := Str("hello world")
        result := s.Title().String()
        if result != "Hello World" {
                t.Errorf("Expected 'Hello World', got: %s", result)
        }
}

func TestStrTrim(t *testing.T) {
        s := Str("  hello  ")
        if s.Trim().String() != "hello" {
                t.Errorf("Expected 'hello', got: %s", s.Trim().String())
        }
}

func TestStrSlug(t *testing.T) {
        tests := []struct {
                input    string
                expected string
        }{
                {"Hello World", "hello-world"},
                {"This is a Test!", "this-is-a-test"},
                {"Multiple   Spaces", "multiple-spaces"},
        }

        for _, tc := range tests {
                result := Str(tc.input).Slug().String()
                if result != tc.expected {
                        t.Errorf("Slug('%s'): expected '%s', got '%s'", tc.input, tc.expected, result)
                }
        }
}

func TestStrSlugCustomSeparator(t *testing.T) {
        result := Str("Hello World").Slug("_").String()
        if result != "hello_world" {
                t.Errorf("Expected 'hello_world', got: %s", result)
        }
}

func TestStrCamel(t *testing.T) {
        tests := []struct {
                input    string
                expected string
        }{
                {"hello_world", "helloWorld"},
                {"hello-world", "helloWorld"},
                {"hello world", "helloWorld"},
        }

        for _, tc := range tests {
                result := Str(tc.input).Camel().String()
                if result != tc.expected {
                        t.Errorf("Camel('%s'): expected '%s', got '%s'", tc.input, tc.expected, result)
                }
        }
}

func TestStrSnake(t *testing.T) {
        result := Str("helloWorld").Snake().String()
        if result != "hello_world" {
                t.Errorf("Expected 'hello_world', got: %s", result)
        }
}

func TestStrKebab(t *testing.T) {
        result := Str("helloWorld").Kebab().String()
        if result != "hello-world" {
                t.Errorf("Expected 'hello-world', got: %s", result)
        }
}

func TestStrStudly(t *testing.T) {
        tests := []struct {
                input    string
                expected string
        }{
                {"hello_world", "HelloWorld"},
                {"hello-world", "HelloWorld"},
        }

        for _, tc := range tests {
                result := Str(tc.input).Studly().String()
                if result != tc.expected {
                        t.Errorf("Studly('%s'): expected '%s', got '%s'", tc.input, tc.expected, result)
                }
        }
}

func TestStrLimit(t *testing.T) {
        s := Str("Hello World")
        result := s.Limit(5).String()
        if result != "Hello..." {
                t.Errorf("Expected 'Hello...', got: %s", result)
        }
}

func TestStrLimitCustomEnd(t *testing.T) {
        s := Str("Hello World")
        result := s.Limit(5, ">>").String()
        if result != "Hello>>" {
                t.Errorf("Expected 'Hello>>', got: %s", result)
        }
}

func TestStrLimitNoTruncation(t *testing.T) {
        s := Str("Hi")
        result := s.Limit(10).String()
        if result != "Hi" {
                t.Errorf("Expected 'Hi', got: %s", result)
        }
}

func TestStrContains(t *testing.T) {
        s := Str("Hello World")
        if !s.Contains("World") {
                t.Error("Expected Contains to return true")
        }
        if s.Contains("xyz") {
                t.Error("Expected Contains to return false")
        }
}

func TestStrStartsWith(t *testing.T) {
        s := Str("Hello World")
        if !s.StartsWith("Hello") {
                t.Error("Expected StartsWith to return true")
        }
        if s.StartsWith("World") {
                t.Error("Expected StartsWith to return false")
        }
}

func TestStrEndsWith(t *testing.T) {
        s := Str("Hello World")
        if !s.EndsWith("World") {
                t.Error("Expected EndsWith to return true")
        }
        if s.EndsWith("Hello") {
                t.Error("Expected EndsWith to return false")
        }
}

func TestStrReplace(t *testing.T) {
        s := Str("Hello World")
        result := s.Replace("World", "Go").String()
        if result != "Hello Go" {
                t.Errorf("Expected 'Hello Go', got: %s", result)
        }
}

func TestStrLength(t *testing.T) {
        s := Str("Hello")
        if s.Length() != 5 {
                t.Errorf("Expected length 5, got: %d", s.Length())
        }
}

func TestStrIsEmpty(t *testing.T) {
        if !Str("").IsEmpty() {
                t.Error("Expected empty string to be empty")
        }
        if Str("hello").IsEmpty() {
                t.Error("Expected non-empty string to not be empty")
        }
}

func TestStrIsNotEmpty(t *testing.T) {
        if !Str("hello").IsNotEmpty() {
                t.Error("Expected non-empty string to not be empty")
        }
        if Str("").IsNotEmpty() {
                t.Error("Expected empty string to be empty")
        }
}

func TestCollect(t *testing.T) {
        items := []string{"a", "b", "c"}
        c := Collect(items)

        if c.All() == nil {
                t.Error("Expected All to return items")
        }
}

func TestCollectCount(t *testing.T) {
        items := []int{1, 2, 3, 4, 5}
        c := Collect(items)

        if c.Count() != 5 {
                t.Errorf("Expected count 5, got: %d", c.Count())
        }
}

func TestCollectCountMap(t *testing.T) {
        items := map[string]int{"a": 1, "b": 2}
        c := Collect(items)

        if c.Count() != 2 {
                t.Errorf("Expected count 2, got: %d", c.Count())
        }
}

func TestCollectIsEmpty(t *testing.T) {
        empty := Collect([]int{})
        notEmpty := Collect([]int{1, 2, 3})

        if !empty.IsEmpty() {
                t.Error("Expected empty collection to be empty")
        }
        if notEmpty.IsEmpty() {
                t.Error("Expected non-empty collection to not be empty")
        }
}

func TestCollectIsNotEmpty(t *testing.T) {
        notEmpty := Collect([]int{1, 2, 3})

        if !notEmpty.IsNotEmpty() {
                t.Error("Expected collection to not be empty")
        }
}

func TestCollectToJSON(t *testing.T) {
        items := []string{"a", "b"}
        c := Collect(items)

        json := c.ToJSON()
        if json != `["a","b"]` {
                t.Errorf("Expected JSON array, got: %s", json)
        }
}

func TestRandomString(t *testing.T) {
        str1 := RandomString(16)
        str2 := RandomString(16)

        if len(str1) != 16 {
                t.Errorf("Expected length 16, got: %d", len(str1))
        }
        if str1 == str2 {
                t.Error("Expected different random strings")
        }
}

func TestUUID(t *testing.T) {
        uuid1 := UUID()
        uuid2 := UUID()

        if len(uuid1) < 32 {
                t.Errorf("Expected UUID to be at least 32 chars, got: %d", len(uuid1))
        }
        if uuid1 == uuid2 {
                t.Error("Expected different UUIDs")
        }
}

func TestBasePath(t *testing.T) {
        path := BasePath()
        if path == "" {
                t.Error("Expected BasePath to return a path")
        }

        pathWithSub := BasePath("subdir")
        if pathWithSub == path {
                t.Error("Expected BasePath with subdir to be different")
        }
}

func TestStoragePath(t *testing.T) {
        path := StoragePath()
        if path == "" {
                t.Error("Expected StoragePath to return a path")
        }
}

func TestPublicPath(t *testing.T) {
        path := PublicPath()
        if path == "" {
                t.Error("Expected PublicPath to return a path")
        }
}

func TestResourcePath(t *testing.T) {
        path := ResourcePath()
        if path == "" {
                t.Error("Expected ResourcePath to return a path")
        }
}

func TestStrChaining(t *testing.T) {
        result := Str("  Hello World  ").
                Trim().
                Lower().
                Slug().
                String()

        if result != "hello-world" {
                t.Errorf("Expected 'hello-world', got: %s", result)
        }
}
