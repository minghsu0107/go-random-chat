package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minghsu0107/go-random-chat/internal/wire"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var uploaderCmd = &cobra.Command{
	Use:   "uploader",
	Short: "uploader server",
	Run: func(cmd *cobra.Command, args []string) {
		runUploadServer()
	},
}

func runUploadServer() {
	router, err := wire.InitializeUploaderRouter()
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

func init() {
	rootCmd.AddCommand(uploaderCmd)
}
