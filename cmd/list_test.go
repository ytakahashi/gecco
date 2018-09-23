package cmd

import (
	"errors"
	"testing"

	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
)

type mockedListCommand1 struct{}

func (c mockedListCommand1) runCommand() error {
	return nil
}

func (c mockedListCommand1) initListCommand(o config.ListOption, client aws.Ec2Client) {
}

func TestNewListCmd(t *testing.T) {
	command := newListCmd(mockedListCommand1{})

	validate := func(name string, actual string, expected string) {
		if actual != expected {
			t.Errorf("Result of %v was '%v', ecpected: '%v'", name, actual, expected)
		}
	}

	name := "Use"
	expectedUse := "list"
	actualUse := command.Use
	validate(name, actualUse, expectedUse)

	name = "Short"
	expectedShort := "lists EC2 instances"
	actualShort := command.Short
	validate(name, actualShort, expectedShort)

	actualFlags := command.Flags()
	statusFlag := actualFlags.Lookup("status")
	tagKeyFlag := actualFlags.Lookup("tagKey")
	tagValueFlag := actualFlags.Lookup("tagValue")

	name = "Flags.status.Name"
	expectedStatusFlagName := "status"
	actualStatusFlagName := statusFlag.Name
	validate(name, actualStatusFlagName, expectedStatusFlagName)

	name = "Flags.status.Usage"
	expectedStatusFlagUsage := "filters instances by status"
	actualStatusFlagUsage := statusFlag.Usage
	validate(name, actualStatusFlagUsage, expectedStatusFlagUsage)

	name = "Flags.tagKey.Name"
	expectedTagKeyFlagName := "tagKey"
	actualTagKeyFlagName := tagKeyFlag.Name
	validate(name, expectedTagKeyFlagName, actualTagKeyFlagName)

	name = "Flags.tagKey.Usage"
	expectedTagKeyFlagUsage := "filters instances by tag key"
	actualTagKeyFlagUsage := tagKeyFlag.Usage
	validate(name, expectedTagKeyFlagUsage, actualTagKeyFlagUsage)

	name = "Flags.tagValue.Name"
	expectedTagValueFlagName := "tagValue"
	actualTagValueFlagName := tagValueFlag.Name
	validate(name, expectedTagValueFlagName, actualTagValueFlagName)

	name = "Flags.tagValue.Usage"
	expectedTagValueFlagUsage := "filters instances by tag value"
	actualTagValueFlagUsage := tagValueFlag.Usage
	validate(name, expectedTagValueFlagUsage, actualTagValueFlagUsage)

	err := command.RunE(nil, nil)
	if err != nil {
		t.Errorf("Error")
	}
}

type mockedEc2_1 struct{}

func (e mockedEc2_1) GetInstances(o config.ListOption, s aws.IAwsService) (instances aws.Ec2Instances, err error) {
	return aws.Ec2Instances{}, nil
}

func (e mockedEc2_1) StartInstance(target string, s aws.IAwsService) error {
	return nil
}

func (e mockedEc2_1) StopInstance(target string, s aws.IAwsService) error {
	return nil
}

type mockedEc2_2 struct{}

func (e mockedEc2_2) GetInstances(o config.ListOption, s aws.IAwsService) (instances aws.Ec2Instances, err error) {
	return nil, errors.New("error")
}

func (e mockedEc2_2) StartInstance(target string, s aws.IAwsService) error {
	return nil
}

func (e mockedEc2_2) StopInstance(target string, s aws.IAwsService) error {
	return nil
}

func TestInitListCommand(t *testing.T) {
	o := config.ListOption{
		Status: "status",
	}
	sut := listCommand{}
	sut.initListCommand(o, mockedEc2_1{})

	if sut.options != o {
		t.Errorf("Error %v", sut.options)
	}

	if sut.ec2Client == nil {
		t.Error("Error")
	}
}

func TestRunCommand1(t *testing.T) {
	sut := listCommand{
		options: config.ListOption{Status: "foo"},
	}

	err := sut.runCommand()

	if err == nil {
		t.Error("Error should be thrown")
	}
}

func TestRunCommand2(t *testing.T) {
	sut := listCommand{
		options:   config.ListOption{},
		ec2Client: mockedEc2_2{},
	}

	err := sut.runCommand()

	if err == nil {
		t.Error("Error should be thrown")
	}
}

func TestRunCommand3(t *testing.T) {
	sut := listCommand{
		options:   config.ListOption{},
		ec2Client: mockedEc2_1{},
	}

	err := sut.runCommand()

	if err != nil {
		t.Errorf("%v", err)
	}
}
