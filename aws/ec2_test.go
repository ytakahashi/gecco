package aws

import (
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ytakahashi/gecco/config"
)

func Test_Tag_ToString(t *testing.T) {
	tag := Tag{
		Key:   "k",
		Value: "v",
	}

	actual := tag.toString()
	expected := "{\"k\": \"v\"}"
	if actual != expected {
		t.Errorf("\nExpected: '%s'\n Actual '%s'", expected, actual)
	}
}

func Test_Tags_ToString1(t *testing.T) {
	tags := Tags{}
	actual := tags.ToString()
	expected := ""
	if actual != expected {
		t.Errorf("\nExpected: '%s'\n Actual '%s'", expected, actual)
	}
}

func Test_Tags_ToString2(t *testing.T) {
	tags := Tags{{Key: "k1", Value: "v1"}, {Key: "k2", Value: "v2"}}
	actual := tags.ToString()
	expected := "[{\"k1\": \"v1\"}, {\"k2\": \"v2\"}]"
	if actual != expected {
		t.Errorf("\nExpected: '%s'\n Actual '%s'", expected, actual)
	}
}

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
			InstanceID: "instanceId",
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

func Test_Ec2Instances_ToString1(t *testing.T) {
	i := Ec2Instance{
		InstanceID:   "instance1",
		InstanceType: "type1",
		Status:       "status1",
		Tags:         []Tag{{Key: "k", Value: "v"}},
	}
	instanceList := Ec2Instances{i}

	expected := "instance1 type1 status1 [{\"k\": \"v\"}]\n"

	actual, err := instanceList.ToString(config.Text)
	if actual != expected {
		t.Errorf("\nExpected: '%s'\n Actual '%s'", expected, actual)
	}

	if err != nil {
		t.Errorf("ERR: %v", err)
	}
}

func Test_Ec2Instances_ToString2(t *testing.T) {
	i := Ec2Instance{
		InstanceID:   "instance1",
		InstanceType: "type1",
		Status:       "status1",
		Tags:         []Tag{{Key: "k", Value: "v"}},
	}
	instanceList := Ec2Instances{i}

	b, _ := json.MarshalIndent(&instanceList, "", "    ")
	expected := string(b)

	actual, err := instanceList.ToString(config.JSON)

	if actual != expected {
		t.Errorf("\nExpected: '%s'\n Actual '%s'", expected, actual)
	}

	if err != nil {
		t.Errorf("ERR: %v", err)
	}
}

func Test_Ec2Instances_ToString3(t *testing.T) {
	i := Ec2Instance{
		InstanceID:   "instance1",
		InstanceType: "type1",
		Status:       "status1",
		Tags:         []Tag{{Key: "k", Value: "v"}},
	}
	instanceList := Ec2Instances{i}

	_, err := instanceList.ToString(config.Unknown)

	if err == nil {
		t.Errorf("ERR: %v", err)
	}
}
