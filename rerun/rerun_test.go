package rerun

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/unirita/cutotest/util"
)

var masterLog = filepath.Join(util.GetCutoRoot(), "log", "master.log")

func TestRerunSerial(t *testing.T) {
	defer util.SaveEvidence("rerun", "flow_serial")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("flow_serial.ini")
	rc, err := m.Rerun("1")
	if err != nil {
		t.Fatalf("Master start failed: %")
	}
	if rc != 0 {
		t.Fatalf("Master RC is not %d", 0)
	}

	if isExecuted(1, "job1", masterLog) {
		t.Errorf("JOB [%s] must not be executed, but it was.", "job1")
	}
	if isExecuted(1, "job2", masterLog) {
		t.Errorf("JOB [%s] must not be executed, but it was.", "job2")
	}
	if !isExecuted(1, "job3", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job3")
	}
	if !isExecuted(1, "job4", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job4")
	}
}

func isExecuted(instanceID int, jobName string, logfile string) bool {
	file, err := os.Open(logfile)
	if err != nil {
		return false
	}
	defer file.Close()

	ptnStr := fmt.Sprintf(`CTM023I JOB \[%s\].*INSTANCE \[%d\]`, jobName, instanceID)
	matcher := regexp.MustCompile(ptnStr)

	s := bufio.NewScanner(file)
	for s.Scan() {
		if matcher.MatchString(s.Text()) {
			return true
		}
	}

	return false
}
