package cmd

import (
	log "log/slog"
	"os"

	"github.com/minghsu0107/go-random-chat/internal/wire"
	"github.com/spf13/cobra"
)

var uploaderCmd = &cobra.Command{
	Use:   "uploader",
	Short: "uploader server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := wire.InitializeUploaderServer("uploader")
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(uploaderCmd)
}
