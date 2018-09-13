package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists EC2 instances",
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

type tag struct {
	key   string
	value string
}

type ec2Instance struct {
	instanceID       string
	instanceType     string
	availabilityZone string
	privateIPAdress  string
	tags             []tag
}

type ec2Instances []ec2Instance

func (instances ec2Instances) print() {
	for _, i := range instances {
		var tag string
		for _, t := range i.tags {
			tag = "{key:" + t.key + ",value:" + t.value + "}"
		}
		fmt.Println(
			i.instanceID,
			i.instanceType,
			i.availabilityZone,
			tag,
		)
	}
}

func newEc2Instance(i ec2.Instance) ec2Instance {
	tags := make([]tag, 0)
	for _, t := range i.Tags {
		tags = append(tags, tag{key: *t.Key, value: *t.Value})
	}

	return ec2Instance{
		instanceID:       *i.InstanceId,
		instanceType:     *i.InstanceType,
		availabilityZone: *i.Placement.AvailabilityZone,
		privateIPAdress:  *i.PrivateIpAddress,
		tags:             tags,
	}

}

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ec2Svc := ec2.New(sess)

	result, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		for _, r := range result.Reservations {
			instances := make(ec2Instances, 0)
			for _, i := range r.Instances {
				instances = append(instances, newEc2Instance(*i))
			}
			instances.print()
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
