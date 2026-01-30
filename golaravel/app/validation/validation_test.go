package validation

import (
	"testing"
)

func TestValidationRequired(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
	}
	rules := map[string][]string{
		"name":  {"required"},
		"email": {"required"},
	}

	v := Make(data, rules)

	if v.Passes() {
		t.Error("Expected validation to fail")
	}
	if v.First("email") == "" {
		t.Error("Expected email error")
	}
	if v.First("name") != "" {
		t.Error("Expected no error for name")
	}
}

func TestValidationRequiredEmpty(t *testing.T) {
	data := map[string]interface{}{
		"name": "",
	}
	rules := map[string][]string{
		"name": {"required"},
	}

	v := Make(data, rules)

	if v.Passes() {
		t.Error("Expected validation to fail for empty string")
	}
}

func TestValidationEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user@domain.org", true},
		{"invalid-email", false},
		{"@nodomain.com", false},
		{"", true},
	}

	for _, tc := range tests {
		data := map[string]interface{}{"email": tc.email}
		rules := map[string][]string{"email": {"email"}}

		v := Make(data, rules)
		if v.Passes() != tc.valid {
			t.Errorf("Email '%s': expected valid=%v, got valid=%v", tc.email, tc.valid, v.Passes())
		}
	}
}

func TestValidationMin(t *testing.T) {
	data := map[string]interface{}{
		"short": "ab",
		"long":  "abcdefgh",
	}
	rules := map[string][]string{
		"short": {"min:5"},
		"long":  {"min:5"},
	}

	v := Make(data, rules)

	if v.First("short") == "" {
		t.Error("Expected error for short string")
	}
	if v.First("long") != "" {
		t.Error("Expected no error for long string")
	}
}

func TestValidationMinNumeric(t *testing.T) {
	data := map[string]interface{}{
		"small": 3,
		"big":   10,
	}
	rules := map[string][]string{
		"small": {"min:5"},
		"big":   {"min:5"},
	}

	v := Make(data, rules)

	if v.First("small") == "" {
		t.Error("Expected error for small number")
	}
	if v.First("big") != "" {
		t.Error("Expected no error for big number")
	}
}

func TestValidationMax(t *testing.T) {
	data := map[string]interface{}{
		"short": "ab",
		"long":  "abcdefghij",
	}
	rules := map[string][]string{
		"short": {"max:5"},
		"long":  {"max:5"},
	}

	v := Make(data, rules)

	if v.First("short") != "" {
		t.Error("Expected no error for short string")
	}
	if v.First("long") == "" {
		t.Error("Expected error for long string")
	}
}

func TestValidationBetween(t *testing.T) {
	data := map[string]interface{}{
		"valid":   "abcd",
		"tooShort": "ab",
		"tooLong": "abcdefghij",
	}
	rules := map[string][]string{
		"valid":    {"between:3,8"},
		"tooShort": {"between:3,8"},
		"tooLong":  {"between:3,8"},
	}

	v := Make(data, rules)

	if v.First("valid") != "" {
		t.Error("Expected no error for valid length")
	}
	if v.First("tooShort") == "" {
		t.Error("Expected error for too short")
	}
	if v.First("tooLong") == "" {
		t.Error("Expected error for too long")
	}
}

func TestValidationNumeric(t *testing.T) {
	data := map[string]interface{}{
		"number":    123,
		"float":     12.5,
		"strNumber": "456",
		"notNumber": "abc",
	}
	rules := map[string][]string{
		"number":    {"numeric"},
		"float":     {"numeric"},
		"strNumber": {"numeric"},
		"notNumber": {"numeric"},
	}

	v := Make(data, rules)

	if v.First("number") != "" {
		t.Error("Expected no error for integer")
	}
	if v.First("float") != "" {
		t.Error("Expected no error for float")
	}
	if v.First("strNumber") != "" {
		t.Error("Expected no error for numeric string")
	}
	if v.First("notNumber") == "" {
		t.Error("Expected error for non-numeric string")
	}
}

func TestValidationInteger(t *testing.T) {
	data := map[string]interface{}{
		"int":    42,
		"float":  12.5,
		"strInt": "123",
	}
	rules := map[string][]string{
		"int":    {"integer"},
		"float":  {"integer"},
		"strInt": {"integer"},
	}

	v := Make(data, rules)

	if v.First("int") != "" {
		t.Error("Expected no error for integer")
	}
	if v.First("float") == "" {
		t.Error("Expected error for float")
	}
	if v.First("strInt") != "" {
		t.Error("Expected no error for integer string")
	}
}

func TestValidationBoolean(t *testing.T) {
	data := map[string]interface{}{
		"bool":    true,
		"strTrue": "true",
		"strOne":  "1",
		"invalid": "maybe",
	}
	rules := map[string][]string{
		"bool":    {"boolean"},
		"strTrue": {"boolean"},
		"strOne":  {"boolean"},
		"invalid": {"boolean"},
	}

	v := Make(data, rules)

	if v.First("bool") != "" {
		t.Error("Expected no error for bool")
	}
	if v.First("strTrue") != "" {
		t.Error("Expected no error for 'true' string")
	}
	if v.First("strOne") != "" {
		t.Error("Expected no error for '1' string")
	}
	if v.First("invalid") == "" {
		t.Error("Expected error for 'maybe'")
	}
}

func TestValidationIn(t *testing.T) {
	data := map[string]interface{}{
		"valid":   "apple",
		"invalid": "grape",
	}
	rules := map[string][]string{
		"valid":   {"in:apple,banana,orange"},
		"invalid": {"in:apple,banana,orange"},
	}

	v := Make(data, rules)

	if v.First("valid") != "" {
		t.Error("Expected no error for valid option")
	}
	if v.First("invalid") == "" {
		t.Error("Expected error for invalid option")
	}
}

func TestValidationNotIn(t *testing.T) {
	data := map[string]interface{}{
		"valid":   "grape",
		"invalid": "apple",
	}
	rules := map[string][]string{
		"valid":   {"not_in:apple,banana,orange"},
		"invalid": {"not_in:apple,banana,orange"},
	}

	v := Make(data, rules)

	if v.First("valid") != "" {
		t.Error("Expected no error for value not in list")
	}
	if v.First("invalid") == "" {
		t.Error("Expected error for value in list")
	}
}

func TestValidationAlpha(t *testing.T) {
	data := map[string]interface{}{
		"letters": "abcXYZ",
		"numbers": "abc123",
		"spaces":  "abc def",
	}
	rules := map[string][]string{
		"letters": {"alpha"},
		"numbers": {"alpha"},
		"spaces":  {"alpha"},
	}

	v := Make(data, rules)

	if v.First("letters") != "" {
		t.Error("Expected no error for letters only")
	}
	if v.First("numbers") == "" {
		t.Error("Expected error for letters with numbers")
	}
	if v.First("spaces") == "" {
		t.Error("Expected error for letters with spaces")
	}
}

func TestValidationAlphaNum(t *testing.T) {
	data := map[string]interface{}{
		"valid":   "abc123",
		"invalid": "abc-123",
	}
	rules := map[string][]string{
		"valid":   {"alpha_num"},
		"invalid": {"alpha_num"},
	}

	v := Make(data, rules)

	if v.First("valid") != "" {
		t.Error("Expected no error for alphanumeric")
	}
	if v.First("invalid") == "" {
		t.Error("Expected error for non-alphanumeric")
	}
}

func TestValidationAlphaDash(t *testing.T) {
	data := map[string]interface{}{
		"valid":   "abc-123_xyz",
		"invalid": "abc 123",
	}
	rules := map[string][]string{
		"valid":   {"alpha_dash"},
		"invalid": {"alpha_dash"},
	}

	v := Make(data, rules)

	if v.First("valid") != "" {
		t.Error("Expected no error for alpha_dash")
	}
	if v.First("invalid") == "" {
		t.Error("Expected error for spaces")
	}
}

func TestValidationURL(t *testing.T) {
	data := map[string]interface{}{
		"valid":   "https://example.com/path",
		"invalid": "not-a-url",
	}
	rules := map[string][]string{
		"valid":   {"url"},
		"invalid": {"url"},
	}

	v := Make(data, rules)

	if v.First("valid") != "" {
		t.Error("Expected no error for valid URL")
	}
	if v.First("invalid") == "" {
		t.Error("Expected error for invalid URL")
	}
}

func TestValidationRegex(t *testing.T) {
	data := map[string]interface{}{
		"valid":   "ABC-123",
		"invalid": "abc123",
	}
	rules := map[string][]string{
		"valid":   {"regex:^[A-Z]+-[0-9]+$"},
		"invalid": {"regex:^[A-Z]+-[0-9]+$"},
	}

	v := Make(data, rules)

	if v.First("valid") != "" {
		t.Error("Expected no error for matching regex")
	}
	if v.First("invalid") == "" {
		t.Error("Expected error for non-matching regex")
	}
}

func TestValidationConfirmed(t *testing.T) {
	data := map[string]interface{}{
		"password":              "secret123",
		"password_confirmation": "secret123",
	}
	rules := map[string][]string{
		"password": {"confirmed"},
	}

	v := Make(data, rules)

	if v.Fails() {
		t.Error("Expected validation to pass when confirmation matches")
	}
}

func TestValidationConfirmedFails(t *testing.T) {
	data := map[string]interface{}{
		"password":              "secret123",
		"password_confirmation": "different",
	}
	rules := map[string][]string{
		"password": {"confirmed"},
	}

	v := Make(data, rules)

	if v.Passes() {
		t.Error("Expected validation to fail when confirmation doesn't match")
	}
}

func TestValidationSame(t *testing.T) {
	data := map[string]interface{}{
		"field1": "value",
		"field2": "value",
	}
	rules := map[string][]string{
		"field1": {"same:field2"},
	}

	v := Make(data, rules)

	if v.Fails() {
		t.Error("Expected validation to pass when fields match")
	}
}

func TestValidationDifferent(t *testing.T) {
	data := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}
	rules := map[string][]string{
		"field1": {"different:field2"},
	}

	v := Make(data, rules)

	if v.Fails() {
		t.Error("Expected validation to pass when fields are different")
	}
}

func TestValidationNullable(t *testing.T) {
	data := map[string]interface{}{
		"name": "",
	}
	rules := map[string][]string{
		"name": {"nullable", "min:3"},
	}

	v := Make(data, rules)

	if v.Fails() {
		t.Error("Expected validation to pass when nullable and empty")
	}
}

func TestValidationSize(t *testing.T) {
	data := map[string]interface{}{
		"code": "ABC",
	}
	rules := map[string][]string{
		"code": {"size:3"},
	}

	v := Make(data, rules)

	if v.Fails() {
		t.Error("Expected validation to pass for exact size")
	}
}

func TestValidationSizeFails(t *testing.T) {
	data := map[string]interface{}{
		"code": "ABCD",
	}
	rules := map[string][]string{
		"code": {"size:3"},
	}

	v := Make(data, rules)

	if v.Passes() {
		t.Error("Expected validation to fail for wrong size")
	}
}

func TestValidationDate(t *testing.T) {
	data := map[string]interface{}{
		"date1": "2024-01-15",
		"date2": "01/15/2024",
		"date3": "not-a-date",
	}
	rules := map[string][]string{
		"date1": {"date"},
		"date2": {"date"},
		"date3": {"date"},
	}

	v := Make(data, rules)

	if v.First("date1") != "" {
		t.Error("Expected no error for YYYY-MM-DD format")
	}
	if v.First("date2") != "" {
		t.Error("Expected no error for MM/DD/YYYY format")
	}
	if v.First("date3") == "" {
		t.Error("Expected error for invalid date")
	}
}

func TestValidatorErrors(t *testing.T) {
	data := map[string]interface{}{}
	rules := map[string][]string{
		"email": {"required", "email"},
		"name":  {"required"},
	}

	v := Make(data, rules)

	errors := v.Errors()
	if len(errors) != 2 {
		t.Errorf("Expected 2 field errors, got: %d", len(errors))
	}
}

func TestValidatorAll(t *testing.T) {
	data := map[string]interface{}{}
	rules := map[string][]string{
		"email": {"required"},
		"name":  {"required"},
	}

	v := Make(data, rules)

	all := v.All()
	if len(all) != 2 {
		t.Errorf("Expected 2 total errors, got: %d", len(all))
	}
}

func TestValidatorFirst(t *testing.T) {
	data := map[string]interface{}{
		"email": "invalid",
	}
	rules := map[string][]string{
		"email": {"email"},
	}

	v := Make(data, rules)

	first := v.First("email")
	if first == "" {
		t.Error("Expected First to return error message")
	}
}

func TestValidatorFirstEmpty(t *testing.T) {
	data := map[string]interface{}{
		"email": "test@example.com",
	}
	rules := map[string][]string{
		"email": {"email"},
	}

	v := Make(data, rules)

	first := v.First("email")
	if first != "" {
		t.Error("Expected First to return empty for valid field")
	}
}

func TestValidationMultipleRules(t *testing.T) {
	data := map[string]interface{}{
		"email": "a",
	}
	rules := map[string][]string{
		"email": {"required", "email", "min:5"},
	}

	v := Make(data, rules)

	errors := v.Errors()["email"]
	if len(errors) < 2 {
		t.Error("Expected multiple errors for multiple failed rules")
	}
}
