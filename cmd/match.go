package cmd

import (
	"github.com/minghsu0107/go-random-chat/internal/wire"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "match server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := wire.InitializeMatchServer("match")
		if err != nil {
			log.Fatal(err)
		}
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(matchCmd)
}
