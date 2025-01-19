package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nilhiu/srlivechat/client"
	"github.com/nilhiu/srlivechat/server"
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

		log.Info().Msgf("connected to srlivechat server at %s as %s", serverHost, username)

		go messageHandler(c)

		for {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if scanner.Err() != nil {
				log.Error().Msgf("could not read written message, %s", err.Error())
			}

			if err := c.Write(scanner.Text()); err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	connectCmd.PersistentFlags().
		StringVarP(&username, "username", "u", "default_user", "the username to use in the live chat")
	connectCmd.PersistentFlags().
		StringVar(&serverHost, "host", "localhost:3000", "the server address to connect to")
	rootCmd.AddCommand(connectCmd)
}

func messageHandler(c *client.Client) {
	for {
		msg, err := c.Read()
		if err != nil {
			log.Fatal().Msgf("failed to read from server, %s", err.Error())
		}

		switch msg.Type() {
		case server.UserMessage:
			fmt.Printf("[%s]: %s\n", msg.Sender(), msg.Message())
		case server.ServerMessage:
			fmt.Printf("<SERVER>: %s\n", msg.Message())
		case server.ConnectMessage:
			fmt.Printf("<CONNECTED>: %s\n", msg.Sender())
		case server.DisconnectMessage:
			fmt.Printf("<DISCONNECTED>: %s\n", msg.Sender())
		}
	}
}
