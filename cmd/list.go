package cmd

import (
	"encoding/json"
	"errors"
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

func printInstances(instances aws.Ec2Instances, w io.Writer, outputFormat config.OutputFormat) error {
	switch outputFormat {
	case config.Text:
		for _, i := range instances {
			fmt.Fprintln(w,
				i.InstanceID,
				i.InstanceType,
				i.Status,
				i.Tags.ToString(),
			)
		}
		return nil

	case config.JSON:
		b, _ := json.MarshalIndent(&instances, "", "    ")
		fmt.Fprint(w, string(b))
		return nil

	default:
		return errors.New("Unexpected OutputFormat")
	}
}
