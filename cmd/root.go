package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gecco",
		Short: "A Command Line Tool To Oprtate AWS EC2.",
		Long:  "A Command Line Tool To Oprtate AWS EC2.",
	}

	rootCmd.AddCommand(newListCmd(&listCommand{}))
	rootCmd.AddCommand(newConnectCmd(&connectCommand{}))

	return rootCmd
}

// Execute command
func Execute() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
