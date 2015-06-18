package util

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"syscall"
)

type Show struct {
	Command    string
	ConfigPath string
	Stdout     string
	Stderr     string
	cmd        *exec.Cmd
	params     []string
}

func NewShow() *Show {
	s := new(Show)
	if osname == "windows" {
		s.Command = filepath.Join(cutoroot, "bin", "show.exe")
	} else {
		s.Command = filepath.Join(cutoroot, "bin", "show")
	}
	return s
}

func (s *Show) SetConfig(filename string) {
	s.ConfigPath = filepath.Join(cutoroot, "bin", filename)
}

func (s *Show) Help() (int, error) {
	s.cmd = exec.Command(s.Command, "-help")
	return s.exec()
}

func (s *Show) Version() (int, error) {
	s.cmd = exec.Command(s.Command, "-v")
	return s.exec()
}

func (s *Show) Run() (int, error) {
	if len(s.ConfigPath) == 0 {
		s.cmd = exec.Command(s.Command)
	} else {
		if len(s.params) == 0 {
			s.cmd = exec.Command(s.Command, "-c", s.ConfigPath)
		} else {
			s.params = append(s.params, "-c="+s.ConfigPath)
			s.cmd = exec.Command(s.Command, s.params...)
		}
	}
	return s.exec()
}

func (s *Show) exec() (int, error) {
	outbuf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	s.cmd.Stdout = outbuf
	s.cmd.Stderr = errbuf
	if err := s.cmd.Start(); err != nil {
		return -1, err
	}
	defer func() {
		s.Stdout = outbuf.String()
		s.Stderr = errbuf.String()
	}()

	err := s.cmd.Wait()
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

func (s *Show) AddJobnet(jobnet string) {
	s.params = append(s.params, "-jobnet="+jobnet)
}

func (s *Show) AddFrom(date string) {
	s.params = append(s.params, "-from="+date)
}

func (s *Show) AddTo(date string) {
	s.params = append(s.params, "-to="+date)
}

func (s *Show) AddStatus(status string) {
	s.params = append(s.params, "-status="+status)
}

func (s *Show) AddFormat(format string) {
	s.params = append(s.params, "-format="+format)
}
