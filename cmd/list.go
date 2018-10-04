package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
)

func newListCmd(command iListCommand) *cobra.Command {
	listOpts := &config.FilterOption{}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "lists EC2 instances",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			command.initListCommand(*listOpts, aws.Ec2{})
			err = command.runCommand()
			return
		},
	}

	listCmd.Flags().StringVarP(&listOpts.TagKey, "tagKey", "", "", "filters instances by tag key")
	listCmd.Flags().StringVarP(&listOpts.TagValue, "tagValue", "", "", "filters instances by tag value")
	listCmd.Flags().StringVarP(&listOpts.Status, "status", "", "", "filters instances by status")
	listCmd.Flags().StringVarP(&listOpts.OutputFormat, "output", "o", "text", "filters instances by status")

	return listCmd
}

type iListCommand interface {
	initListCommand(config.FilterOption, aws.Ec2Client)
	runCommand() error
}

type listCommand struct {
	options   config.FilterOption
	ec2Client aws.Ec2Client
}

func (c *listCommand) initListCommand(o config.FilterOption, client aws.Ec2Client) {
	c.ec2Client = client
	c.options = o
}

func (c listCommand) runCommand() (err error) {
	options := c.options

	err = options.IsValid()
	if err != nil {
		return err
	}

	instances, err := c.ec2Client.GetInstances(options, aws.Ec2Service{})
	if err != nil {
		return err
	}

	return printInstances(instances, os.Stdout, options.GetOutputFormat())
}

func printInstances(instances aws.Instances, w io.Writer, outputFormat config.OutputFormat) error {
	str, err := instances.ToString(outputFormat)
	if err != nil {
		return err
	}

	fmt.Fprint(w, str)
	return nil
}
