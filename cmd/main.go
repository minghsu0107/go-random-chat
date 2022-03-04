package main

import (
	"log"
	"os"

	"github.com/minghsu0107/go-random-chat/cmd/chat"
	"github.com/minghsu0107/go-random-chat/cmd/upload"
	"github.com/minghsu0107/go-random-chat/cmd/web"
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
	case "web":
		web.RunWebServer()
	default:
		log.Fatalf("invalid app name: %s. Should be 'chat', 'upload'", app)
	}
}
