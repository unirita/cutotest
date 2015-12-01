package util

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

type command struct {
	Path        string
	Timeout     int
	Stdout      string
	Stderr      string
	cmd         *exec.Cmd
	asyncStdout *bytes.Buffer
	asyncStderr *bytes.Buffer
}

func (c *command) GeneratePath(name string) {
	c.Path = filepath.Join(cutoroot, "bin", name)
	if runtime.GOOS == "windows" {
		c.Path += ".exe"
	}
}

func (c *command) Exec(arg ...string) (int, error) {
	c.cmd = exec.Command(c.Path, arg...)
	outbuf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	c.cmd.Stdout = outbuf
	c.cmd.Stderr = errbuf
	if err := c.cmd.Start(); err != nil {
		return -1, err
	}
	defer func() {
		c.Stdout = outbuf.String()
		c.Stderr = errbuf.String()
	}()

	err := c.waitTimeout()
	if err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				return s.ExitStatus(), nil
			} else {
				err = errors.New("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus.")
				return -1, err
			}
		}
	}
	return 0, nil
}

func (c *command) AsyncExec(arg ...string) error {
	c.asyncStdout = new(bytes.Buffer)
	c.asyncStderr = new(bytes.Buffer)
	c.cmd.Stdout = c.asyncStdout
	c.cmd.Stderr = c.asyncStderr
	c.cmd = exec.Command(c.Path, arg...)
	return c.cmd.Start()
}

func (c *command) Kill() {
	c.Stdout = c.asyncStdout.String()
	c.Stderr = c.asyncStderr.String()
	c.cmd.Process.Kill()
}

func (c *command) waitTimeout() error {
	if c.Timeout <= 0 {
		return c.cmd.Wait()
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- c.cmd.Wait()
	}()

	t := time.Duration(c.Timeout) * time.Second
	select {
	case err := <-ch:
		return err
	case <-time.After(t):
		c.cmd.Process.Kill()
		return errors.New("Process timeout.")
	}

	return nil
}
