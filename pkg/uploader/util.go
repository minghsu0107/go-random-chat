package uploader

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func newObjectKey(channelID uint64, extension string) string {
	return joinStrs(strconv.FormatUint(channelID, 10), "/", uuid.New().String(), extension)
}

func joinStrs(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
