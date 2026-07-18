package crypt

import "testing"

func TestSha256(t *testing.T) {
	// Known SHA-256 vector for the input "hello".
	const want = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got := Sha256("hello"); got != want {
		t.Errorf("Sha256(%q) = %q, want %q", "hello", got, want)
	}

	// The empty string has a well-known digest too.
	const wantEmpty = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if got := Sha256(""); got != wantEmpty {
		t.Errorf("Sha256(%q) = %q, want %q", "", got, wantEmpty)
	}
}

func TestEncodeDecodeBase36RoundTrip(t *testing.T) {
	cases := []int64{0, 1, 35, 36, 1234567890, 9223372036854775807}
	for _, n := range cases {
		enc := EncodeBase36(n)
		if got := DecondeBase36(enc); got != n {
			t.Errorf("DecondeBase36(EncodeBase36(%d)) = %d, want %d (encoded=%q)", n, got, n, enc)
		}
	}
}

func TestEncodeBase36Known(t *testing.T) {
	cases := map[int64]string{
		0:  "0",
		35: "z",
		36: "10",
	}
	for in, want := range cases {
		if got := EncodeBase36(in); got != want {
			t.Errorf("EncodeBase36(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestDecodeBase36Invalid(t *testing.T) {
	// '.' is not a valid base-36 digit, so decoding must fail and return -1.
	if got := DecondeBase36("!!not-base36!!"); got != -1 {
		t.Errorf("DecondeBase36(invalid) = %d, want -1", got)
	}
}
