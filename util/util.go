package util

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
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
	//ClearDir(filepath.Join(cutoroot, "temp"))

	emptyDbPath := filepath.Join(basePath, "util", "empty.sqlite")
	testDbPath := filepath.Join(cutoroot, "data", "cuto.sqlite")
	CopyFile(emptyDbPath, testDbPath)
}

func DeployTestData(testname string) {
	srcDir := filepath.Join(basePath, testname, "data", osname)
	targetDir := cutoroot
	CopyDir(srcDir, targetDir)
	PermitExecRecursive(targetDir)
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

func HasLogError(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return true
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	hasError := false
	for s.Scan() {
		line := s.Text()

		if strings.Contains(line, "[ERR]") {
			hasError = true
			break
		}
		errCodePtn := regexp.MustCompile("CT[MS][0-9]{3}E")
		if errCodePtn.MatchString(line) {
			hasError = true
			break
		}
	}

	return hasError
}

func GetCutoRoot() string {
	return cutoroot
}

func GetLogPath(filename string) string {
	return filepath.Join(cutoroot, "log", filename)
}

func GetDBDirPath() string {
	return filepath.Join(cutoroot, "data", "cuto.sqlite")
}

func FindJoblog(dirname string, nid int, jobname string) []string {
	now := time.Now()
	datestr := fmt.Sprintf("%04d%02d%02d", now.Year(), now.Month(), now.Day())
	dirpath := filepath.Join(cutoroot, dirname, datestr)
	joblogs := make([]string, 0)
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return joblogs
	}
	if !dirinfo.IsDir() {
		return joblogs
	}

	fis, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return joblogs
	}

	prefix := fmt.Sprintf("%d.%s", nid, jobname)
	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), prefix) {
			joblogs = append(joblogs, filepath.Join(dirpath, fi.Name()))
		}
	}

	return joblogs
}

func SaveEvidence(names ...string) {
	base := filepath.Join(os.Getenv("GOPATH"), "evidence")
	for _, name := range names {
		base = filepath.Join(base, name)
	}

	dataFrom := filepath.Join(cutoroot, "data")
	dataTo := filepath.Join(base, "data")
	joblogFrom := filepath.Join(cutoroot, "joblog")
	joblogTo := filepath.Join(base, "joblog")
	logFrom := filepath.Join(cutoroot, "log")
	logTo := filepath.Join(base, "log")

	ClearDir(dataTo)
	ClearDir(joblogTo)
	ClearDir(logTo)
	CopyDir(dataFrom, dataTo)
	CopyDir(joblogFrom, joblogTo)
	CopyDir(logFrom, logTo)
}
