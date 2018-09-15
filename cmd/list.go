package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists EC2 instances",
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateListOpts(); err != nil {
			fmt.Println("err", err)
			os.Exit(1)
		}
		main()
	},
}

type listOptions struct {
	tagKey   string
	tagValue string
	status   string
}

var listOpts = &listOptions{}

func isValidStatus(e string) bool {
	status := []string{"running", "stopping", "pending", "shutting-down", "terminated", "stopped"}
	for _, v := range status {
		if e == v {
			return true
		}
	}
	return false
}

func validateListOpts() error {
	if listOpts.status != "" && !isValidStatus(listOpts.status) {
		return fmt.Errorf("Invalid status (%s)", listOpts.status)
	}
	return nil
}

type tag struct {
	key   string
	value string
}

type tags []tag

func (tags tags) matches(k string, v string) bool {
	if k == "" && v == "" {
		return true
	}

	if k != "" && v == "" {
		for _, t := range tags {
			if k == t.key {
				return true
			}
		}
		return false
	}

	if k == "" && v != "" {
		for _, t := range tags {
			if v == t.value {
				return true
			}
		}
		return false
	}

	if len(tags) == 0 {
		return false
	}

	for _, t := range tags {
		if k == t.key && v == t.value {
			return true
		}
	}

	return false
}

type ec2Instance struct {
	instanceID       string
	instanceType     string
	availabilityZone string
	privateIPAdress  string
	status           string
	tags             tags
}

type ec2Instances []ec2Instance

func (instances ec2Instances) print() {
	for _, i := range instances {
		var tag string
		for _, t := range i.tags {
			tag += "{" + t.key + ":" + t.value + "}"
		}
		fmt.Println(
			i.instanceID,
			i.instanceType,
			i.availabilityZone,
			i.status,
			tag,
		)
	}
}

func newEc2Instance(i ec2.Instance) (ret ec2Instance) {
	status := *i.State.Name
	if listOpts.status != "" && listOpts.status != status {
		return
	}

	tags := make(tags, 0)
	for _, t := range i.Tags {
		tag := tag{key: *t.Key, value: *t.Value}
		tags = append(tags, tag)
	}

	if !tags.matches(listOpts.tagKey, listOpts.tagValue) {
		return
	}

	ret = ec2Instance{
		instanceID:       *i.InstanceId,
		instanceType:     *i.InstanceType,
		availabilityZone: *i.Placement.AvailabilityZone,
		privateIPAdress:  *i.PrivateIpAddress,
		status:           *i.State.Name,
		tags:             tags,
	}

	return
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
				instance := newEc2Instance(*i)
				if instance.instanceID != "" {
					instances = append(instances, instance)
				}
			}
			instances.print()
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
