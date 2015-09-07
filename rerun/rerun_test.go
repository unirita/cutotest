package rerun

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

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

func TestRerunParallel(t *testing.T) {
	defer util.SaveEvidence("rerun", "flow_parallel")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("flow_parallel.ini")
	rc, err := m.Rerun("1")
	if err != nil {
		t.Fatalf("Master start failed: %")
	}
	if rc != 0 {
		t.Fatalf("Master RC is not %d", 0)
	}

	if isExecuted(1, "job2", masterLog) {
		t.Errorf("JOB [%s] must not be executed, but it was.", "job2")
	}
	if isExecuted(1, "job4", masterLog) {
		t.Errorf("JOB [%s] must not be executed, but it was.", "job4")
	}
	if !isExecuted(1, "job3", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job3")
	}
	if !isExecuted(1, "job5", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job5")
	}
}

func TestRerun_Failcase(t *testing.T) {
	defer util.SaveEvidence("rerun", "failcase")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	correctPath := filepath.Join(util.GetCutoRoot(), "jobscript", getScriptFileName("job4"))
	wrongPath := filepath.Join(util.GetCutoRoot(), "jobscript", getScriptFileName("_job4"))

	os.Rename(correctPath, wrongPath)
	testRerun_Failcase_first(t)

	time.Sleep(500 * time.Millisecond)

	logPath := filepath.Join(util.GetCutoRoot(), "log", "master.log")
	bkupPath := filepath.Join(util.GetCutoRoot(), "log", "master_first.log")
	os.Rename(logPath, bkupPath)

	os.Rename(wrongPath, correctPath)
	testRerun_Failcase_second(t)
}

func testRerun_Failcase_first(t *testing.T) {
	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("failcase.ini")
	rc, err := m.Rerun("1")
	if err != nil {
		t.Fatalf("Master start failed: %")
	}
	if rc == 0 {
		t.Fatalf("Master RC must not be %d", 0)
	}

	if !isExecuted(1, "job1", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job1")
	}
	if !isExecuted(1, "job2", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job2")
	}
	if !isExecuted(1, "job3", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job3")
	}
	if !isExecuted(1, "job4", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job4")
	}
	if isExecuted(1, "job5", masterLog) {
		t.Errorf("JOB [%s] must not be executed, but it was.", "job5")
	}
	if isExecuted(1, "job6", masterLog) {
		t.Errorf("JOB [%s] must not be executed, but it was.", "job6")
	}
}

func testRerun_Failcase_second(t *testing.T) {
	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("failcase.ini")
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
	if isExecuted(1, "job3", masterLog) {
		t.Errorf("JOB [%s] must not be executed, but it was.", "job3")
	}
	if !isExecuted(1, "job4", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job4")
	}
	if !isExecuted(1, "job5", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job5")
	}
	if !isExecuted(1, "job6", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job6")
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
