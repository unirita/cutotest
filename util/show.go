package util

import (
	"path/filepath"
)

type Show struct {
	command

	Config string
	params []string
}

func NewShow() *Show {
	s := new(Show)
	s.GeneratePath("show")
	return s
}

func (s *Show) UseConfig(filename string) {
	s.Config = filepath.Join(cutoroot, "bin", filename)
}

func (s *Show) Help() (int, error) {
	return s.Exec("-help")
}

func (s *Show) Version() (int, error) {
	return s.Exec("-v")
}

func (s *Show) Run() (int, error) {
	if len(s.Config) == 0 {
		return s.Exec(s.Path)
	} else {
		if len(s.params) == 0 {
			return s.Exec("-c", s.Config)
		} else {
			s.params = append(s.params, "-c="+s.Config)
			return s.Exec(s.params...)
		}
	}
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

func (s *Show) AddUTCOption() {
	s.params = append(s.params, "-utc")
}
