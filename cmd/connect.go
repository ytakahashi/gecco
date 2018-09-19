package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
	"github.com/ytakahashi/gecco/ext"
)

var connectOpts = &config.ConnectOptions{}

func newConnectCmd() *cobra.Command {
	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "connect to EC2 instance",
		Long:  "connect to EC2 instance using 'aws cli start-session' command",
		Run: func(cmd *cobra.Command, args []string) {
			e := aws.Ec2{}
			err := connect(*connectOpts, doRun, initConfig, e)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	connectCmd.Flags().StringVarP(&connectOpts.Target, "target", "", "", "target instanceId to start session")
	connectCmd.Flags().BoolVarP(&connectOpts.Interactive, "interactive", "i", false, "Select a value interactively (requires config file)")

	return connectCmd
}

func connect(
	options config.ConnectOptions,
	run func(string) error,
	init func() error,
	client aws.Ec2Client,
) error {
	var target string
	if options.Interactive {
		if err := init(); err != nil {
			return err
		}

		instances, err := client.GetInstances(config.ListOption{})
		if err != nil {
			return err
		}

		filter := ext.Command{
			Args: []string{config.Conf.InteractiveFilterCommand},
		}

		target, err = instances.GetFilteredInstances(filter)
		if err != nil {
			return err
		}
	} else {
		if options.Target == "" {
			return errors.New("Option '--target' is not specified")
		}
		target = options.Target
	}
	return run(target)

}

var startSession = ext.Command{
	Args: []string{"aws", "ssm", "start-session", "--target"},
}

func doRun(target string) error {
	startSession.Args = append(startSession.Args, target)
	return startSession.CreateCommand(os.Stdin, os.Stdout, os.Stderr).Run()
}
