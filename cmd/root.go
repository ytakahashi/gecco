package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gecco",
	Short: "A Command Line Tool To Oprtate AWS EC2.",
	Long:  "A Command Line Tool To Oprtate AWS EC2.",
}

// Execute command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")

	listCmd.Flags().StringVarP(&listOpts.tagKey, "tagKey", "", "", "filters by tag key")
	listCmd.Flags().StringVarP(&listOpts.tagValue, "tagValue", "", "", "filters by tag value")
	listCmd.Flags().StringVarP(&listOpts.status, "status", "", "", "filters by tag value")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".gecco")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
