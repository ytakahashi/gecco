package aws

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ytakahashi/gecco/config"
)

func TestPrint1(t *testing.T) {
	buf := &bytes.Buffer{}
	i := Ec2Instance{
		instanceID:   "instance id",
		instanceType: "instance type",
		status:       "status",
	}
	instances := Ec2Instances{i}

	instances.Print(buf)

	expected := fmt.Sprintln(
		i.instanceID,
		i.instanceType,
		i.status,
		"",
	)
	actual := buf.String()

	if expected != actual {
		t.Errorf("\nExpected: '%s'\n Actual '%s'", expected, actual)
	}
}

func TestPrint2(t *testing.T) {
	buf := &bytes.Buffer{}
	i := Ec2Instance{
		instanceID:   "instance id",
		instanceType: "instance type",
		status:       "status",
		tags:         []tag{{key: "k", value: "v"}},
	}
	instances := Ec2Instances{i}

	instances.Print(buf)

	expected := fmt.Sprintln(
		i.instanceID,
		i.instanceType,
		i.status,
		"[{\"k\": \"v\"}]",
	)
	actual := buf.String()

	if expected != actual {
		t.Errorf("\nExpected: '%s'\n Actual '%s'", expected, actual)
	}
}

func TestCreateInput1(t *testing.T) {
	options := config.FilterOption{}

	input := createInput(options)

	filters := input.Filters

	expectedLength := 1
	actualLength := len(filters)

	if actualLength != expectedLength {
		t.Errorf("Number of filters should not be '%v'", actualLength)
	}

	expectedFilterName := "instance-state-name"
	actualFilterName := filters[0].Name

	actualFilterValues := filters[0].Values
	actualFilterlLength := len(actualFilterValues)

	validateFilter(expectedFilterName, aws.StringValue(actualFilterName), actualFilterlLength, 6, t)
}

func TestCreateInput2(t *testing.T) {
	options := config.FilterOption{
		Status:   "foo",
		TagKey:   "k",
		TagValue: "v",
	}

	input := createInput(options)

	filters := input.Filters

	expectedLength := 2
	actualLength := len(filters)

	if actualLength != expectedLength {
		t.Errorf("Number of filters should not be '%v'", actualLength)
	}

	expectedFilterName0 := "instance-state-name"
	actualFilterName0 := filters[0].Name

	actualFilterValues0 := filters[0].Values
	actuaFilterlLength0 := len(actualFilterValues0)

	validateFilter(expectedFilterName0, aws.StringValue(actualFilterName0), actuaFilterlLength0, 1, t)

	expectedFilterName1 := "tag:" + options.TagKey
	actualFilterName1 := filters[1].Name

	actualFilterValues1 := filters[1].Values
	actuaFilterlLength1 := len(actualFilterValues1)

	validateFilter(expectedFilterName1, aws.StringValue(actualFilterName1), actuaFilterlLength1, 1, t)

}

func validateFilter(expectedFilterName, actualFilterName string, expectedFilterlLength, actualFilterlLength int, t *testing.T) {
	if actualFilterName != expectedFilterName {
		t.Errorf("Filter name mismatch:\nActual: %v\nExpected: %v", actualFilterName, expectedFilterName)
	}

	if actualFilterlLength != expectedFilterlLength {
		t.Errorf("Invalid number of filters:\nActual: %v\nExpected: %v", actualFilterlLength, expectedFilterlLength)
	}
}

func TestNew(t *testing.T) {
	i := ec2.Instance{}

	ta := ec2.Tag{}
	ta.SetKey("k").SetValue("v")
	tags := []*ec2.Tag{&ta}

	expectedAvailabilityZone := "zone 1"
	p := ec2.Placement{}
	p.SetAvailabilityZone(expectedAvailabilityZone)

	expectedStatus := "state name"
	s := ec2.InstanceState{}
	s.SetName(expectedStatus)

	expectedInstanceID := "id"
	i.SetInstanceId(expectedInstanceID).SetTags(tags).SetPlacement(&p).SetState(&s)

	actual := newEc2Instance(i)

	if actual.instanceID != expectedInstanceID {
		t.Errorf("InstanceID:\nActual: %v\nExpected: %v", actual.instanceID, expectedInstanceID)
	}

	if actual.instanceType != "" {
		t.Errorf("instanceType:\nActual: %v\nExpected: %v", actual.instanceType, "expectedFilterlLength")
	}

	if actual.status != expectedStatus {
		t.Errorf("status:\nActual: %v\nExpected: %v", actual.status, expectedStatus)
	}

	if len(actual.tags) > 1 {
		t.Errorf("Error")
	}

	key := "k"
	if actual.tags[0].key != key {
		t.Errorf("tag key:\nActual: %v\nExpected: %v", actual.tags[0].key, key)
	}

	val := "v"
	if actual.tags[0].value != val {
		t.Errorf("tag value:\nActual: %v\nExpected: %v", actual.tags[0].value, val)
	}
}

type mockedCommand1 struct{}
type mockedCommand2 struct{}

func (c mockedCommand1) CreateCommand(i io.Reader, o io.Writer, e io.Writer) *exec.Cmd {
	command := exec.Command("echo")
	command.Args = []string{"echo", "foo"}
	command.Stdin = i
	command.Stdout = o
	command.Stderr = e
	return command
}

func (c mockedCommand2) CreateCommand(i io.Reader, o io.Writer, e io.Writer) *exec.Cmd {
	command := exec.Command("foo")
	command.Args = []string{"foo"}
	command.Stdin = i
	command.Stdout = o
	command.Stderr = e
	return command
}

func TestGetFilteredInstances1(t *testing.T) {
	i := Ec2Instances{
		Ec2Instance{
			instanceID: "instanceId",
		},
	}

	mockedFilter := mockedCommand1{}

	str, err := i.GetFilteredInstances(mockedFilter)
	if err != nil {
		t.Errorf("err value:\nActual: %v", err)
	}
	if str == "" {
		t.Errorf("Error")
	}
	if str != "foo" {
		t.Errorf("str value:\nActual: %v", str)
	}

}

func TestGetFilteredInstances2(t *testing.T) {
	i := Ec2Instances{}

	mockedFilter := mockedCommand2{}

	str, err := i.GetFilteredInstances(mockedFilter)
	if err == nil {
		t.Errorf("Error should be thrown.")
	}
	if str != "" {
		t.Errorf("str value:\nActual: %v", str)
	}
}
