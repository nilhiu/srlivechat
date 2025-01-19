package cmd

import (
	"github.com/nilhiu/srlivechat/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var port string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts a srlivechat server",
	Long:  `starts a srlivechat users can connect to using the "srlivechat connect" command.`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.New(port)

		log.Info().Msgf("starting srlivechat server at port %s", port)

		if err := srv.Run(); err != nil {
			log.Fatal().Msgf("server error: %s", err.Error())
		}
	},
}

func init() {
	startCmd.PersistentFlags().
		StringVarP(&port, "port", "p", ":3000", "specifies the port to run the server")
	rootCmd.AddCommand(startCmd)
}
