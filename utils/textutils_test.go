package utils

import "testing"

func TestStringExists(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if !StringExists("b", slice) {
		t.Error("StringExists returned false for present element")
	}
	if StringExists("z", slice) {
		t.Error("StringExists returned true for absent element")
	}
	if StringExists("a", nil) {
		t.Error("StringExists returned true for nil slice")
	}
}

func TestSubstring(t *testing.T) {
	cases := []struct {
		text string
		max  int
		want string
	}{
		{"hello", 3, "hel"},
		{"hello", 5, "hello"},
		{"hi", 10, "hi"},
		{"", 3, ""},
	}
	for _, c := range cases {
		if got := Substring(c.text, c.max); got != c.want {
			t.Errorf("Substring(%q, %d) = %q, want %q", c.text, c.max, got, c.want)
		}
	}
}

func TestAtoi(t *testing.T) {
	if got := Atoi("42"); got != 42 {
		t.Errorf("Atoi(\"42\") = %d, want 42", got)
	}
	// Non-numeric input returns the sentinel -100.
	if got := Atoi("nope"); got != -100 {
		t.Errorf("Atoi(\"nope\") = %d, want -100", got)
	}
}

func TestAtoiPtr(t *testing.T) {
	if got := AtoiPtr("7"); got == nil || *got != 7 {
		t.Errorf("AtoiPtr(\"7\") = %v, want *7", got)
	}
	if got := AtoiPtr("nan"); got != nil {
		t.Errorf("AtoiPtr(\"nan\") = %v, want nil", got)
	}
}

func TestCaseConverters(t *testing.T) {
	if got := PascalCase("hello"); got != "Hello" {
		t.Errorf("PascalCase(\"hello\") = %q, want %q", got, "Hello")
	}
	if got := SnakeCase("helloWorld"); got != "hello_world" {
		t.Errorf("SnakeCase(\"helloWorld\") = %q, want %q", got, "hello_world")
	}
	if got := KebabCase("helloWorld"); got != "hello-world" {
		t.Errorf("KebabCase(\"helloWorld\") = %q, want %q", got, "hello-world")
	}
	if got := CamelCase("hello_world"); got != "helloWorld" {
		t.Errorf("CamelCase(\"hello_world\") = %q, want %q", got, "helloWorld")
	}
}

func TestCheckStringCategory(t *testing.T) {
	cases := map[string]string{
		"12345":    "all_digits",
		"abcDEF":   "all_alphabets",
		"abc123":   "mixed",
		"abc 123":  "invalid", // space is neither letter nor digit
		"user@x":   "invalid",
		"":         "invalid",
	}
	for in, want := range cases {
		if got := CheckStringCategory(in); got != want {
			t.Errorf("CheckStringCategory(%q) = %q, want %q", in, got, want)
		}
	}
}
