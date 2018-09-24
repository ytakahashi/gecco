package config

import (
	"errors"
	"fmt"
)

// IOption is an interface for option.
type IOption interface {
	IsValid() error
}

// FilterOption stores options for list command.
type FilterOption struct {
	Status   string
	TagKey   string
	TagValue string
}

// TargetOption stores options for connect/start/stop command.
type TargetOption struct {
	// Target instanceId.
	Target string
	// A flag whether to select ec2 interactively.
	Interactive bool
}

// IsValid validates given options.
func (option TargetOption) IsValid() error {
	if option.Target == "" && option.Interactive == false {
		return errors.New("Option '--target' or '-i' is required")
	}

	if option.Target != "" && option.Interactive == true {
		return errors.New("Options '--target' and '-i' cannot be used at the same time")
	}

	return nil
}

// IsValid validates given options.
func (option FilterOption) IsValid() error {
	if option.Status != "" && !isValidStatus(option.Status) {
		return fmt.Errorf("Invalid status (%s)", option.Status)
	}

	if option.TagKey == "" && option.TagValue != "" {
		return errors.New("Option '--tagKey' is required when '--tagValue' is specified")
	}

	if option.TagKey != "" && option.TagValue == "" {
		return errors.New("Option '--tagValue' is required when '--tagKey' is specified")
	}

	return nil
}

func isValidStatus(e string) bool {
	status := []string{"running", "stopping", "pending", "shutting-down", "terminated", "stopped"}
	for _, v := range status {
		if e == v {
			return true
		}
	}
	return false
}
