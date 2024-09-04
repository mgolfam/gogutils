package utils

import (
	"strings"

	"github.com/google/uuid"
)

func UUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func CleanUUID() string {
	return strings.Replace(UUID(), "-", "", -1)
}
