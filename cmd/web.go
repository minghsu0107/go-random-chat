package cmd

import (
	"github.com/minghsu0107/go-random-chat/internal/wire"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "web server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := wire.InitializeWebServer("web")
		if err != nil {
			log.Fatal(err)
		}
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}
