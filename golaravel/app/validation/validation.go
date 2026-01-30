package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Validator struct {
	data   map[string]interface{}
	rules  map[string][]string
	errors map[string][]string
	passed bool
}

func Make(data map[string]interface{}, rules map[string][]string) *Validator {
	v := &Validator{
		data:   data,
		rules:  rules,
		errors: make(map[string][]string),
		passed: true,
	}
	v.validate()
	return v
}

func (v *Validator) validate() {
	for field, fieldRules := range v.rules {
		for _, rule := range fieldRules {
			v.applyRule(field, rule)
		}
	}
}

func (v *Validator) applyRule(field, rule string) {
	value, exists := v.data[field]

	parts := strings.SplitN(rule, ":", 2)
	ruleName := parts[0]
	var ruleParam string
	if len(parts) > 1 {
		ruleParam = parts[1]
	}

	if ruleName == "nullable" && (!exists || value == nil || value == "") {
		return
	}

	switch ruleName {
	case "required":
		if !exists || value == nil || value == "" {
			v.addError(field, "The %s field is required")
		}
	case "email":
		if str, ok := v.getString(value); ok && str != "" {
			if _, err := mail.ParseAddress(str); err != nil {
				v.addError(field, "The %s field must be a valid email address")
			}
		}
	case "min":
		v.validateMin(field, value, ruleParam)
	case "max":
		v.validateMax(field, value, ruleParam)
	case "between":
		v.validateBetween(field, value, ruleParam)
	case "string":
		if _, ok := v.getString(value); !ok && exists && value != nil {
			v.addError(field, "The %s field must be a string")
		}
	case "numeric":
		if !v.isNumeric(value) && exists && value != nil {
			v.addError(field, "The %s field must be numeric")
		}
	case "integer":
		if !v.isInteger(value) && exists && value != nil {
			v.addError(field, "The %s field must be an integer")
		}
	case "boolean":
		if !v.isBoolean(value) && exists && value != nil {
			v.addError(field, "The %s field must be true or false")
		}
	case "alpha":
		if str, ok := v.getString(value); ok && !v.isAlpha(str) {
			v.addError(field, "The %s field must only contain letters")
		}
	case "alpha_num":
		if str, ok := v.getString(value); ok && !v.isAlphaNumeric(str) {
			v.addError(field, "The %s field must only contain letters and numbers")
		}
	case "alpha_dash":
		if str, ok := v.getString(value); ok && !v.isAlphaDash(str) {
			v.addError(field, "The %s field must only contain letters, numbers, dashes and underscores")
		}
	case "url":
		if str, ok := v.getString(value); ok && !v.isURL(str) {
			v.addError(field, "The %s field must be a valid URL")
		}
	case "regex":
		if str, ok := v.getString(value); ok {
			if matched, _ := regexp.MatchString(ruleParam, str); !matched {
				v.addError(field, "The %s field format is invalid")
			}
		}
	case "confirmed":
		confirmField := field + "_confirmation"
		confirmValue, confirmExists := v.data[confirmField]
		if !confirmExists || value != confirmValue {
			v.addError(field, "The %s field confirmation does not match")
		}
	case "in":
		allowed := strings.Split(ruleParam, ",")
		found := false
		strVal, _ := v.getString(value)
		for _, a := range allowed {
			if strVal == a {
				found = true
				break
			}
		}
		if !found {
			v.addError(field, "The selected %s is invalid")
		}
	case "not_in":
		disallowed := strings.Split(ruleParam, ",")
		strVal, _ := v.getString(value)
		for _, d := range disallowed {
			if strVal == d {
				v.addError(field, "The selected %s is invalid")
				break
			}
		}
	case "size":
		v.validateSize(field, value, ruleParam)
	case "date":
		if str, ok := v.getString(value); ok && !v.isDate(str) {
			v.addError(field, "The %s field must be a valid date")
		}
	case "same":
		otherValue, _ := v.data[ruleParam]
		if value != otherValue {
			v.addError(field, "The %s field must match "+ruleParam)
		}
	case "different":
		otherValue, _ := v.data[ruleParam]
		if value == otherValue {
			v.addError(field, "The %s field must be different from "+ruleParam)
		}
	case "nullable":
	}
}

func (v *Validator) validateMin(field string, value interface{}, param string) {
	minVal, _ := strconv.ParseFloat(param, 64)

	switch val := value.(type) {
	case string:
		if float64(len(val)) < minVal {
			v.addError(field, fmt.Sprintf("The %%s field must be at least %s characters", param))
		}
	case int, int64, float64:
		numVal := v.toFloat(val)
		if numVal < minVal {
			v.addError(field, fmt.Sprintf("The %%s field must be at least %s", param))
		}
	case []interface{}:
		if float64(len(val)) < minVal {
			v.addError(field, fmt.Sprintf("The %%s field must have at least %s items", param))
		}
	}
}

func (v *Validator) validateMax(field string, value interface{}, param string) {
	maxVal, _ := strconv.ParseFloat(param, 64)

	switch val := value.(type) {
	case string:
		if float64(len(val)) > maxVal {
			v.addError(field, fmt.Sprintf("The %%s field must not exceed %s characters", param))
		}
	case int, int64, float64:
		numVal := v.toFloat(val)
		if numVal > maxVal {
			v.addError(field, fmt.Sprintf("The %%s field must not exceed %s", param))
		}
	case []interface{}:
		if float64(len(val)) > maxVal {
			v.addError(field, fmt.Sprintf("The %%s field must not have more than %s items", param))
		}
	}
}

func (v *Validator) validateBetween(field string, value interface{}, param string) {
	parts := strings.Split(param, ",")
	if len(parts) != 2 {
		return
	}
	minVal, _ := strconv.ParseFloat(parts[0], 64)
	maxVal, _ := strconv.ParseFloat(parts[1], 64)

	switch val := value.(type) {
	case string:
		length := float64(len(val))
		if length < minVal || length > maxVal {
			v.addError(field, fmt.Sprintf("The %%s field must be between %s and %s characters", parts[0], parts[1]))
		}
	case int, int64, float64:
		numVal := v.toFloat(val)
		if numVal < minVal || numVal > maxVal {
			v.addError(field, fmt.Sprintf("The %%s field must be between %s and %s", parts[0], parts[1]))
		}
	}
}

func (v *Validator) validateSize(field string, value interface{}, param string) {
	sizeVal, _ := strconv.ParseFloat(param, 64)

	switch val := value.(type) {
	case string:
		if float64(len(val)) != sizeVal {
			v.addError(field, fmt.Sprintf("The %%s field must be exactly %s characters", param))
		}
	case int, int64, float64:
		numVal := v.toFloat(val)
		if numVal != sizeVal {
			v.addError(field, fmt.Sprintf("The %%s field must be exactly %s", param))
		}
	case []interface{}:
		if float64(len(val)) != sizeVal {
			v.addError(field, fmt.Sprintf("The %%s field must contain exactly %s items", param))
		}
	}
}

func (v *Validator) addError(field, message string) {
	v.passed = false
	formattedMessage := fmt.Sprintf(message, field)
	if v.errors[field] == nil {
		v.errors[field] = make([]string, 0)
	}
	v.errors[field] = append(v.errors[field], formattedMessage)
}

func (v *Validator) getString(value interface{}) (string, bool) {
	if str, ok := value.(string); ok {
		return str, true
	}
	return "", false
}

func (v *Validator) isNumeric(value interface{}) bool {
	switch val := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case string:
		_, err := strconv.ParseFloat(val, 64)
		return err == nil
	}
	return false
}

func (v *Validator) isInteger(value interface{}) bool {
	switch val := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case string:
		_, err := strconv.ParseInt(val, 10, 64)
		return err == nil
	case float64:
		return val == float64(int64(val))
	}
	return false
}

func (v *Validator) isBoolean(value interface{}) bool {
	switch val := value.(type) {
	case bool:
		return true
	case string:
		lower := strings.ToLower(val)
		return lower == "true" || lower == "false" || lower == "1" || lower == "0"
	case int:
		return val == 0 || val == 1
	}
	return false
}

func (v *Validator) isAlpha(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func (v *Validator) isAlphaNumeric(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func (v *Validator) isAlphaDash(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return false
		}
	}
	return true
}

func (v *Validator) isURL(str string) bool {
	pattern := `^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`
	matched, _ := regexp.MatchString(pattern, str)
	return matched
}

func (v *Validator) isDate(str string) bool {
	patterns := []string{
		`^\d{4}-\d{2}-\d{2}$`,
		`^\d{2}/\d{2}/\d{4}$`,
		`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`,
	}
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, str); matched {
			return true
		}
	}
	return false
}

func (v *Validator) toFloat(value interface{}) float64 {
	switch val := value.(type) {
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case float64:
		return val
	}
	return 0
}

func (v *Validator) Fails() bool {
	return !v.passed
}

func (v *Validator) Passes() bool {
	return v.passed
}

func (v *Validator) Errors() map[string][]string {
	return v.errors
}

func (v *Validator) First(field string) string {
	if errs, exists := v.errors[field]; exists && len(errs) > 0 {
		return errs[0]
	}
	return ""
}

func (v *Validator) All() []string {
	all := make([]string, 0)
	for _, errs := range v.errors {
		all = append(all, errs...)
	}
	return all
}
