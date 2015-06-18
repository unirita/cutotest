package util

import (
	"os/exec"
)

type Container struct {
	image     string
	name      string
	isRunning bool
}

func New(image string, name string) *Container {
	c := new(Container)
	c.image = image
	c.name = name
	c.isRunning = false
	return c
}

func (c *Container) Start() error {
	err := exec.Command("docker", "run", "-itd", "--name", c.name, c.image).Run()
	if err != nil {
		return err
	}

	c.isRunning = true
	return nil
}

func (c *Container) IPAddress() (string, error) {
	out, err := exec.Command("docker", "inspect",
		"--format='{{.NetworkSettings.IPAddress}}'", c.name).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (c *Container) Terminate() {
	if c.isRunning {
		exec.Command("docker", "stop", c.name).Run()
		exec.Command("docker", "rm", "-f", c.name).Run()
		c.isRunning = false
	}
}
