package database

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func findMissingEnvValues(required map[string]string) []string {
	missing := make([]string, 0)
	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)
	return missing
}

func envInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func envDuration(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	value, err := time.ParseDuration(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
