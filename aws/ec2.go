package aws

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ytakahashi/gecco/config"
	"github.com/ytakahashi/gecco/ext"
)

type tag struct {
	key   string
	value string
}

func (t tag) toString() string {
	return "{\"" + t.key + "\": \"" + t.value + "\"}"
}

type tags []tag

func (tags tags) toString() (str string) {
	if len(tags) > 0 {
		str = "["
		for _, t := range tags {
			str += t.toString() + ", "
		}
		str = strings.TrimRight(str, ", ")
		str += "]"
	}
	return
}

// Ec2Instance ec2 instance used in this app
type Ec2Instance struct {
	instanceID   string
	instanceType string
	status       string
	tags         tags
}

// Ec2 contains EC2 instance info
type Ec2 struct{}

// Ec2Client Ec2 Client
type Ec2Client interface {
	GetInstances(config.FilterOption, IEc2Service) (Ec2Instances, error)
	StartInstance(string, IEc2Service) error
	StopInstance(string, IEc2Service) error
}

// StartInstance starts target instance
func (e Ec2) StartInstance(target string, service IEc2Service) error {
	ec2Svc := service.initEc2Service()
	result, err := service.start(ec2Svc, true, target)

	if service.handleError(err) {
		result, err = service.start(ec2Svc, false, target)
		if err != nil {
			return err
		}

		fmt.Println("Success:", result.StartingInstances)
		return nil
	}
	return err
}

// StopInstance stops target instance
func (e Ec2) StopInstance(target string, service IEc2Service) error {
	ec2Svc := service.initEc2Service()

	result, err := service.stop(ec2Svc, true, target)

	if service.handleError(err) {
		result, err = service.stop(ec2Svc, false, target)
		if err != nil {
			return err
		}

		fmt.Println("Success:", result.StoppingInstances)
		return nil
	}
	return err
}

// GetInstances returns Instances
func (e Ec2) GetInstances(options config.FilterOption, service IEc2Service) (instances Ec2Instances, err error) {
	ec2Svc := service.initEc2Service()

	input := createInput(options)

	result, err := ec2Svc.DescribeInstances(&input)
	if err != nil {
		return nil, err
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

// Instances instances
type Instances interface {
	GetFilteredInstances(ext.ICommand) (string, error)
	ToString(config.OutputFormat) (string, error)
}

// Ec2Instances contains EC2 instance info
type Ec2Instances []Ec2Instance

// GetFilteredInstances returns isntanceID of a selected instance
func (instances Ec2Instances) GetFilteredInstances(filter ext.ICommand) (selected string, err error) {
	records := instances.toStringSlice()
	var text string
	for _, r := range records {
		text += r + "\n"
	}

	var buf bytes.Buffer
	cmd := filter.CreateCommand(strings.NewReader(text), &buf, os.Stderr)
	err = cmd.Run()
	if err != nil {
		return
	}

	if buf.Len() == 0 {
		err = errors.New("No line is selected")
		return
	}

	selected = strings.TrimSpace(buf.String())
	selected = strings.Split(selected, " ")[0]
	return
}

// ToString returns string value of Ec2Instance
func (instances Ec2Instances) ToString(outputFormat config.OutputFormat) (string, error) {
	switch outputFormat {
	case config.Text:
		b := make([]byte, 0, len(instances))
		for _, i := range instances {
			s := i.instanceID + " " + i.instanceType + " " + i.status + " " + i.tags.toString() + "\n"
			b = append(b, s...)
		}
		return string(b), nil

	case config.JSON:
		b, _ := json.MarshalIndent(&instances, "", "    ")
		return string(b), nil

	default:
		return "", errors.New("Unexpected OutputFormat")
	}
}
