package utils

import (
	"fmt"
	"strconv"
	"time"
	"unicode"

	"github.com/mgolfam/gogutils/crypt"
)

const (
	TIME_FORMAT_TS   = "2006-01-02 15:04:05"
	TIME_FORMAT_DATE = "2006-01-02"
	TIME_FORMAT_TIME = "15:04:05"
)

func NowUnixSeconds() int64 {
	currentTime := time.Now()
	return currentTime.Unix()
}

func Today() string {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")
	return currentDate
}

func TodayTime() time.Time {
	return time.Now()
}

func TodayZeroClockTime() time.Time {
	today, _ := TimeFormatCheck(Today(), "2006-01-02")

	return today
}

func Now(format string) string {
	currentTime := time.Now()
	currentDate := currentTime.Format(format)
	return currentDate
}

func TodayBase36() string {
	currentTime := time.Now()
	currentDate := currentTime.Format("060102")
	for i := 0; i < 5; i++ {
		if len(currentDate) > 1 && currentDate[0] == '0' {
			// Remove the '0' from the first character
			currentDate = currentDate[1:]
		} else {
			break
		}
	}

	number, _ := strconv.ParseInt(currentDate, 10, 64)
	base36 := crypt.EncodeBase36(number)
	return base36
}

func TimeFormatCheck(timestr, format string) (time.Time, error) {
	// Parse the date string
	parsed, err := time.Parse(format, timestr)
	return parsed, err
}

func ShortDate(dateString string) (string, error) {
	// Parse the date string
	parsedDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return "", err
	}

	// Format the date as "231006"
	formattedDate := parsedDate.Format("20060102")
	return formattedDate, nil
}

func ShortDateInt(dateString string) (int, error) {
	// Parse the date string
	parsedDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return -1, err
	}

	// Format the date as "231006"
	formattedDate := parsedDate.Format("20060102")
	intNumber, err := strconv.Atoi(formattedDate)

	return intNumber, err
}

func CalculateAge(dateString, dateFormat string) (int, error) {
	// Parse the date string into a time.Time object
	dob, err := time.Parse(dateFormat, dateString)
	if err != nil {
		return 0, err
	}

	// Get the current date
	currentDate := time.Now()

	// Calculate the age
	age := currentDate.Year() - dob.Year()
	if currentDate.YearDay() < dob.YearDay() {
		age--
	}

	return age, nil
}

func FilterAlphanumeric(input string) string {
	result := make([]rune, 0, len(input))
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result = append(result, r)
		}
	}
	return string(result)
}

func Time2Unix(timeString, layout string) (int64, error) {
	// Parse the string into a time.Time object
	parsedTime, err := time.Parse(layout, timeString)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return 0, err
	}

	// Convert the time to Unix timestamp (seconds)
	unixTime := parsedTime.Unix()

	return unixTime, nil
}

// Duration helpers

// ParseDurationSafe wraps time.ParseDuration, returning 0 on error.
func ParseDurationSafe(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}

// FormatDurationShort formats a duration like "1h2m3s".
func FormatDurationShort(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	seconds := int64(d.Seconds())
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// Time window helpers

// StartOfDay returns the start of the day for t in its location.
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day for t in its location.
func EndOfDay(t time.Time) time.Time {
	return StartOfDay(t).Add(24*time.Hour - time.Nanosecond)
}

// StartOfWeek returns the start of the week (Monday) for t.
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return StartOfDay(t).AddDate(0, 0, -weekday+1)
}

// EndOfWeek returns the end of the week (Sunday) for t.
func EndOfWeek(t time.Time) time.Time {
	return StartOfWeek(t).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// StartOfMonth returns the start of the month for t.
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of the month for t.
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// ExponentialBackoff returns a backoff duration for a given attempt (0-based).
// Example: base=100ms, factor=2, max=5s.
func ExponentialBackoff(base time.Duration, factor float64, max time.Duration, attempt int) time.Duration {
	if attempt <= 0 {
		return base
	}
	backoff := float64(base)
	for i := 0; i < attempt; i++ {
		backoff *= factor
		if time.Duration(backoff) >= max {
			return max
		}
	}
	return time.Duration(backoff)
}
