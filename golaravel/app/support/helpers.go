package support

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func Env(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func DD(values ...interface{}) {
	for _, v := range values {
		fmt.Printf("%#v\n", v)
	}
	os.Exit(1)
}

func Dump(values ...interface{}) {
	for _, v := range values {
		fmt.Printf("%#v\n", v)
	}
}

func Now() time.Time {
	return time.Now()
}

func Carbon(t time.Time) *CarbonTime {
	return &CarbonTime{Time: t}
}

type CarbonTime struct {
	time.Time
}

func (c *CarbonTime) Format(layout string) string {
	return c.Time.Format(layout)
}

func (c *CarbonTime) DiffForHumans() string {
	diff := time.Since(c.Time)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	years := int(diff.Hours() / 24 / 365)
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

func (c *CarbonTime) AddDays(days int) *CarbonTime {
	return &CarbonTime{Time: c.Time.AddDate(0, 0, days)}
}

func (c *CarbonTime) AddMonths(months int) *CarbonTime {
	return &CarbonTime{Time: c.Time.AddDate(0, months, 0)}
}

func (c *CarbonTime) AddYears(years int) *CarbonTime {
	return &CarbonTime{Time: c.Time.AddDate(years, 0, 0)}
}

func (c *CarbonTime) SubDays(days int) *CarbonTime {
	return c.AddDays(-days)
}

func (c *CarbonTime) SubMonths(months int) *CarbonTime {
	return c.AddMonths(-months)
}

func (c *CarbonTime) SubYears(years int) *CarbonTime {
	return c.AddYears(-years)
}

func Str(s string) *Stringable {
	return &Stringable{value: s}
}

type Stringable struct {
	value string
}

func (s *Stringable) String() string {
	return s.value
}

func (s *Stringable) Upper() *Stringable {
	return &Stringable{value: strings.ToUpper(s.value)}
}

func (s *Stringable) Lower() *Stringable {
	return &Stringable{value: strings.ToLower(s.value)}
}

func (s *Stringable) Title() *Stringable {
	return &Stringable{value: strings.Title(s.value)}
}

func (s *Stringable) Trim() *Stringable {
	return &Stringable{value: strings.TrimSpace(s.value)}
}

func (s *Stringable) Slug(separator ...string) *Stringable {
	sep := "-"
	if len(separator) > 0 {
		sep = separator[0]
	}

	result := strings.ToLower(s.value)
	result = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`[\s-]+`).ReplaceAllString(result, sep)
	result = strings.Trim(result, sep)

	return &Stringable{value: result}
}

func (s *Stringable) Camel() *Stringable {
	parts := regexp.MustCompile(`[\s_-]+`).Split(s.value, -1)
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			runes := []rune(parts[i])
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return &Stringable{value: strings.Join(parts, "")}
}

func (s *Stringable) Snake(separator ...string) *Stringable {
	sep := "_"
	if len(separator) > 0 {
		sep = separator[0]
	}

	result := regexp.MustCompile(`([A-Z])`).ReplaceAllString(s.value, sep+"$1")
	result = strings.ToLower(result)
	result = strings.TrimPrefix(result, sep)

	return &Stringable{value: result}
}

func (s *Stringable) Kebab() *Stringable {
	return s.Snake("-")
}

func (s *Stringable) Studly() *Stringable {
	parts := regexp.MustCompile(`[\s_-]+`).Split(s.value, -1)
	for i := range parts {
		if len(parts[i]) > 0 {
			runes := []rune(parts[i])
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return &Stringable{value: strings.Join(parts, "")}
}

func (s *Stringable) Limit(limit int, end ...string) *Stringable {
	suffix := "..."
	if len(end) > 0 {
		suffix = end[0]
	}

	if len(s.value) <= limit {
		return s
	}

	return &Stringable{value: s.value[:limit] + suffix}
}

func (s *Stringable) Contains(needle string) bool {
	return strings.Contains(s.value, needle)
}

func (s *Stringable) StartsWith(prefix string) bool {
	return strings.HasPrefix(s.value, prefix)
}

func (s *Stringable) EndsWith(suffix string) bool {
	return strings.HasSuffix(s.value, suffix)
}

func (s *Stringable) Replace(old, new string) *Stringable {
	return &Stringable{value: strings.ReplaceAll(s.value, old, new)}
}

func (s *Stringable) Length() int {
	return len(s.value)
}

func (s *Stringable) IsEmpty() bool {
	return s.value == ""
}

func (s *Stringable) IsNotEmpty() bool {
	return s.value != ""
}

func Collect(items interface{}) *Collection {
	return &Collection{items: items}
}

type Collection struct {
	items interface{}
}

func (c *Collection) All() interface{} {
	return c.items
}

func (c *Collection) ToJSON() string {
	bytes, _ := json.Marshal(c.items)
	return string(bytes)
}

func (c *Collection) Count() int {
	v := reflect.ValueOf(c.items)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return v.Len()
	}
	if v.Kind() == reflect.Map {
		return v.Len()
	}
	return 0
}

func (c *Collection) IsEmpty() bool {
	return c.Count() == 0
}

func (c *Collection) IsNotEmpty() bool {
	return c.Count() > 0
}

func RandomString(length int) string {
	bytes := make([]byte, length/2+1)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

func UUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func BasePath(path ...string) string {
	base, _ := os.Getwd()
	if len(path) > 0 {
		return filepath.Join(base, path[0])
	}
	return base
}

func StoragePath(path ...string) string {
	base := BasePath("storage")
	if len(path) > 0 {
		return filepath.Join(base, path[0])
	}
	return base
}

func PublicPath(path ...string) string {
	base := BasePath("public")
	if len(path) > 0 {
		return filepath.Join(base, path[0])
	}
	return base
}

func ResourcePath(path ...string) string {
	base := BasePath("resources")
	if len(path) > 0 {
		return filepath.Join(base, path[0])
	}
	return base
}
