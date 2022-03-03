package main

import (
	"log"
	"os"

	"github.com/minghsu0107/go-random-chat/cmd/chat"
	"github.com/minghsu0107/go-random-chat/cmd/upload"
)

var (
	app = os.Getenv("APP")
)

func main() {
	switch app {
	case "chat":
		chat.RunChatServer()
	case "upload":
		upload.RunUploadServer()
	default:
		log.Fatalf("invalid app name: %s. Should be 'chat', 'upload'", app)
	}
}
