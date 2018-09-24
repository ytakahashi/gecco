package aws

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type mockedEc2Service1 struct{}

func (s mockedEc2Service1) initEc2Service() *ec2.EC2 {
	return &ec2.EC2{}
}

func (s mockedEc2Service1) start(ec2Svc *ec2.EC2, dryRun bool, instanceID string) (*ec2.StartInstancesOutput, error) {
	if dryRun {
		return nil, errors.New("error")
	}
	result := &ec2.StartInstancesOutput{
		StartingInstances: []*ec2.InstanceStateChange{},
	}

	return result, nil
}

func (s mockedEc2Service1) stop(ec2Svc *ec2.EC2, dryRun bool, instanceID string) (*ec2.StopInstancesOutput, error) {
	if dryRun {
		return nil, errors.New("error")
	}

	result := &ec2.StopInstancesOutput{
		StoppingInstances: []*ec2.InstanceStateChange{},
	}

	return result, nil
}

func (s mockedEc2Service1) handleError(err error) bool {
	return true
}

type mockedEc2Service2 struct{}

func (s mockedEc2Service2) initEc2Service() *ec2.EC2 {
	return &ec2.EC2{}
}

func (s mockedEc2Service2) start(ec2Svc *ec2.EC2, dryRun bool, instanceID string) (*ec2.StartInstancesOutput, error) {
	if dryRun {
		return nil, errors.New("error")
	}
	result := &ec2.StartInstancesOutput{
		StartingInstances: []*ec2.InstanceStateChange{},
	}

	return result, nil
}

func (s mockedEc2Service2) stop(ec2Svc *ec2.EC2, dryRun bool, instanceID string) (*ec2.StopInstancesOutput, error) {
	if dryRun {
		return nil, errors.New("error")
	}

	result := &ec2.StopInstancesOutput{
		StoppingInstances: []*ec2.InstanceStateChange{},
	}

	return result, nil
}

func (s mockedEc2Service2) handleError(err error) bool {
	return false
}

type mockedEc2Service3 struct{}

func (s mockedEc2Service3) initEc2Service() *ec2.EC2 {
	return &ec2.EC2{}
}

func (s mockedEc2Service3) start(ec2Svc *ec2.EC2, dryRun bool, instanceID string) (*ec2.StartInstancesOutput, error) {
	return nil, errors.New("error")

}

func (s mockedEc2Service3) stop(ec2Svc *ec2.EC2, dryRun bool, instanceID string) (*ec2.StopInstancesOutput, error) {
	return nil, errors.New("error")
}

func (s mockedEc2Service3) handleError(err error) bool {
	return true
}

func Test_Ec2Client_StartInstance_Normal(t *testing.T) {
	sut := Ec2{}
	target := ""
	svc := mockedEc2Service1{}

	actual := sut.StartInstance(target, svc)
	if actual != nil {
		t.Errorf("err: %v", actual)
	}
}

func Test_Ec2Client_StartInstance_Error1(t *testing.T) {
	sut := Ec2{}
	target := ""
	svc := mockedEc2Service2{}

	actual := sut.StartInstance(target, svc)
	if actual == nil {
		t.Errorf("err: %v", actual)
	}
}

func Test_Ec2Client_StartInstance_Error2(t *testing.T) {
	sut := Ec2{}
	target := ""
	svc := mockedEc2Service3{}

	actual := sut.StartInstance(target, svc)
	if actual == nil {
		t.Errorf("err: %v", actual)
	}
}

func Test_Ec2Client_StopInstance_Normal(t *testing.T) {
	sut := Ec2{}
	target := ""
	svc := mockedEc2Service1{}

	actual := sut.StopInstance(target, svc)
	if actual != nil {
		t.Errorf("err: %v", actual)
	}
}

func Test_Ec2Client_StopInstance_Error1(t *testing.T) {
	sut := Ec2{}
	target := ""
	svc := mockedEc2Service2{}

	actual := sut.StopInstance(target, svc)
	if actual == nil {
		t.Errorf("err: %v", actual)
	}
}

func Test_Ec2Client_StopInstance_Error2(t *testing.T) {
	sut := Ec2{}
	target := ""
	svc := mockedEc2Service3{}

	actual := sut.StopInstance(target, svc)
	if actual == nil {
		t.Errorf("err: %v", actual)
	}
}

func Test_Ec2Instances_Print1(t *testing.T) {
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

func Test_Ec2Instances_Print2(t *testing.T) {
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
