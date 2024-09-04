package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

const EN_UPPER = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const EN_LOWER = "abcdefghijklmnopqrstuvwxyz"
const EN_DIGIT = "0123456789"

const base36 = "0123456789abcdefghijklmnopqrstuvwxyz"

func ConvertTo36Base(number int64) string {
	if number <= 0 {
		return ""
	}

	var lst []rune
	for number > 0 {
		lst = append(lst, rune(base36[int(number%36)]))
		number = number / 36
	}

	sb := strings.Builder{}
	for i := len(lst) - 1; i >= 0; i-- {
		sb.WriteRune(lst[i])
	}

	return sb.String()
}

func ConvertToBase10From36(base36String string) int64 {
	if base36String == "" {
		return -1
	}

	base10Number := int64(0)
	base := int64(36)
	str := strings.ToLower(base36String)

	for _, char := range str {
		digit := strings.IndexRune(base36, char)
		if digit == -1 {
			// Character not found in base36, handle error or return 0
			return -1
		}
		base10Number = base10Number*base + int64(digit)
	}

	return base10Number
}

func RandomString(length int, upper bool, lower bool, digits bool) string {
	if !upper && !lower && !digits {
		return ""
	}

	tmpString := ""
	if upper {
		tmpString += EN_UPPER
	}

	if lower {
		tmpString += EN_LOWER
	}

	if digits {
		tmpString += EN_DIGIT
	}

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	l := len(tmpString)
	for i := range b {
		b[i] = tmpString[seededRand.Intn(l)]
	}
	return string(b)
}

func RandomInt(maxN int) int {
	// rand.Seed(time.Now().UnixNano())
	// Generate a random integer between 0 and 99
	randomNumber := rand.Intn(maxN)
	return randomNumber
}

func HashMd5(input string) string {
	// MD5 Hashing
	md5Hash := md5.Sum([]byte(input))
	md5HashString := hex.EncodeToString(md5Hash[:])
	return md5HashString
}

func HashSha1(input string) string {
	// SHA-1 Hashing
	sha1Hash := sha1.Sum([]byte(input))
	sha1HashString := hex.EncodeToString(sha1Hash[:])
	return sha1HashString
}

func HashSha256(input string) string {
	// SHA256 Hashing
	sha256Hash := sha256.Sum256([]byte(input))
	sha1HashString := hex.EncodeToString(sha256Hash[:])
	return sha1HashString
}
