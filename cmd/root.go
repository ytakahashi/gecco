package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ytakahashi/gecco/config"
)

var conf config.Config

var cfgFile string

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gecco",
		Short: "A Command Line Tool To Oprtate AWS EC2.",
		Long:  "A Command Line Tool To Oprtate AWS EC2.",
	}

	rootCmd.AddCommand(newListCmd())
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

	err = viper.ReadInConfig()
	if err != nil {
		return
		// panic(fmt.Errorf("error config file: %s", err))
	}

	if err = viper.Unmarshal(&conf); err != nil {
		return
		// fmt.Println(err)
		// os.Exit(1)
	}
	return

	// if cfgFile != "" {
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}

	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(".gecco")
	// }

	// if err := viper.ReadInConfig(); err != nil {
	// 	fmt.Println("Can't read config:", err)
	// 	os.Exit(1)
	// }
}
