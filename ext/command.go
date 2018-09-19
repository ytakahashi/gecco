package ext

import (
	"io"
	"os/exec"
)

// ICommand is an interface of command
type ICommand interface {
	CreateCommand(io.Reader, io.Writer, io.Writer) *exec.Cmd
}

// Command holds command arguments (including command itself)
type Command struct {
	Args []string
}

// CreateCommand created command
func (c Command) CreateCommand(i io.Reader, o io.Writer, e io.Writer) *exec.Cmd {
	command := exec.Command(c.Args[0])
	command.Args = c.Args
	command.Stdin = i
	command.Stdout = o
	command.Stderr = e
	return command
}
