package util

import (
	"path/filepath"
)

type Master struct {
	command

	config string
}

func NewMaster() *Master {
	m := new(Master)
	m.GeneratePath("master")
	return m
}

func (m *Master) UseConfig(filename string) {
	m.config = filepath.Join(cutoroot, "bin", filename)
}

func (m *Master) SyntaxCheck(jobnet string) (int, error) {
	if len(m.config) == 0 {
		return m.Exec("-n", jobnet)
	} else {
		return m.Exec("-n", jobnet, "-c", m.config)
	}
}

func (m *Master) Run(jobnet string) (int, error) {
	if len(m.config) == 0 {
		return m.Exec("-n", jobnet, "-s")
	} else {
		return m.Exec("-n", jobnet, "-s", "-c", m.config)
	}
}

func (m *Master) Rerun(instanceID string) (int, error) {
	if len(m.config) == 0 {
		return m.Exec("-r", instanceID)
	} else {
		return m.Exec("-r", instanceID, "-c", m.config)
	}
}
