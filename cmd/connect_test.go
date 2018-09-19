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

func TestConnect1(t *testing.T) {
	c := config.ConnectOptions{}

	fn := func(target string) error {
		return nil
	}

	fn2 := func() error {
		return nil
	}

	e := mockedEc2_3{}

	expected := "Option '--target' is not specified"
	actual := connect(c, fn, fn2, e)

	if actual == nil {
		t.Errorf("Error should be thrown.")
	}

	if actual.Error() != expected {
		t.Errorf("Error:\n Actual: %v\n Expected: %v", actual, expected)
	}
}

func TestConnect2(t *testing.T) {
	c := config.ConnectOptions{
		Target: "target",
	}

	fn := func(target string) error {
		return nil
	}

	fn2 := func() error {
		return nil
	}

	e := mockedEc2_3{}

	actual := connect(c, fn, fn2, e)

	if actual != nil {
		t.Errorf("Error should not be thrown.")
	}
}

func TestNewConnectCmd(t *testing.T) {
	command := newConnectCmd()

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
}
