package cmd

import (
	"errors"
	"testing"

	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
)

type mockedEc2ForStateTest1 struct{}

func (e mockedEc2ForStateTest1) GetInstances(options config.FilterOption, s aws.IAwsService) (instances aws.Ec2Instances, err error) {
	return nil, nil
}

func (e mockedEc2ForStateTest1) StartInstance(target string, s aws.IAwsService) error {
	return nil
}

func (e mockedEc2ForStateTest1) StopInstance(target string, s aws.IAwsService) error {
	return nil
}

func TestGetFunc1(t *testing.T) {
	operation := startInstance

	_, err := operation.getFunc(mockedEc2ForStateTest1{})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetFunc2(t *testing.T) {
	operation := stopInstance

	_, err := operation.getFunc(mockedEc2ForStateTest1{})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetFunc3(t *testing.T) {
	operation := unknown

	_, err := operation.getFunc(mockedEc2ForStateTest1{})
	if err == nil {
		t.Error("Error")
	}
}

type mockedStateCommand1 struct {
}

func (c *mockedStateCommand1) initStateCommand(o config.TargetOption, client aws.Ec2Client, conf config.IConfig) (err error) {
	return
}

func (c mockedStateCommand1) runCommand(o instanceOperation) (err error) {
	return
}

type mockedStateCommand2 struct {
}

func (c *mockedStateCommand2) initStateCommand(o config.TargetOption, client aws.Ec2Client, conf config.IConfig) error {
	return errors.New("error")
}

func (c mockedStateCommand2) runCommand(o instanceOperation) (err error) {
	return
}

func TestNewStartCmd(t *testing.T) {
	command := newStartCmd(&mockedStateCommand1{})

	validate := func(name string, actual string, expected string) {
		if actual != expected {
			t.Errorf("Result of %v was '%v', ecpected: '%v'", name, actual, expected)
		}
	}

	name := "Use"
	expectedUse := "start"
	actualUse := command.Use
	validate(name, actualUse, expectedUse)

	name = "Short"
	expectedShort := "start specified EC2 instance"
	actualShort := command.Short
	validate(name, actualShort, expectedShort)

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
		t.Errorf("Error: %v", err)
	}
}

func TestNewStartCmd_Error(t *testing.T) {
	command := newStartCmd(&mockedStateCommand2{})

	err := command.RunE(nil, nil)
	if err == nil {
		t.Error("Error")
	}
}

func TestNewStopCmd(t *testing.T) {
	command := newStopCmd(&mockedStateCommand1{})

	validate := func(name string, actual string, expected string) {
		if actual != expected {
			t.Errorf("Result of %v was '%v', ecpected: '%v'", name, actual, expected)
		}
	}

	name := "Use"
	expectedUse := "stop"
	actualUse := command.Use
	validate(name, actualUse, expectedUse)

	name = "Short"
	expectedShort := "stop specified EC2 instance"
	actualShort := command.Short
	validate(name, actualShort, expectedShort)

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
		t.Errorf("Error: %v", err)
	}
}

func TestNewStopCmd_Error(t *testing.T) {
	command := newStopCmd(&mockedStateCommand2{})

	err := command.RunE(nil, nil)
	if err == nil {
		t.Error("Error")
	}
}

func TestInitStateCommand_Normal1(t *testing.T) {
	opts := config.TargetOption{
		Interactive: false,
		Target:      "foo",
	}

	command := stateCommand{}

	err := command.initStateCommand(opts, mockedEc2ForStateTest1{}, &config.Config{})

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
}

func TestInitStateCommand_Normal2(t *testing.T) {
	opts := config.TargetOption{
		Interactive: true,
	}

	command := stateCommand{}

	err := command.initStateCommand(opts, mockedEc2ForStateTest1{}, &mockedConfig1{})

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
}

func TestInitStateCommand_Error(t *testing.T) {
	o := config.TargetOption{
		Interactive: true,
	}

	command := &stateCommand{
		config: &mockedConfig1{},
	}

	actual := command.initStateCommand(o, aws.Ec2{}, mockedConfig2{})

	if actual == nil {
		t.Error("error")
	}
}

type mockedEc2ForStateTest2 struct{}

func (e mockedEc2ForStateTest2) GetInstances(options config.FilterOption, s aws.IAwsService) (instances aws.Ec2Instances, err error) {
	return nil, errors.New("error")
}

func (e mockedEc2ForStateTest2) StartInstance(target string, s aws.IAwsService) error {
	return nil
}

func (e mockedEc2ForStateTest2) StopInstance(target string, s aws.IAwsService) error {
	return nil
}

func Test_StateCommand_RunCommand_Normal(t *testing.T) {
	opts := config.TargetOption{
		Target: "target",
	}
	command := stateCommand{
		option:                   opts,
		ec2Client:                mockedEc2ForStateTest1{},
		interactiveFilterCommand: "echo",
	}

	err := command.runCommand(startInstance)

	if err != nil {
		t.Errorf("%v", err)
	}
}

func Test_StateCommand_RunCommand_Error1(t *testing.T) {
	opts := config.TargetOption{}

	command := stateCommand{
		option:                   opts,
		ec2Client:                mockedEc2ForStateTest2{},
		interactiveFilterCommand: "echo",
	}

	err := command.runCommand(startInstance)

	if err == nil {
		t.Errorf("error")
	}
}

func Test_StateCommand_RunCommand_Error2(t *testing.T) {
	opts := config.TargetOption{
		Interactive: true,
	}

	command := stateCommand{
		option:                   opts,
		ec2Client:                mockedEc2ForStateTest2{},
		interactiveFilterCommand: "echo",
	}

	err := command.runCommand(startInstance)

	if err == nil {
		t.Errorf("error")
	}
}

func Test_StateCommand_RunCommand_Error3(t *testing.T) {
	opts := config.TargetOption{
		Target: "target",
	}

	command := stateCommand{
		option:                   opts,
		ec2Client:                mockedEc2ForStateTest2{},
		interactiveFilterCommand: "echo",
	}

	err := command.runCommand(unknown)

	if err == nil {
		t.Errorf("error")
	}
}
