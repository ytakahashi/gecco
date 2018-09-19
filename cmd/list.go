package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
)

var listOpts = &config.ListOption{}

func newListCmd(command iListCommand) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "lists EC2 instances",
		Run: func(cmd *cobra.Command, args []string) {
			if err := command.runCommand(aws.Ec2{}); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	listCmd.Flags().StringVarP(&listOpts.TagKey, "tagKey", "", "", "filters instances by tag key")
	listCmd.Flags().StringVarP(&listOpts.TagValue, "tagValue", "", "", "filters instances by tag value")
	listCmd.Flags().StringVarP(&listOpts.Status, "status", "", "", "filters instances by status")

	return listCmd
}

type listCommand struct{}

type iListCommand interface {
	runCommand(aws.Ec2Client) error
}

func (c listCommand) runCommand(awsec2 aws.Ec2Client) (err error) {
	options := *listOpts
	err = options.IsValid()
	if err != nil {
		return err
	}

	instances, err := awsec2.GetInstances(options)
	if err != nil {
		return err
	}

	instances.Print(os.Stdout)
	return
}
