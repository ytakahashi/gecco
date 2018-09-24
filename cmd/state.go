package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
)

type instanceOperation int

const (
	unknown instanceOperation = iota
	startInstance
	stopInstance
)

func (c instanceOperation) getFunc(e aws.Ec2Client) (func(string, aws.IAwsService) error, error) {
	switch c {
	case startInstance:
		return e.StartInstance, nil
	case stopInstance:
		return e.StopInstance, nil
	default:
		return nil, errors.New("undefined operation")
	}
}

func newStartCmd(command iStateCommand) *cobra.Command {
	var startOpts = &config.TargetOption{}
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start specified EC2 instance",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			err = command.initStateCommand(*startOpts, aws.Ec2{}, &config.Config{})
			if err != nil {
				return
			}
			return command.runCommand(startInstance)
		},
	}

	startCmd.Flags().StringVarP(&startOpts.Target, "target", "", "", "target instanceId to start session")
	startCmd.Flags().BoolVarP(&startOpts.Interactive, "interactive", "i", false, "Select a value interactively (requires config file)")

	return startCmd
}

func newStopCmd(command iStateCommand) *cobra.Command {
	var stopOpts = &config.TargetOption{}
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "stop specified EC2 instance",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			err = command.initStateCommand(*stopOpts, aws.Ec2{}, &config.Config{})
			if err != nil {
				return
			}
			return command.runCommand(stopInstance)
		},
	}

	stopCmd.Flags().StringVarP(&stopOpts.Target, "target", "", "", "target instanceId to start session")
	stopCmd.Flags().BoolVarP(&stopOpts.Interactive, "interactive", "i", false, "Select a value interactively (requires config file)")

	return stopCmd
}

type iStateCommand interface {
	initStateCommand(o config.TargetOption, client aws.Ec2Client, conf config.IConfig) error
	runCommand(instanceOperation) error
}

type stateCommand struct {
	option                   config.TargetOption
	ec2Client                aws.Ec2Client
	interactiveFilterCommand string
	config                   config.IConfig
}

func (c *stateCommand) initStateCommand(o config.TargetOption, client aws.Ec2Client, conf config.IConfig) (err error) {
	c.ec2Client = client
	c.option = o

	if o.Interactive {
		if err = conf.InitConfig(); err != nil {
			return err
		}

		c.interactiveFilterCommand = conf.GetConfig().InteractiveFilterCommand
	}

	return
}

func (c stateCommand) runCommand(op instanceOperation) error {
	if err := c.option.IsValid(); err != nil {
		return err
	}

	target, err := getTarget(c.option, c.ec2Client, c.interactiveFilterCommand)
	if err != nil {
		return err
	}

	fn, err := op.getFunc(c.ec2Client)
	if err != nil {
		return err
	}

	return fn(target, aws.Ec2Service{})
}
