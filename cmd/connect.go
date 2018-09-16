package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/config"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect to EC2 instance",
	Long:  "connect to EC2 instance using aws cli start-session",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Connect")

		err := connect()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

var connectOpts = &config.ConnectOptions{}

func connect() error {
	if connectOpts.Target == "" {
		return errors.New("Target is not specified")
	}

	command := exec.Command("aws", "ssm", "start-session", "--target", connectOpts.Target)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
