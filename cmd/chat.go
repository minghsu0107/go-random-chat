package cmd

import (
	log "log/slog"
	"os"

	"github.com/minghsu0107/go-random-chat/internal/wire"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "chat server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := wire.InitializeChatServer("chat")
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
