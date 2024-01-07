package cmd

import (
	log "log/slog"
	"os"

	"github.com/minghsu0107/go-random-chat/internal/wire"
	"github.com/spf13/cobra"
)

var forwarderCmd = &cobra.Command{
	Use:   "forwarder",
	Short: "forwarder server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := wire.InitializeForwarderServer("forwarder")
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(forwarderCmd)
}
