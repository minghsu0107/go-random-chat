package uploader

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func newObjectKey(channelID uint64, extension string) string {
	return joinStrs(strconv.FormatUint(channelID, 10), "/", uuid.New().String(), extension)
}

func getChannelIDFromObjectKey(objectKey string) (uint64, error) {
	channelIDStr := strings.Split(objectKey, "/")[0]
	channelID, err := strconv.ParseUint(channelIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse channel ID from object key: %v, error: %v", objectKey, err)
	}
	return channelID, nil
}

func joinStrs(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
