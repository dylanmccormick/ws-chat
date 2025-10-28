package wschat

import (
	"fmt"

	"github.com/dylanmccormick/ws-chat/cmd/client"
	"github.com/dylanmccormick/ws-chat/cmd/client/tui"
	"github.com/dylanmccormick/ws-chat/cmd/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ws-chat",
	Short: "this is my ws-chat program",
	Long:  `Here is a long description. Don't read it`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(replCmd)
	rootCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(startTui)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of ws-chat",
	Long:  `This is the long desc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ws-chat version 0.0.1")
	},
}

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "start the repl for the server.",
	Long: `A repl for testing the web socket chat without having to launch the whole client.
	Very basic and does not get real time updates to chat messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		client.StartREPL()
	},
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "a command to start the server",
	Long:  `Will update these later with some polish`,
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer()
	},
}

var startTui = &cobra.Command{
	Use:   "tui",
	Short: "a command to start the client tui",
	Long:  `Will update these later with some polish`,
	Run: func(cmd *cobra.Command, args []string) {
		tui.Start()
	},
}
