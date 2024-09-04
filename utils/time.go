package utils

import (
	"strconv"
	"time"
	"unicode"

	"github.com/mgolfam/gogutils/crypt"
)

const (
	TIME_FORMAT_TS   = "2006-01-02 15:04:05"
	TIME_FORMAT_DATE = "2006-01-02 15:04:05"
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
	currentTime := time.Now()
	return currentTime
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
