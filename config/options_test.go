package config

import (
	"fmt"
	"testing"
)

func TestOutputFormat_FromString_Error1(t *testing.T) {
	assert := func(actual OutputFormat, expected OutputFormat) {
		if actual != expected {
			t.Errorf("Actual '%v', ecpected: '%v'", actual, expected)
		}
	}

	value1 := "text"
	value2 := "json"
	value3 := "foo"

	expected1 := Text
	expected2 := JSON
	expected3 := Unknown

	actual1 := NewOutputFormat(value1)
	actual2 := NewOutputFormat(value2)
	actual3 := NewOutputFormat(value3)

	assert(actual1, expected1)
	assert(actual2, expected2)
	assert(actual3, expected3)
}

func TestTargetOption_IsValid_Error1(t *testing.T) {
	c := TargetOption{}
	expected := "Option '--target' or '-i' is required"
	actual := c.IsValid()

	if actual == nil {
		t.Errorf("Error should be thrown.")
	}

	if actual.Error() != expected {
		t.Errorf("Error:\n Actual: %v\n Expected: %v", actual, expected)
	}
}

func TestTargetOption_IsValid_Error2(t *testing.T) {
	c := TargetOption{
		Target:      "foo",
		Interactive: true,
	}
	expected := "Options '--target' and '-i' cannot be used at the same time"
	actual := c.IsValid()

	if actual == nil {
		t.Errorf("Error should be thrown.")
	}

	if actual.Error() != expected {
		t.Errorf("Error:\n Actual: %v\n Expected: %v", actual, expected)
	}
}

func TestTargetOption_IsValid_Ok(t *testing.T) {
	c := TargetOption{
		Interactive: true,
	}

	actual := c.IsValid()

	if actual != nil {
		t.Errorf("Error should not be thrown.")
	}
}

func Test_FilterOption_IsValid_Normal1(t *testing.T) {
	options := FilterOption{
		OutputFormat: "text",
	}

	actual := options.IsValid()

	if actual != nil {
		t.Errorf("Error")
	}
}

func Test_FilterOption_IsValid_Normal2(t *testing.T) {
	options := FilterOption{
		OutputFormat: "JSON",
	}

	actual := options.IsValid()

	if actual != nil {
		t.Errorf("Error")
	}
}

func Test_FilterOption_IsValid_InvalidStatus(t *testing.T) {
	status := "foo"
	options := FilterOption{Status: status}

	expected := fmt.Sprintf("Invalid status (%v)", status)
	actual := options.IsValid()

	if actual == nil {
		t.Errorf("Error")
	}

	if actual.Error() != expected {
		t.Errorf("Error:\n Actual: %v\n Expected: %v", actual, expected)
	}
}

func Test_FilterOption_IsValid_InvalidTags1(t *testing.T) {
	tagKey := "foo"
	options := FilterOption{TagKey: tagKey}

	expected := "Option '--tagValue' is required when '--tagKey' is specified"
	actual := options.IsValid()

	if actual == nil {
		t.Errorf("Error")
	}

	if actual.Error() != expected {
		t.Errorf("Error:\n Actual: %v\n Expected: %v", actual, expected)
	}
}

func Test_FilterOption_IsValid_InvalidTags2(t *testing.T) {
	tagValue := "foo"
	options := FilterOption{TagValue: tagValue}

	expected := "Option '--tagKey' is required when '--tagValue' is specified"
	actual := options.IsValid()

	if actual == nil {
		t.Errorf("Error")
	}

	if actual.Error() != expected {
		t.Errorf("Error:\n Actual: %v\n Expected: %v", actual, expected)
	}
}

func Test_FilterOption_IsValid_InvalidFormat(t *testing.T) {
	options := FilterOption{
		OutputFormat: "foo",
	}
	expected := "Option '--output' should be one of 'text' or 'json'"
	actual := options.IsValid()

	if actual == nil {
		t.Errorf("Error")
	}

	if actual.Error() != expected {
		t.Errorf("Error:\n Actual: %v\n Expected: %v", actual, expected)
	}
}

func Test_FilterOption_GetOutputFormat(t *testing.T) {
	options := FilterOption{
		OutputFormat: "json",
	}

	expected := JSON
	actual := options.GetOutputFormat()

	if actual != expected {
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
