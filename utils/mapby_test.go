package utils

import "testing"

func TestFillFromMapByTags(t *testing.T) {
	type Src struct {
		FullName string
		Years    int
	}
	type Dst struct {
		Name string `mapby:"FullName"`
		Age  int    `mapby:"Years"`
	}

	src := Src{FullName: "Ada", Years: 36}
	var dst Dst
	if err := FillFromMapByTags(&dst, src); err != nil {
		t.Fatalf("FillFromMapByTags returned error: %v", err)
	}
	if dst.Name != "Ada" || dst.Age != 36 {
		t.Errorf("FillFromMapByTags produced %+v, want {Name:Ada Age:36}", dst)
	}
}

func TestFillFromMapByTagsMissingField(t *testing.T) {
	type Src struct{ A string }
	type Dst struct {
		B string `mapby:"DoesNotExist"`
	}
	var dst Dst
	if err := FillFromMapByTags(&dst, Src{A: "x"}); err == nil {
		t.Error("FillFromMapByTags expected an error for a missing source field, got nil")
	}
}

func TestFillByFieldNameAndType(t *testing.T) {
	type Src struct {
		ID   int
		Name string
		Skip float64 // different type on dest -> ignored
	}
	type Dst struct {
		ID      int
		Name    string
		Skip    string
		Missing bool // not in src -> untouched
	}

	src := Src{ID: 7, Name: "go", Skip: 1.5}
	dst := Dst{Skip: "keep", Missing: true}
	if err := FillByFieldNameAndType(&dst, src); err != nil {
		t.Fatalf("FillByFieldNameAndType returned error: %v", err)
	}
	if dst.ID != 7 || dst.Name != "go" {
		t.Errorf("FillByFieldNameAndType did not copy matching fields: %+v", dst)
	}
	if dst.Skip != "keep" {
		t.Errorf("FillByFieldNameAndType overwrote a type-mismatched field: Skip=%q", dst.Skip)
	}
	if !dst.Missing {
		t.Errorf("FillByFieldNameAndType touched a field absent from source")
	}
}

func TestFillByFieldNameAndTypeSlice(t *testing.T) {
	type Src struct{ N int }
	type Dst struct{ N int }

	srcs := []Src{{N: 1}, {N: 2}, {N: 3}}
	dsts := make([]Dst, 3)
	if err := FillByFieldNameAndTypeSlice(dsts, srcs); err != nil {
		t.Fatalf("FillByFieldNameAndTypeSlice returned error: %v", err)
	}
	for i, d := range dsts {
		if d.N != srcs[i].N {
			t.Errorf("element %d = %d, want %d", i, d.N, srcs[i].N)
		}
	}

	// Length mismatch must error.
	if err := FillByFieldNameAndTypeSlice(make([]Dst, 1), srcs); err == nil {
		t.Error("FillByFieldNameAndTypeSlice expected an error on length mismatch, got nil")
	}
}

func TestFillNonStructError(t *testing.T) {
	var dst struct{ A int }
	notAStruct := 5
	if err := FillFromMapByTags(&dst, notAStruct); err == nil {
		t.Error("FillFromMapByTags expected an error for non-struct source, got nil")
	}
}
