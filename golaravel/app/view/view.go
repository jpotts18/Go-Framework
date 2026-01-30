package view

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type View struct {
	basePath  string
	extension string
	templates map[string]*template.Template
	shared    map[string]interface{}
	funcMap   template.FuncMap
	mu        sync.RWMutex
}

func New(basePath string) *View {
	return &View{
		basePath:  basePath,
		extension: ".html",
		templates: make(map[string]*template.Template),
		shared:    make(map[string]interface{}),
		funcMap:   defaultFuncMap(),
	}
}

func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"upper":    strings.ToUpper,
		"lower":    strings.ToLower,
		"title":    strings.Title,
		"trim":     strings.TrimSpace,
		"split":    strings.Split,
		"join":     strings.Join,
		"contains": strings.Contains,
		"replace":  strings.ReplaceAll,
		"default": func(defaultVal, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultVal
			}
			return value
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"css": func(s string) template.CSS {
			return template.CSS(s)
		},
		"js": func(s string) template.JS {
			return template.JS(s)
		},
		"url": func(s string) template.URL {
			return template.URL(s)
		},
		"json": func(v interface{}) template.JS {
			return template.JS(fmt.Sprintf("%v", v))
		},
		"loop": func(n int) []int {
			result := make([]int, n)
			for i := 0; i < n; i++ {
				result[i] = i
			}
			return result
		},
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int { return a % b },
		"eq":  func(a, b interface{}) bool { return a == b },
		"ne":  func(a, b interface{}) bool { return a != b },
		"lt":  func(a, b int) bool { return a < b },
		"le":  func(a, b int) bool { return a <= b },
		"gt":  func(a, b int) bool { return a > b },
		"ge":  func(a, b int) bool { return a >= b },
	}
}

func (v *View) SetExtension(ext string) {
	v.extension = ext
}

func (v *View) AddFunc(name string, fn interface{}) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.funcMap[name] = fn
}

func (v *View) Share(key string, value interface{}) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.shared[key] = value
}

func (v *View) ShareMultiple(data map[string]interface{}) {
	v.mu.Lock()
	defer v.mu.Unlock()
	for key, value := range data {
		v.shared[key] = value
	}
}

func (v *View) getTemplatePath(name string) string {
	name = strings.ReplaceAll(name, ".", string(os.PathSeparator))
	return filepath.Join(v.basePath, name+v.extension)
}

func (v *View) loadTemplate(name string) (*template.Template, error) {
	v.mu.RLock()
	if tmpl, exists := v.templates[name]; exists {
		v.mu.RUnlock()
		return tmpl, nil
	}
	v.mu.RUnlock()

	path := v.getTemplatePath(name)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", name, err)
	}

	tmpl, err := template.New(name).Funcs(v.funcMap).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	v.mu.Lock()
	v.templates[name] = tmpl
	v.mu.Unlock()

	return tmpl, nil
}

func (v *View) Render(name string, data map[string]interface{}) (string, error) {
	tmpl, err := v.loadTemplate(name)
	if err != nil {
		return "", err
	}

	mergedData := make(map[string]interface{})
	v.mu.RLock()
	for k, val := range v.shared {
		mergedData[k] = val
	}
	v.mu.RUnlock()

	for k, val := range data {
		mergedData[k] = val
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, mergedData); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.String(), nil
}

func (v *View) RenderTo(w io.Writer, name string, data map[string]interface{}) error {
	content, err := v.Render(name, data)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(content))
	return err
}

func (v *View) Exists(name string) bool {
	path := v.getTemplatePath(name)
	_, err := os.Stat(path)
	return err == nil
}

func (v *View) ClearCache() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.templates = make(map[string]*template.Template)
}

type Component struct {
	Name    string
	Props   map[string]interface{}
	Slots   map[string]string
	Content string
}

func (v *View) Component(name string, props map[string]interface{}) *Component {
	return &Component{
		Name:  name,
		Props: props,
		Slots: make(map[string]string),
	}
}

func (c *Component) Slot(name, content string) *Component {
	c.Slots[name] = content
	return c
}

func (c *Component) SetContent(content string) *Component {
	c.Content = content
	return c
}
