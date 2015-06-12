package util

import (
	"errors"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

type Master struct {
	Command    string
	ConfigPath string
	Timeout    int
	cmd        *exec.Cmd
}

func NewMaster() *Master {
	m := new(Master)
	if osname == "windows" {
		m.Command = filepath.Join(cutoroot, "bin", "master.exe")
	} else {
		m.Command = filepath.Join(cutoroot, "bin", "master")
	}
	return m
}

func (m *Master) SetConfig(filename string) {
	m.ConfigPath = filepath.Join(cutoroot, "bin", filename)
}

func (m *Master) SyntaxCheck(jobnet string) (int, error) {
	if len(m.ConfigPath) == 0 {
		m.cmd = exec.Command(m.Command, "-n", jobnet)
	} else {
		m.cmd = exec.Command(m.Command, "-n", jobnet, "-c", m.ConfigPath)
	}
	return m.exec()
}

func (m *Master) Run(jobnet string) (int, error) {
	if len(m.ConfigPath) == 0 {
		m.cmd = exec.Command(m.Command, "-n", jobnet, "-s")
	} else {
		m.cmd = exec.Command(m.Command, "-n", jobnet, "-s", "-c", m.ConfigPath)
	}
	return m.exec()
}

func (m *Master) exec() (int, error) {
	if err := m.cmd.Start(); err != nil {
		return -1, err
	}

	err := m.waitTimeout()
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

func (m *Master) waitTimeout() error {
	if m.Timeout == 0 {
		return m.cmd.Wait()
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- m.cmd.Wait()
	}()

	t := time.Duration(m.Timeout) * time.Second
	select {
	case err := <-ch:
		return err
	case <-time.After(t):
		m.cmd.Process.Kill()
		return errors.New("Process timeout.")
	}

	return nil
}
