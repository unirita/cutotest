package util

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"
)

type ConfigParam struct {
	CutoRoot string
}

var goPath string = os.Getenv("GOPATH")
var basePath string = filepath.Join(goPath, "src", "github.com", "unirita", "cutotest")

func InitCutoRoot() {
	os.MkdirAll(cutoroot, 0755)
	ClearDir(filepath.Join(cutoroot, "bpmn"))
	ClearDir(filepath.Join(cutoroot, "data"))
	ClearDir(filepath.Join(cutoroot, "joblog"))
	ClearDir(filepath.Join(cutoroot, "jobscript"))
	ClearDir(filepath.Join(cutoroot, "log"))
	ClearDir(filepath.Join(cutoroot, "temp"))

	emptyDbPath := filepath.Join(basePath, "util", "empty.sqlite")
	testDbPath := filepath.Join(cutoroot, "data", "cuto.sqlite")
	CopyFile(emptyDbPath, testDbPath)
}

func DeployTestData(testname string) {
	srcDir := filepath.Join(basePath, testname, "data", osname)
	targetDir := cutoroot
	CopyDir(srcDir, targetDir)
}

func ComplementConfig(filename string) error {
	path := filepath.Join(cutoroot, "bin", filename)
	tpl, err := template.ParseFiles(path)
	if err != nil {
		return errors.New("Config parse error.")
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.New("Failed to open config file.")
	}
	defer file.Close()

	param := ConfigParam{CutoRoot: cutoroot}
	if err := tpl.Execute(file, param); err != nil {
		return err
	}

	return nil
}
