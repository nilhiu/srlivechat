package cmd

import (
	"github.com/nilhiu/srlivechat/client"
	"github.com/nilhiu/srlivechat/tui"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	username   string
	serverHost string
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connects to a srlivechat server",
	Long:  `connects to a srlivechat server started by the "srlivechat start" command.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New(serverHost, username)
		if err != nil {
			log.Fatal().Msgf("client could not be created, %s", err.Error())
		}
		defer c.Close()

		tui.New(cmd.Context(), c).Run()
	},
}

func init() {
	connectCmd.PersistentFlags().
		StringVarP(&username, "username", "u", "default_user", "the username to use in the live chat")
	connectCmd.PersistentFlags().
		StringVar(&serverHost, "host", "localhost:3000", "the server address to connect to")
	rootCmd.AddCommand(connectCmd)
}
