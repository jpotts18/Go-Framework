package view

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestView(t *testing.T) (*View, string) {
	tempDir, err := os.MkdirTemp("", "view_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	return New(tempDir), tempDir
}

func cleanupTestView(tempDir string) {
	os.RemoveAll(tempDir)
}

func TestNewView(t *testing.T) {
	v := New("/tmp/views")
	if v == nil {
		t.Fatal("Expected view to be created")
	}
	if v.basePath != "/tmp/views" {
		t.Errorf("Expected basePath to be '/tmp/views', got: %s", v.basePath)
	}
	if v.extension != ".html" {
		t.Errorf("Expected default extension to be '.html', got: %s", v.extension)
	}
}

func TestViewSetExtension(t *testing.T) {
	v := New("/tmp/views")
	v.SetExtension(".tmpl")

	if v.extension != ".tmpl" {
		t.Errorf("Expected extension to be '.tmpl', got: %s", v.extension)
	}
}

func TestViewShare(t *testing.T) {
	v := New("/tmp/views")
	v.Share("appName", "MyApp")

	if v.shared["appName"] != "MyApp" {
		t.Error("Expected shared data to be set")
	}
}

func TestViewShareMultiple(t *testing.T) {
	v := New("/tmp/views")
	v.ShareMultiple(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})

	if v.shared["key1"] != "value1" || v.shared["key2"] != "value2" {
		t.Error("Expected multiple shared data to be set")
	}
}

func TestViewAddFunc(t *testing.T) {
	v := New("/tmp/views")
	v.AddFunc("double", func(n int) int { return n * 2 })

	if v.funcMap["double"] == nil {
		t.Error("Expected custom function to be added")
	}
}

func TestViewRender(t *testing.T) {
	v, tempDir := setupTestView(t)
	defer cleanupTestView(tempDir)

	templateContent := `Hello, {{.Name}}!`
	templatePath := filepath.Join(tempDir, "greeting.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	result, err := v.Render("greeting", map[string]interface{}{
		"Name": "World",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got: %s", result)
	}
}

func TestViewRenderWithSharedData(t *testing.T) {
	v, tempDir := setupTestView(t)
	defer cleanupTestView(tempDir)

	templateContent := `{{.AppName}} - {{.PageTitle}}`
	templatePath := filepath.Join(tempDir, "page.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	v.Share("AppName", "MyApp")

	result, err := v.Render("page", map[string]interface{}{
		"PageTitle": "Home",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != "MyApp - Home" {
		t.Errorf("Expected 'MyApp - Home', got: %s", result)
	}
}

func TestViewRenderNotFound(t *testing.T) {
	v := New("/tmp/nonexistent")

	_, err := v.Render("missing", nil)
	if err == nil {
		t.Error("Expected error for missing template")
	}
}

func TestViewExists(t *testing.T) {
	v, tempDir := setupTestView(t)
	defer cleanupTestView(tempDir)

	templatePath := filepath.Join(tempDir, "exists.html")
	if err := os.WriteFile(templatePath, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	if !v.Exists("exists") {
		t.Error("Expected Exists to return true for existing template")
	}
	if v.Exists("missing") {
		t.Error("Expected Exists to return false for missing template")
	}
}

func TestViewClearCache(t *testing.T) {
	v, tempDir := setupTestView(t)
	defer cleanupTestView(tempDir)

	templatePath := filepath.Join(tempDir, "cached.html")
	if err := os.WriteFile(templatePath, []byte("original"), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	v.Render("cached", nil)

	if len(v.templates) != 1 {
		t.Error("Expected template to be cached")
	}

	v.ClearCache()

	if len(v.templates) != 0 {
		t.Error("Expected cache to be cleared")
	}
}

func TestViewDefaultFunctions(t *testing.T) {
	v, tempDir := setupTestView(t)
	defer cleanupTestView(tempDir)

	tests := []struct {
		template string
		data     map[string]interface{}
		expected string
	}{
		{`{{upper .Text}}`, map[string]interface{}{"Text": "hello"}, "HELLO"},
		{`{{lower .Text}}`, map[string]interface{}{"Text": "HELLO"}, "hello"},
		{`{{add 2 3}}`, nil, "5"},
		{`{{sub 5 3}}`, nil, "2"},
		{`{{mul 3 4}}`, nil, "12"},
		{`{{div 10 2}}`, nil, "5"},
		{`{{mod 7 3}}`, nil, "1"},
		{`{{if eq .A .B}}yes{{end}}`, map[string]interface{}{"A": 1, "B": 1}, "yes"},
		{`{{if ne .A .B}}yes{{end}}`, map[string]interface{}{"A": 1, "B": 2}, "yes"},
		{`{{if lt .A .B}}yes{{end}}`, map[string]interface{}{"A": 1, "B": 2}, "yes"},
		{`{{if gt .A .B}}yes{{end}}`, map[string]interface{}{"A": 2, "B": 1}, "yes"},
	}

	for i, tc := range tests {
		templatePath := filepath.Join(tempDir, "test"+string(rune('0'+i))+".html")
		if err := os.WriteFile(templatePath, []byte(tc.template), 0644); err != nil {
			t.Fatalf("Failed to write template: %v", err)
		}

		result, err := v.Render("test"+string(rune('0'+i)), tc.data)
		if err != nil {
			t.Errorf("Test %d: unexpected error: %v", i, err)
			continue
		}
		if result != tc.expected {
			t.Errorf("Test %d: expected '%s', got '%s'", i, tc.expected, result)
		}
	}
}

func TestViewNestedTemplates(t *testing.T) {
	v, tempDir := setupTestView(t)
	defer cleanupTestView(tempDir)

	subDir := filepath.Join(tempDir, "layouts")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	templateContent := `Layout: {{.Content}}`
	templatePath := filepath.Join(subDir, "base.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	result, err := v.Render("layouts.base", map[string]interface{}{
		"Content": "Main Content",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != "Layout: Main Content" {
		t.Errorf("Expected 'Layout: Main Content', got: %s", result)
	}
}

func TestViewComponent(t *testing.T) {
	v := New("/tmp/views")

	comp := v.Component("button", map[string]interface{}{
		"color": "blue",
		"size":  "large",
	})

	if comp.Name != "button" {
		t.Errorf("Expected component name to be 'button', got: %s", comp.Name)
	}
	if comp.Props["color"] != "blue" {
		t.Error("Expected color prop to be 'blue'")
	}
}

func TestViewComponentSlot(t *testing.T) {
	v := New("/tmp/views")

	comp := v.Component("card", nil).
		Slot("header", "Card Header").
		Slot("body", "Card Body")

	if comp.Slots["header"] != "Card Header" {
		t.Error("Expected header slot to be set")
	}
	if comp.Slots["body"] != "Card Body" {
		t.Error("Expected body slot to be set")
	}
}

func TestViewLoopFunction(t *testing.T) {
	v, tempDir := setupTestView(t)
	defer cleanupTestView(tempDir)

	templateContent := `{{range loop 3}}x{{end}}`
	templatePath := filepath.Join(tempDir, "loop.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	result, err := v.Render("loop", nil)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != "xxx" {
		t.Errorf("Expected 'xxx', got: %s", result)
	}
}

func TestViewDefaultFunction(t *testing.T) {
	v, tempDir := setupTestView(t)
	defer cleanupTestView(tempDir)

	templateContent := `{{default "fallback" .Value}}`
	templatePath := filepath.Join(tempDir, "default.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	result1, _ := v.Render("default", map[string]interface{}{"Value": ""})
	if result1 != "fallback" {
		t.Errorf("Expected 'fallback' for empty value, got: %s", result1)
	}

	v.ClearCache()
	result2, _ := v.Render("default", map[string]interface{}{"Value": "actual"})
	if result2 != "actual" {
		t.Errorf("Expected 'actual' for non-empty value, got: %s", result2)
	}
}
