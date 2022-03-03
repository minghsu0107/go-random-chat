package chat

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	rc "github.com/minghsu0107/go-random-chat"
	log "github.com/sirupsen/logrus"
)

func RunChatServer() {
	router, err := rc.InitializeRouter()
	if err != nil {
		log.Fatal(err)
	}

	router.Run()

	done := make(chan bool, 1)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		router.GracefulStop(ctx, done)
	}()

	<-done
}
