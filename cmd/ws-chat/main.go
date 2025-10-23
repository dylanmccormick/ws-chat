package wschat

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "ws-chat",
	Short: "this is my ws-chat program",
	Long: `Here is a long description. Don't read it`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Print the version of ws-chat",
	Long: `This is the long desc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ws-chat version 0.0.1")
	},
}

