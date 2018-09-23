package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
	"github.com/ytakahashi/gecco/ext"
)

var connectOpts = &config.ConnectOption{}

func newConnectCmd(command iConnectCommand) *cobra.Command {
	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "connect to EC2 instance",
		Long:  "connect to EC2 instance using 'aws cli start-session' command",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			err = command.initConnectCommand(*connectOpts, aws.Ec2{}, &config.Config{})
			if err != nil {
				return
			}
			err = command.runCommand()
			return
		},
	}

	connectCmd.Flags().StringVarP(&connectOpts.Target, "target", "", "", "target instanceId to start session")
	connectCmd.Flags().BoolVarP(&connectOpts.Interactive, "interactive", "i", false, "Select a value interactively (requires config file)")

	return connectCmd
}

type iConnectCommand interface {
	initConnectCommand(config.ConnectOption, aws.Ec2Client, config.IConfig) error
	runCommand() error
}

type connectCommand struct {
	option                   config.ConnectOption
	ec2Client                aws.Ec2Client
	interactiveFilterCommand string
	command                  ext.ICommand
	config                   config.IConfig
}

func (c *connectCommand) initConnectCommand(o config.ConnectOption, client aws.Ec2Client, conf config.IConfig) (err error) {
	c.ec2Client = client
	c.option = o

	if o.Interactive {
		if err = conf.InitConfig(); err != nil {
			return err
		}

		c.interactiveFilterCommand = conf.GetConfig().InteractiveFilterCommand
	}

	c.command = ext.Command{
		Args: []string{"aws", "ssm", "start-session", "--target"},
	}

	return
}

func (c connectCommand) runCommand() (err error) {
	if err = c.option.IsValid(); err != nil {
		return err
	}

	var target string
	if c.option.Interactive {
		instances, err := c.ec2Client.GetInstances(config.ListOption{})
		if err != nil {
			return err
		}

		filter := ext.Command{
			Args: []string{c.interactiveFilterCommand},
		}

		target, err = instances.GetFilteredInstances(filter)
		if err != nil {
			return err
		}
	} else {
		target = c.option.Target
	}

	startSession.Args = append(startSession.Args, target)
	return startSession.CreateCommand(os.Stdin, os.Stdout, os.Stderr).Run()
}

var startSession = ext.Command{
	Args: []string{"aws", "ssm", "start-session", "--target"},
}
