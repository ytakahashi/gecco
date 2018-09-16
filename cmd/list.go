package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists EC2 instances",
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateListOpts(*listOpts); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		list()
	},
}

var listOpts = &config.ListOptions{}

func isValidStatus(e string) bool {
	status := []string{"running", "stopping", "pending", "shutting-down", "terminated", "stopped"}
	for _, v := range status {
		if e == v {
			return true
		}
	}
	return false
}

func validateListOpts(listOpts config.ListOptions) error {
	if listOpts.Status != "" && !isValidStatus(listOpts.Status) {
		return fmt.Errorf("Invalid status (%s)", listOpts.Status)
	}

	if listOpts.TagKey == "" && listOpts.TagValue != "" {
		return errors.New("Option '--tagKey' is required when '--tagValue' is specified")
	}

	if listOpts.TagKey != "" && listOpts.TagValue == "" {
		return errors.New("Option '--tagValue' is required when '--tagKey' is specified")
	}

	return nil
}

func list() error {
	instances, err := aws.DescribeEC2(*listOpts)
	if err != nil {
		return err
	}
	instances.Print()

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
