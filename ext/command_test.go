package ext

import (
	"os"
	"testing"
)

func TestCreateCommand(t *testing.T) {

	var cmd = Command{
		[]string{"foo", "bar"},
	}

	// target := "TARGET"
	actual := cmd.CreateCommand(os.Stdin, os.Stdout, os.Stderr)

	validate := func(act, expc string) {
		if act != expc {
			t.Errorf("Error:\n Actual: %v\nExpected: %v", act, expc)
		}
	}

	if len(actual.Args) != 2 {
		t.Errorf("Error")
	}

	validate(actual.Args[0], "foo")
	validate(actual.Args[1], "bar")

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
