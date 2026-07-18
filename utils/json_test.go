package utils

import (
	"encoding/json"
	"testing"
)

func TestFlattenJSON(t *testing.T) {
	nested := map[string]interface{}{
		"name": "root",
		"child": map[string]interface{}{
			"a": 1,
			"deep": map[string]interface{}{
				"b": "x",
			},
		},
	}

	flat := FlattenJSON(nested)

	want := map[string]interface{}{
		"name":         "root",
		"child.a":      1,
		"child.deep.b": "x",
	}

	if len(flat) != len(want) {
		t.Fatalf("FlattenJSON produced %d keys, want %d: %#v", len(flat), len(want), flat)
	}
	for k, v := range want {
		got, ok := flat[k]
		if !ok {
			t.Errorf("FlattenJSON missing key %q", k)
			continue
		}
		if got != v {
			t.Errorf("FlattenJSON[%q] = %v, want %v", k, got, v)
		}
	}
}

func TestToJsonString(t *testing.T) {
	in := map[string]int{"a": 1}
	out := ToJsonString(in)

	// Round-trip back to a map to avoid asserting on exact whitespace.
	var back map[string]int
	if err := json.Unmarshal([]byte(out), &back); err != nil {
		t.Fatalf("ToJsonString produced invalid JSON %q: %v", out, err)
	}
	if back["a"] != 1 {
		t.Errorf("ToJsonString round trip = %v, want map[a:1]", back)
	}
}

func TestToJsonStringUnmarshalable(t *testing.T) {
	// Channels cannot be marshaled; the function must return "" rather than panic.
	if got := ToJsonString(make(chan int)); got != "" {
		t.Errorf("ToJsonString(chan) = %q, want empty string", got)
	}
}
