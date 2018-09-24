package aws

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ytakahashi/gecco/config"
	"github.com/ytakahashi/gecco/ext"
)

type tag struct {
	key   string
	value string
}

func (t tag) toString() string {
	return "{\"" + t.key + "\": \"" + t.value + "\"}"
}

type tags []tag

func (tags tags) toString() (str string) {
	if len(tags) > 0 {
		str = "["
		for _, t := range tags {
			str += t.toString() + ", "
		}
		str = strings.TrimRight(str, ", ")
		str += "]"
	}
	return
}

// IAwsService is an interface for Aws services
type IAwsService interface {
	initEc2Service() *ec2.EC2
}

// Ec2Service Ec2 Service
type Ec2Service struct{}

func (s Ec2Service) initEc2Service() *ec2.EC2 {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return ec2.New(sess)
}

// Ec2Instance Ec2Instance
type Ec2Instance struct {
	instanceID   string
	instanceType string
	status       string
	tags         tags
}

// Ec2 contains EC2 instance info
type Ec2 struct{}

// Ec2Client Ec2 Client
type Ec2Client interface {
	GetInstances(config.FilterOption, IAwsService) (Ec2Instances, error)
	StartInstance(string, IAwsService) error
	StopInstance(string, IAwsService) error
}

// StartInstance starts target instance
func (e Ec2) StartInstance(target string, service IAwsService) error {
	ec2Svc := service.initEc2Service()

	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(target),
		},
		DryRun: aws.Bool(true),
	}
	result, err := ec2Svc.StartInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = ec2Svc.StartInstances(input)
		if err != nil {
			return err
		}
		fmt.Println("Success:", result.StartingInstances)
		return nil

	}
	return err
}

// StopInstance stops target instance
func (e Ec2) StopInstance(target string, service IAwsService) error {
	ec2Svc := service.initEc2Service()

	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(target),
		},
		DryRun: aws.Bool(true),
	}
	result, err := ec2Svc.StopInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = ec2Svc.StopInstances(input)
		if err != nil {
			return err
		}
		fmt.Println("Success:", result.StoppingInstances)
		return nil

	}
	return err
}

// GetInstances Get Instances
func (e Ec2) GetInstances(options config.FilterOption, service IAwsService) (instances Ec2Instances, err error) {
	ec2Svc := service.initEc2Service()

	input := createInput(options)

	result, err := ec2Svc.DescribeInstances(&input)
	if err != nil {
		return nil, err
	}

	instances = make(Ec2Instances, 0)
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			instance := newEc2Instance(*i)
			if instance.instanceID != "" {
				instances = append(instances, instance)
			}
		}
	}
	return instances, nil
}

// Instances instances
type Instances interface {
	Print(w io.Writer)
	GetFilteredInstances(ext.ICommand) (string, error)
}

// Ec2Instances contains EC2 instance info
type Ec2Instances []Ec2Instance

func (instances Ec2Instances) toStringSlice() []string {
	sl := make([]string, 0)
	for _, i := range instances {
		sl = append(sl, i.instanceID+",Tags="+i.tags.toString())
	}
	return sl
}

// Print instances
func (instances Ec2Instances) Print(w io.Writer) {
	for _, i := range instances {
		fmt.Fprintln(w,
			i.instanceID,
			i.instanceType,
			i.status,
			i.tags.toString(),
		)
	}
}

// GetFilteredInstances GetFilteredInstances
func (instances Ec2Instances) GetFilteredInstances(filter ext.ICommand) (selected string, err error) {
	records := instances.toStringSlice()
	var text string
	for _, r := range records {
		text += r + "\n"
	}

	var buf bytes.Buffer
	cmd := filter.CreateCommand(strings.NewReader(text), &buf, os.Stderr)
	err = cmd.Run()
	if err != nil {
		return
	}

	if buf.Len() == 0 {
		err = errors.New("No line is selected")
		return
	}

	selected = strings.TrimSpace(buf.String())
	selected = strings.Split(selected, ",")[0]
	return
}

func createInput(options config.FilterOption) ec2.DescribeInstancesInput {
	filters := make([]*ec2.Filter, 0)

	if options.Status != "" {
		filter := ec2.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []*string{aws.String(options.Status)},
		}
		filters = append(filters, &filter)
	} else {
		filter := ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
				aws.String("pending"),
				aws.String("stopping"),
				aws.String("shutting-down"),
				aws.String("terminated"),
				aws.String("stopped"),
			},
		}
		filters = append(filters, &filter)
	}

	if options.TagKey != "" && options.TagValue != "" {
		filter := ec2.Filter{
			Name:   aws.String("tag:" + options.TagKey),
			Values: []*string{aws.String(options.TagValue)},
		}
		filters = append(filters, &filter)
	}

	return ec2.DescribeInstancesInput{
		Filters: filters,
	}
}

func newEc2Instance(i ec2.Instance) Ec2Instance {
	tags := make(tags, 0)
	for _, t := range i.Tags {
		tag := tag{key: *t.Key, value: *t.Value}
		tags = append(tags, tag)
	}

	var status string
	s := i.State
	if s != nil {
		status = aws.StringValue(s.Name)
	}

	return Ec2Instance{
		instanceID:   aws.StringValue(i.InstanceId),
		instanceType: aws.StringValue(i.InstanceType),
		status:       status,
		tags:         tags,
	}
}
