package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ytakahashi/gecco/config"
)

// IEc2Service is an interface for ec2 services
type IEc2Service interface {
	initEc2Service() *ec2.EC2
	start(*ec2.EC2, bool, string) (*ec2.StartInstancesOutput, error)
	stop(*ec2.EC2, bool, string) (*ec2.StopInstancesOutput, error)
	handleError(error) bool
}

// Ec2Service Ec2 Service
type Ec2Service struct{}

func (s Ec2Service) initEc2Service() *ec2.EC2 {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return ec2.New(sess)
}

func (s Ec2Service) start(ec2Svc *ec2.EC2, dryRun bool, instanceID string) (*ec2.StartInstancesOutput, error) {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		DryRun: aws.Bool(dryRun),
	}
	return ec2Svc.StartInstances(input)
}

func (s Ec2Service) stop(ec2Svc *ec2.EC2, dryRun bool, instanceID string) (*ec2.StopInstancesOutput, error) {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		DryRun: aws.Bool(dryRun),
	}
	return ec2Svc.StopInstances(input)
}

func (s Ec2Service) handleError(err error) bool {
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		return true
	}
	return false
}

func (instances Ec2Instances) toStringSlice() []string {
	sl := make([]string, 0)
	for _, i := range instances {
		sl = append(sl, i.instanceID+" ("+i.status+"), Tags="+i.tags.toString())
	}
	return sl
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
