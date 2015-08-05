package container

import (
	"fmt"
	"os/exec"
	"strings"
)

type Container struct {
	image     string
	name      string
	isRunning bool
	timezone  string
}

const (
	defaultTimezone = "/etc/localtime"
	zoneInfoPath    = "/usr/share/zoneinfo"
)

func New(image string, name string) *Container {
	c := new(Container)
	c.image = image
	c.name = name
	c.isRunning = false
	c.timezone = defaultTimezone
	return c
}

func (c *Container) Start() error {
	tzOption := fmt.Sprintf("%s:/etc/localtime:ro", c.timezone)
	cmd := exec.Command("docker", "run", "-itd", "-v", tzOption, "--name", c.name, c.image)
	err := cmd.Run()
	if err != nil {
		return err
	}

	c.isRunning = true
	return nil
}

func (c *Container) SetTimezone(zoneName string) {
	c.timezone = fmt.Sprintf("%s/%s", zoneInfoPath, zoneName)
}

func (c *Container) ResetTimezone() {
	c.timezone = defaultTimezone
}

func (c *Container) IPAddress() (string, error) {
	out, err := exec.Command("docker", "inspect",
		"--format='{{.NetworkSettings.IPAddress}}'", c.name).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func (c *Container) Terminate() {
	if c.isRunning {
		exec.Command("docker", "stop", c.name).Run()
		exec.Command("docker", "rm", "-f", c.name).Run()
		c.isRunning = false
	}
}
