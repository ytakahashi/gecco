package cmd

import (
	"testing"

	"github.com/ytakahashi/gecco/config"
)

func TestConnect1(t *testing.T) {
	c := config.ConnectOptions{}

	fn := func(target string) error {
		return nil
	}

	fn2 := func() error {
		return nil
	}

	expected := "Option '--target' is not specified"
	actual := connect(c, fn, fn2)

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

	actual := connect(c, fn, fn2)

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

func TestCreateCommand(t *testing.T) {
	target := "TARGET"
	actual := createCommand(target)

	validate := func(act, expc string) {
		if act != expc {
			t.Errorf("Error:\n Actual: %v\nExpected: %v", act, expc)
		}
	}

	if len(actual.Args) != 5 {
		t.Errorf("Error")
	}

	validate(actual.Args[0], "aws")
	validate(actual.Args[1], "ssm")
	validate(actual.Args[2], "start-session")
	validate(actual.Args[3], "--target")
	validate(actual.Args[4], "TARGET")

	if actual.Stdin == nil {
		t.Errorf("Error: should not be nil.")
	}

	if actual.Stdout == nil {
		t.Errorf("Error: should not be nil.")
	}

	if actual.Stderr == nil {
		t.Errorf("Error: should not be nil.")
	}
}
