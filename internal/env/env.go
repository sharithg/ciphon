package env

import (
	"log"
	"os"
	"strconv"
)

func GetString(key string, fail bool, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		if fail {
			log.Fatalf("environment variable %s not set", key)
		}
		return fallback
	}

	return val
}

func GetInt(key string, fail bool, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		if fail {
			log.Fatalf("environment variable %s not set", key)
		}
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		if fail {
			log.Fatalf("environment variable %s is not a valid integer", key)
		}
		return fallback
	}

	return valAsInt
}

func GetBool(key string, fail bool, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		if fail {
			log.Fatalf("environment variable %s not set", key)
		}
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		if fail {
			log.Fatalf("environment variable %s is not a valid boolean", key)
		}
		return fallback
	}

	return boolVal
}
