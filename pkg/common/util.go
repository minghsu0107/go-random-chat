package common

import (
	"os"
	"strings"
)

func GetServerAddrs(addrs string) []string {
	return strings.Split(addrs, ",")
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
