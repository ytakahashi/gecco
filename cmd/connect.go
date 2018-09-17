package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/config"
)

var connectOpts = &config.ConnectOptions{}

func newConnectCmd() *cobra.Command {
	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "connect to EC2 instance",
		Long:  "connect to EC2 instance using 'aws cli start-session' command",
		Run: func(cmd *cobra.Command, args []string) {
			err := connect(*connectOpts, doRun)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	connectCmd.Flags().StringVarP(&connectOpts.Target, "target", "", "", "target instanceId to start session")

	return connectCmd
}

func connect(options config.ConnectOptions, run func(string) error) error {
	if options.Target == "" {
		return errors.New("Option '--target' is not specified")
	}

	return run(connectOpts.Target)
}

func createCommand(target string) *exec.Cmd {
	command := exec.Command("aws", "ssm", "start-session", "--target", target)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command
}

func doRun(target string) error {
	return createCommand(target).Run()
}
