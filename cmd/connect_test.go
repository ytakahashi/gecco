package cmd

import (
	"testing"

	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
)

type mockedEc2_3 struct{}

func (e mockedEc2_3) GetInstances(options config.ListOption) (instances aws.Ec2Instances, err error) {
	i := aws.Ec2Instance{}
	return aws.Ec2Instances{i}, nil
}

type mockedConnectCommand struct {
}

func (c *mockedConnectCommand) initConnectCommand(o config.ConnectOption, client aws.Ec2Client, conf config.IConfig) (err error) {
	return
}

func (c mockedConnectCommand) runCommand() (err error) {
	return
}

func TestNewConnectCmd(t *testing.T) {
	command := newConnectCmd(&mockedConnectCommand{})

	validate := func(name string, actual string, expected string) {
		if actual != expected {
			t.Errorf("Result of %v was '%v', ecpected: '%v'", name, actual, expected)
		}
	}

	name := "Use"
	expectedUse := "connect"
	actualUse := command.Use
	validate(name, actualUse, expectedUse)

	name = "Short"
	expectedShort := "connect to EC2 instance"
	actualShort := command.Short
	validate(name, actualShort, expectedShort)

	name = "Long"
	expectedLong := "connect to EC2 instance using 'aws cli start-session' command"
	actualLong := command.Long
	validate(name, actualLong, expectedLong)

	actualFlags := command.Flags()
	targetFlag := actualFlags.Lookup("target")

	name = "Flags.target.Name"
	expectedStatusFlagName := "target"
	actualStatusFlagName := targetFlag.Name
	validate(name, actualStatusFlagName, expectedStatusFlagName)

	name = "Flags.target.Usage"
	expectedStatusFlagUsage := "target instanceId to start session"
	actualStatusFlagUsage := targetFlag.Usage
	validate(name, actualStatusFlagUsage, expectedStatusFlagUsage)

	err := command.RunE(nil, nil)
	if err != nil {
		t.Errorf("Error")
	}
}
