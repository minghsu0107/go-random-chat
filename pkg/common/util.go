package common

import (
	"strings"
)

func GetServerAddrs(addrs string) []string {
	return strings.Split(addrs, ",")
}
