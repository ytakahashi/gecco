package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
)

var listOpts = &config.ListOption{}

func newListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "lists EC2 instances",
		Run: func(cmd *cobra.Command, args []string) {
			if err := (*listOpts).IsValid(); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			list(*listOpts, aws.DescribeEC2)
		},
	}

	listCmd.Flags().StringVarP(&listOpts.TagKey, "tagKey", "", "", "filters instances by tag key")
	listCmd.Flags().StringVarP(&listOpts.TagValue, "tagValue", "", "", "filters instances by tag value")
	listCmd.Flags().StringVarP(&listOpts.Status, "status", "", "", "filters instances by status")

	return listCmd
}

func list(
	options config.ListOption,
	getInstances func(config.ListOption) (instances aws.Ec2Instances, err error),
) error {
	instances, err := getInstances(options)
	if err != nil {
		return err
	}
	instances.Print(os.Stdout)

	return nil
}
