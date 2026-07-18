package utils

import (
	"testing"
	"time"
)

func TestShortDate(t *testing.T) {
	got, err := ShortDate("2024-01-06")
	if err != nil {
		t.Fatalf("ShortDate returned error: %v", err)
	}
	if got != "20240106" {
		t.Errorf("ShortDate = %q, want %q", got, "20240106")
	}

	if _, err := ShortDate("not-a-date"); err == nil {
		t.Error("ShortDate(invalid) expected an error, got nil")
	}
}

func TestShortDateInt(t *testing.T) {
	got, err := ShortDateInt("2024-01-06")
	if err != nil {
		t.Fatalf("ShortDateInt returned error: %v", err)
	}
	if got != 20240106 {
		t.Errorf("ShortDateInt = %d, want %d", got, 20240106)
	}

	if _, err := ShortDateInt("bad"); err == nil {
		t.Error("ShortDateInt(invalid) expected an error, got nil")
	}
}

func TestTime2Unix(t *testing.T) {
	got, err := Time2Unix("2024-01-01 00:00:00", TIME_FORMAT_TS)
	if err != nil {
		t.Fatalf("Time2Unix returned error: %v", err)
	}
	want := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	if got != want {
		t.Errorf("Time2Unix = %d, want %d", got, want)
	}

	if _, err := Time2Unix("bad", TIME_FORMAT_TS); err == nil {
		t.Error("Time2Unix(invalid) expected an error, got nil")
	}
}

func TestTimeFormatCheck(t *testing.T) {
	if _, err := TimeFormatCheck("2024-01-06", "2006-01-02"); err != nil {
		t.Errorf("TimeFormatCheck valid input returned error: %v", err)
	}
	if _, err := TimeFormatCheck("2024/01/06", "2006-01-02"); err == nil {
		t.Error("TimeFormatCheck mismatched layout expected an error, got nil")
	}
}

func TestFilterAlphanumeric(t *testing.T) {
	cases := map[string]string{
		"a-b_c 1!2": "abc12",
		"hello":     "hello",
		"@#$%":      "",
	}
	for in, want := range cases {
		if got := FilterAlphanumeric(in); got != want {
			t.Errorf("FilterAlphanumeric(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCalculateAge(t *testing.T) {
	// Someone whose birthday was one day ago, 30 years back, is exactly 30.
	dob := time.Now().AddDate(-30, 0, -1)
	age, err := CalculateAge(dob.Format("2006-01-02"), "2006-01-02")
	if err != nil {
		t.Fatalf("CalculateAge returned error: %v", err)
	}
	if age != 30 {
		t.Errorf("CalculateAge = %d, want 30", age)
	}

	if _, err := CalculateAge("bad-date", "2006-01-02"); err == nil {
		t.Error("CalculateAge(invalid) expected an error, got nil")
	}
}
