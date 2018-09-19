package cmd

import (
	"errors"
	"testing"

	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
	"github.com/ytakahashi/gecco/ext"
)

type mockedConnectCommand1 struct {
}

func (c *mockedConnectCommand1) initConnectCommand(o config.ConnectOption, client aws.Ec2Client, conf config.IConfig) (err error) {
	return
}

func (c mockedConnectCommand1) runCommand() (err error) {
	return
}

func TestNewConnectCmd(t *testing.T) {
	command := newConnectCmd(&mockedConnectCommand1{})

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

func TestRunCommand_Error1(t *testing.T) {
	opts := config.ConnectOption{
		Target:      "foo",
		Interactive: true,
	}
	command := connectCommand{
		option: opts,
	}

	err := command.runCommand()

	if err == nil {
		t.Error("Erroe")
	}
}

func TestInitConnectCommand(t *testing.T) {
	opts := config.ConnectOption{
		Interactive: false,
	}

	command := connectCommand{}

	err := command.initConnectCommand(opts, mockedEc2_3{}, config.Config{})

	if err != nil {
		t.Errorf("%v", err)
	}
}

type mockedEc2_3 struct{}

func (e mockedEc2_3) GetInstances(options config.ListOption) (instances aws.Ec2Instances, err error) {
	return nil, errors.New("error")
}

func TestRunCommand_Error2(t *testing.T) {
	opts := config.ConnectOption{
		Interactive: true,
	}
	command := connectCommand{
		option:    opts,
		ec2Client: mockedEc2_3{},
	}

	err := command.runCommand()

	if err == nil {
		t.Error("Error")
	}
}

type mockedEc2_4 struct{}

func (e mockedEc2_4) GetInstances(options config.ListOption) (instances aws.Ec2Instances, err error) {
	return aws.Ec2Instances{}, nil
}

func TestRunCommand_Normal1(t *testing.T) {
	opts := config.ConnectOption{
		Interactive: true,
	}
	command := connectCommand{
		option:                   opts,
		ec2Client:                mockedEc2_4{},
		interactiveFilterCommand: "echo",
	}

	startSession = ext.Command{
		Args: []string{"echo"},
	}

	err := command.runCommand()

	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestRunCommand_Normal2(t *testing.T) {
	opts := config.ConnectOption{
		Target:      "foo",
		Interactive: false,
	}
	command := connectCommand{
		option: opts,
	}

	startSession = ext.Command{
		Args: []string{"echo"},
	}

	err := command.runCommand()

	if err != nil {
		t.Errorf("%v", err)
	}
}
