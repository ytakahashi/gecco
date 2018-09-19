package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ytakahashi/gecco/config"
)

var cfgFile string

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gecco",
		Short: "A Command Line Tool To Oprtate AWS EC2.",
		Long:  "A Command Line Tool To Oprtate AWS EC2.",
	}

	rootCmd.AddCommand(newListCmd(&listCommand{}))
	rootCmd.AddCommand(newConnectCmd())

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

func initConfig() (err error) {
	viper.SetConfigName("gecco")
	viper.AddConfigPath("$HOME/.config")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.Conf); err != nil {
		return
	}

	return
}
