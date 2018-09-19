package config

import (
	"errors"
	"fmt"
)

// Config file
type Config struct {
	InteractiveFilterCommand string
}

// Conf Config
var Conf Config

// ListOption stores options for list command
type ListOption struct {
	TagKey   string
	TagValue string
	Status   string
}

// ConnectOptions stores options for connect command
type ConnectOptions struct {
	Target      string
	Interactive bool
}

// IsValid returns true if optons are valid
func (option ListOption) IsValid() error {
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
