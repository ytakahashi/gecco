package aws

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ytakahashi/gecco/config"
)

type tag struct {
	key   string
	value string
}

type tags []tag

type ec2Instance struct {
	instanceID       string
	instanceType     string
	availabilityZone string
	privateIPAdress  string
	status           string
	tags             tags
}

// Ec2Instances contains EC2 instance info
type Ec2Instances []ec2Instance

// Print instances
func (instances Ec2Instances) Print(w io.Writer) {
	for _, i := range instances {
		var tag string
		if len(i.tags) > 0 {
			tag = "{ "
			for _, t := range i.tags {
				tag += t.key + ":" + t.value + " "
			}
			tag += "}"
		}

		fmt.Fprintln(w,
			i.instanceID,
			i.instanceType,
			i.availabilityZone,
			i.status,
			tag,
		)
	}
}

// DescribeEC2 returns EC2 instance info
func DescribeEC2(options config.ListOption) (instances Ec2Instances, err error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ec2Svc := ec2.New(sess)

	input := createInput(options)

	result, err := ec2Svc.DescribeInstances(&input)
	if err != nil {
		return instances, err
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

func createInput(options config.ListOption) ec2.DescribeInstancesInput {
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

func newEc2Instance(i ec2.Instance) ec2Instance {
	tags := make(tags, 0)
	for _, t := range i.Tags {
		tag := tag{key: *t.Key, value: *t.Value}
		tags = append(tags, tag)
	}

	var az string
	p := i.Placement
	if p != nil {
		az = aws.StringValue(p.AvailabilityZone)
	}

	var status string
	s := i.State
	if s != nil {
		status = aws.StringValue(s.Name)
	}

	return ec2Instance{
		instanceID:       aws.StringValue(i.InstanceId),
		instanceType:     aws.StringValue(i.InstanceType),
		availabilityZone: az,
		privateIPAdress:  aws.StringValue(i.PrivateIpAddress),
		status:           status,
		tags:             tags,
	}
}
