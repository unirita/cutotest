package util

import (
	//"os"
	"os/exec"
	"path/filepath"
)

type Servant struct {
	Command    string
	ConfigPath string

	cmd *exec.Cmd
}

func NewServant() *Servant {
	s := new(Servant)
	s.Command = filepath.Join(cutoroot, "bin", "servant.exe")
	return s
}

func (s *Servant) SetConfig(filename string) {
	s.ConfigPath = filepath.Join(cutoroot, "bin", filename)
}

func (s *Servant) Start() error {
	if len(s.ConfigPath) == 0 {
		s.cmd = exec.Command(s.Command)
	} else {
		s.cmd = exec.Command(s.Command, "-c", s.ConfigPath)
	}
	return s.cmd.Start()
}

func (s *Servant) Kill() {
	s.cmd.Process.Kill()
}
