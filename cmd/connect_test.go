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

func (c *mockedConnectCommand1) initConnectCommand(o config.TargetOption, client aws.Ec2Client, conf config.IConfig) (err error) {
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

func Test_ConnectCommand_InitConnectCommand_Normal1(t *testing.T) {
	opts := config.TargetOption{
		Interactive: false,
		Target:      "foo",
	}

	command := connectCommand{}

	err := command.initConnectCommand(opts, mockedEc2_3{}, &config.Config{})

	if err != nil {
		t.Errorf("%v", err)
	}

	if command.option.Interactive == true {
		t.Error("err")
	}

	if command.option.Target != "foo" {
		t.Errorf("error: %v", command.option.Target)
	}

	if command.interactiveFilterCommand != "" {
		t.Errorf("error: %v", command.interactiveFilterCommand)
	}

	if command.command == nil {
		t.Error("err")
	}
}

type mockedConfig1 struct {
}

func (c *mockedConfig1) InitConfig() (err error) {
	return nil
}

func (c mockedConfig1) GetConfig() config.Config {
	return config.Config{
		InteractiveFilterCommand: "foo",
	}
}

func Test_ConnectCommand_InitConnectCommand_Normal2(t *testing.T) {
	opts := config.TargetOption{
		Interactive: true,
	}

	command := connectCommand{}

	err := command.initConnectCommand(opts, mockedEc2_3{}, &mockedConfig1{})

	if err != nil {
		t.Errorf("%v", err)
	}

	if command.option.Interactive == false {
		t.Error("err")
	}

	if command.option.Target != "" {
		t.Errorf("error: %v", command.option.Target)
	}

	if command.interactiveFilterCommand != "foo" {
		t.Errorf("error: %v", command.interactiveFilterCommand)
	}

	if command.command == nil {
		t.Error("err")
	}
}

type mockedConfig2 struct {
}

func (c mockedConfig2) InitConfig() (err error) {
	return errors.New("error")
}

func (c mockedConfig2) GetConfig() config.Config {
	return config.Config{
		InteractiveFilterCommand: "foo",
	}
}

func Test_ConnectCommand_InitConnectCommand_Error(t *testing.T) {
	o := config.TargetOption{
		Interactive: true,
	}

	command := &connectCommand{
		config: &mockedConfig1{},
	}

	actual := command.initConnectCommand(o, aws.Ec2{}, mockedConfig2{})

	if actual == nil {
		t.Error("error")
	}
}

type mockedEc2_3 struct{}

func (e mockedEc2_3) GetInstances(options config.FilterOption, s aws.IEc2Service) (instances aws.Ec2Instances, err error) {
	return nil, errors.New("error")
}

func (e mockedEc2_3) StartInstance(target string, s aws.IEc2Service) error {
	return nil
}

func (e mockedEc2_3) StopInstance(target string, s aws.IEc2Service) error {
	return nil
}

func Test_ConnectCommand_RunCommand_Error1(t *testing.T) {
	opts := config.TargetOption{
		Target:      "foo",
		Interactive: true,
	}
	command := connectCommand{
		option: opts,
	}

	err := command.runCommand()

	if err == nil {
		t.Error("Error")
	}
}

func Test_ConnectCommand_RunCommand_Error2(t *testing.T) {
	opts := config.TargetOption{
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

func (e mockedEc2_4) GetInstances(options config.FilterOption, s aws.IEc2Service) (instances aws.Ec2Instances, err error) {
	return aws.Ec2Instances{}, nil
}

func (e mockedEc2_4) StartInstance(target string, s aws.IEc2Service) error {
	return nil
}

func (e mockedEc2_4) StopInstance(target string, s aws.IEc2Service) error {
	return nil
}

func Test_ConnectCommand_RunCommand_Normal1(t *testing.T) {
	opts := config.TargetOption{
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

func Test_ConnectCommand_RunCommand_Normal2(t *testing.T) {
	opts := config.TargetOption{
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
