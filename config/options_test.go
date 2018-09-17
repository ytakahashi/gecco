package config

import (
	"fmt"
	"testing"
)

func TestIsValid_One(t *testing.T) {
	options := ListOption{}

	actual := options.IsValid()

	if actual != nil {
		t.Errorf("Error")
	}
}

func TestIsValid_InvalidStatus(t *testing.T) {
	status := "foo"
	options := ListOption{Status: status}

	expected := fmt.Sprintf("Invalid status (%v)", status)
	actual := options.IsValid()

	if actual == nil {
		t.Errorf("Error")
	}

	if actual.Error() != expected {
		t.Errorf("Error:\n Actual: %v\n Expected: %v", actual, expected)
	}
}

func TestIsValid_InvalidTags1(t *testing.T) {
	tagKey := "foo"
	options := ListOption{TagKey: tagKey}

	expected := "Option '--tagValue' is required when '--tagKey' is specified"
	actual := options.IsValid()

	if actual == nil {
		t.Errorf("Error")
	}

	if actual.Error() != expected {
		t.Errorf("Error:\n Actual: %v\n Expected: %v", actual, expected)
	}
}

func TestIsValid_InvalidTags2(t *testing.T) {
	tagValue := "foo"
	options := ListOption{TagValue: tagValue}

	expected := "Option '--tagKey' is required when '--tagValue' is specified"
	actual := options.IsValid()

	if actual == nil {
		t.Errorf("Error")
	}

	if actual.Error() != expected {
		t.Errorf("Error:\n Actual: %v\n Expected: %v", actual, expected)
	}
}

func TestIsValidStatus1(t *testing.T) {
	status := [6]string{"running", "stopping", "pending", "shutting-down", "terminated", "stopped"}

	for i := 0; i < len(status); i++ {
		actual := isValidStatus(status[i])
		if actual != true {
			t.Errorf("Result of isValidStatus for Status '%v' was %v", status[i], actual)
		}
	}
}
