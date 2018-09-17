package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ytakahashi/gecco/aws"
	"github.com/ytakahashi/gecco/config"
)

var connectOpts = &config.ConnectOptions{}

func newConnectCmd() *cobra.Command {
	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "connect to EC2 instance",
		Long:  "connect to EC2 instance using 'aws cli start-session' command",
		Run: func(cmd *cobra.Command, args []string) {
			err := connect(*connectOpts, doRun, initConfig)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	connectCmd.Flags().StringVarP(&connectOpts.Target, "target", "", "", "target instanceId to start session")
	connectCmd.Flags().BoolVarP(&connectOpts.Interactive, "interactive", "i", false, "Select a value interactively (requires config file)")

	return connectCmd
}

func connect(
	options config.ConnectOptions,
	run func(string) error,
	init func() error,
) error {
	var target string
	if options.Interactive {
		if err := init(); err != nil {
			return err
		}

		i, err := aws.DescribeEC2(config.ListOption{})
		if err != nil {
			return err
		}
		sl := i.ToStringSlice()

		target, err = filter(conf, sl)
		if err != nil {
			return err
		}
	} else {
		if options.Target == "" {
			return errors.New("Option '--target' is not specified")
		}
		target = options.Target
	}
	return run(target)

}

func doRun(target string) error {
	return createCommand([]string{"aws", "ssm", "start-session", "--target", target}, os.Stdin, os.Stdout, os.Stderr).Run()
}

func filter(conf config.Config, records []string) (selected string, err error) {
	var text string
	for _, r := range records {
		text += r + "\n"
	}

	buf, err := doFilter(conf, text)

	if buf.Len() == 0 {
		err = errors.New("No line is selected")
		return
	}

	selected = strings.TrimSpace(buf.String())
	return
}

func doFilter(conf config.Config, text string) (buf bytes.Buffer, err error) {
	cmd := createCommand([]string{conf.InteractiveFilterCommand}, os.Stderr, &buf, strings.NewReader(text))
	err = cmd.Run()
	return
}

func createCommand(commandWithArgs []string, e, o io.Writer, i io.Reader) *exec.Cmd {
	command := exec.Command(commandWithArgs[0])
	command.Args = commandWithArgs
	command.Stderr = e
	command.Stdout = o
	command.Stdin = i
	return command
}
