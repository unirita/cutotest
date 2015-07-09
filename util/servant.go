package util

import (
	"path/filepath"
)

type Servant struct {
	command

	config string
}

func NewServant() *Servant {
	s := new(Servant)
	s.GeneratePath("servant")
	return s
}

func (s *Servant) UseConfig(filename string) {
	s.config = filepath.Join(cutoroot, "bin", filename)
}

func (s *Servant) Start() error {
	if len(s.config) == 0 {
		return s.AsyncExec()
	} else {
		return s.AsyncExec("-c", s.config)
	}
}
