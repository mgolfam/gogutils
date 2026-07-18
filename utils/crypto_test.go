package utils

import (
	"strings"
	"testing"
)

func TestHashKnownVectors(t *testing.T) {
	cases := []struct {
		name string
		fn   func(string) string
		in   string
		want string
	}{
		{"md5", HashMd5, "hello", "5d41402abc4b2a76b9719d911017c592"},
		{"sha1", HashSha1, "hello", "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{"sha256", HashSha256, "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.fn(c.in); got != c.want {
				t.Errorf("%s(%q) = %q, want %q", c.name, c.in, got, c.want)
			}
		})
	}
}

func TestConvertBase36RoundTrip(t *testing.T) {
	for _, n := range []int64{1, 35, 36, 71, 1234567890} {
		enc := ConvertTo36Base(n)
		if got := ConvertToBase10From36(enc); got != n {
			t.Errorf("round trip for %d failed: encoded=%q decoded=%d", n, enc, got)
		}
	}
}

func TestConvertTo36BaseNonPositive(t *testing.T) {
	// The implementation returns an empty string for values <= 0.
	for _, n := range []int64{0, -1, -100} {
		if got := ConvertTo36Base(n); got != "" {
			t.Errorf("ConvertTo36Base(%d) = %q, want empty string", n, got)
		}
	}
}

func TestConvertToBase10From36Invalid(t *testing.T) {
	if got := ConvertToBase10From36(""); got != -1 {
		t.Errorf("ConvertToBase10From36(empty) = %d, want -1", got)
	}
	if got := ConvertToBase10From36("!!!"); got != -1 {
		t.Errorf("ConvertToBase10From36(invalid) = %d, want -1", got)
	}
}

func TestRandomStringLengthAndCharset(t *testing.T) {
	const n = 64
	s := RandomString(n, false, true, false) // lowercase only
	if len(s) != n {
		t.Fatalf("RandomString length = %d, want %d", len(s), n)
	}
	for _, r := range s {
		if !strings.ContainsRune(EN_LOWER, r) {
			t.Fatalf("RandomString produced out-of-charset rune %q", r)
		}
	}

	// Digits-only path.
	d := RandomString(n, false, false, true)
	for _, r := range d {
		if !strings.ContainsRune(EN_DIGIT, r) {
			t.Fatalf("digits RandomString produced non-digit rune %q", r)
		}
	}
}

func TestRandomStringNoCharsetsReturnsEmpty(t *testing.T) {
	if got := RandomString(10, false, false, false); got != "" {
		t.Errorf("RandomString with no charsets = %q, want empty", got)
	}
}

func TestRandomIntBounds(t *testing.T) {
	for i := 0; i < 1000; i++ {
		v := RandomInt(10)
		if v < 0 || v >= 10 {
			t.Fatalf("RandomInt(10) = %d, out of [0,10)", v)
		}
	}
}
