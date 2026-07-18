package utils

import (
	"regexp"
	"strings"
	"testing"
)

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func TestUUIDFormat(t *testing.T) {
	id := UUID()
	if !uuidRe.MatchString(id) {
		t.Errorf("UUID() = %q, not a canonical UUID", id)
	}
}

func TestCleanUUID(t *testing.T) {
	id := CleanUUID()
	if len(id) != 32 {
		t.Errorf("CleanUUID() length = %d, want 32", len(id))
	}
	if strings.Contains(id, "-") {
		t.Errorf("CleanUUID() = %q, must not contain dashes", id)
	}
}

func TestUUIDUniqueness(t *testing.T) {
	seen := make(map[string]struct{}, 1000)
	for i := 0; i < 1000; i++ {
		id := UUID()
		if _, dup := seen[id]; dup {
			t.Fatalf("UUID() produced a duplicate: %q", id)
		}
		seen[id] = struct{}{}
	}
}
