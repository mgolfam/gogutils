package utils

import (
	"testing"

	"github.com/mgolfam/gogutils/enums"
)

func TestIsPasswordEqual(t *testing.T) {
	const plain = "s3cret"

	cases := []struct {
		name     string
		dbPass   string
		input    string
		passType int
		want     bool
	}{
		{"plain match", plain, plain, enums.PASS_TYPE_PLAIN, true},
		{"plain mismatch", plain, "wrong", enums.PASS_TYPE_PLAIN, false},
		{"md5 match", HashMd5(plain), plain, enums.PASS_TYPE_MD5, true},
		{"md5 mismatch", HashMd5(plain), "wrong", enums.PASS_TYPE_MD5, false},
		{"sha256 match", HashSha256(plain), plain, enums.PASS_TYPE_SHA256, true},
		{"sha256 mismatch", HashSha256(plain), "wrong", enums.PASS_TYPE_SHA256, false},
		{"unknown type", plain, plain, 99, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsPasswordEqual(c.dbPass, c.input, c.passType); got != c.want {
				t.Errorf("IsPasswordEqual(%q, %q, %d) = %v, want %v",
					c.dbPass, c.input, c.passType, got, c.want)
			}
		})
	}
}
