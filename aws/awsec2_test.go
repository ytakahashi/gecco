package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ytakahashi/gecco/config"
)

func Test_Ec2Service_InitEc2Service1(t *testing.T) {
	sut := Ec2Service{}
	actual := sut.initEc2Service()
	if actual == nil {
		t.Errorf("err: %v", actual)
	}
}

type mockedAwsError1 struct{}

func (e mockedAwsError1) Error() string {
	return ""
}

func (e mockedAwsError1) Code() string {
	return "DryRunOperation"
}

func (e mockedAwsError1) Message() string {
	return ""
}

func (e mockedAwsError1) OrigErr() error {
	return nil
}

func Test_Ec2Service_HandleError1(t *testing.T) {
	sut := Ec2Service{}

	actual := sut.handleError(mockedAwsError1{})
	if actual != true {
		t.Error("Error")
	}
}

type mockedAwsError2 struct{}

func (e mockedAwsError2) Error() string {
	return ""
}

func (e mockedAwsError2) Code() string {
	return ""
}

func (e mockedAwsError2) Message() string {
	return ""
}

func (e mockedAwsError2) OrigErr() error {
	return nil
}

func Test_Ec2Service_HandleError2(t *testing.T) {
	sut := Ec2Service{}

	actual := sut.handleError(mockedAwsError2{})
	if actual != false {
		t.Error("Error")
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
