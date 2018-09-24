package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
	"github.com/ytakahashi/gecco/ext"
)

var connectOpts = &config.TargetOption{}

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
			return command.runCommand()
		},
	}

	connectCmd.Flags().StringVarP(&connectOpts.Target, "target", "", "", "target instanceId to start session")
	connectCmd.Flags().BoolVarP(&connectOpts.Interactive, "interactive", "i", false, "Select a value interactively (requires config file)")

	return connectCmd
}

type iConnectCommand interface {
	initConnectCommand(config.TargetOption, aws.Ec2Client, config.IConfig) error
	runCommand() error
}

type connectCommand struct {
	option                   config.TargetOption
	ec2Client                aws.Ec2Client
	interactiveFilterCommand string
	command                  ext.ICommand
	config                   config.IConfig
}

func (c *connectCommand) initConnectCommand(o config.TargetOption, client aws.Ec2Client, conf config.IConfig) (err error) {
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

func (c connectCommand) runCommand() error {
	if err := c.option.IsValid(); err != nil {
		return err
	}

	target, err := getTarget(c.option, c.ec2Client, c.interactiveFilterCommand)
	if err != nil {
		return err
	}

	startSession.Args = append(startSession.Args, target)
	return startSession.CreateCommand(os.Stdin, os.Stdout, os.Stderr).Run()
}

var startSession = ext.Command{
	Args: []string{"aws", "ssm", "start-session", "--target"},
}

func getTarget(opts config.TargetOption, client aws.Ec2Client, filterCommand string) (string, error) {
	if opts.Interactive {
		instances, err := client.GetInstances(config.FilterOption{}, aws.Ec2Service{})
		if err != nil {
			return "", err
		}

		filter := ext.Command{
			Args: []string{filterCommand},
		}

		return instances.GetFilteredInstances(filter)
	}

	return opts.Target, nil
}
