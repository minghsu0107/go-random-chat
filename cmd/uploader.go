package cmd

import (
	"github.com/minghsu0107/go-random-chat/internal/wire"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var uploaderCmd = &cobra.Command{
	Use:   "uploader",
	Short: "uploader server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := wire.InitializeUploaderServer("uploader")
		if err != nil {
			log.Fatal(err)
		}
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(uploaderCmd)
}
