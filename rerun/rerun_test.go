package rerun

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/unirita/cutotest/util"
	"github.com/unirita/cutotest/util/db"
)

var masterLog = filepath.Join(util.GetCutoRoot(), "log", "master.log")

func TestRerun_FlowSerial(t *testing.T) {
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

func TestRerun_FlowParallel(t *testing.T) {
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

func TestRerun_JobCheck_Normal(t *testing.T) {
	defer util.SaveEvidence("rerun", "jobchk_normal")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	s := util.NewServant()
	s.UseConfig("jobchk_normal_servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("jobchk_normal_master.ini")
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
	if !isExecuted(1, "job2", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job2")
	}

	dbfile := filepath.Join(util.GetCutoRoot(), "data", "jobchk_normal.sqlite")
	conn, err := db.Open(dbfile)
	if err != nil {
		t.Fatalf("Could not open db file: %s", dbfile)
	}
	defer conn.Close()

	job1result, err := conn.SelectJob(1, "j1")
	if err != nil {
		t.Fatalf("Could not read job1 result.")
	}

	if job1result.Var != "correct variable" {
		t.Errorf("Variable of job1 is unexpected: %s", job1result.Var)
	}
}

func TestRerun_JobCheck_Abnormal(t *testing.T) {
	defer util.SaveEvidence("rerun", "jobchk_abnormal")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	s := util.NewServant()
	s.UseConfig("jobchk_abnormal_servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("jobchk_abnormal_master.ini")
	rc, err := m.Rerun("1")
	if err != nil {
		t.Fatalf("Master start failed: %")
	}
	if rc != 0 {
		t.Fatalf("Master RC is not %d", 0)
	}

	if !isExecuted(1, "job1", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job1")
	}
	if !isExecuted(1, "job2", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job2")
	}
}

func TestRerun_JobCheck_Running(t *testing.T) {
	defer util.SaveEvidence("rerun", "jobchk_running")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	s := util.NewServant()
	s.UseConfig("jobchk_running_servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("jobchk_running_master.ini")
	rc, err := m.Rerun("1")
	if err != nil {
		t.Fatalf("Master start failed: %")
	}
	if rc == 0 {
		t.Fatalf("Master RC must not be %d", 0)
	}

	if isExecuted(1, "job1", masterLog) {
		t.Errorf("JOB [%s] must not be executed, but it was.", "job1")
	}
	if isExecuted(1, "job2", masterLog) {
		t.Errorf("JOB [%s] must not be executed, but it was.", "job2")
	}
}

func TestRerun_Node_Primary(t *testing.T) {
	defer util.SaveEvidence("rerun", "node_primary")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("node.ini")
	rc, err := m.Rerun("1")
	if err != nil {
		t.Fatalf("Master start failed: %")
	}
	if rc == 0 {
		t.Fatalf("Master RC must not be %d", 0)
	}
}

func TestRerun_Node_Secondary(t *testing.T) {
	defer util.SaveEvidence("rerun", "node_secondary")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("node.ini")
	rc, err := m.Rerun("2")
	if err != nil {
		t.Fatalf("Master start failed: %")
	}
	if rc != 0 {
		t.Fatalf("Master RC is not %d", 0)
	}

	if !isExecuted(2, "job1", masterLog) {
		t.Errorf("JOB [%s] must be executed, but it was not.", "job1")
	}
}

func TestRerun_ErrorCase_AlreadyNormalEnd(t *testing.T) {
	defer util.SaveEvidence("rerun", "errorcase_normalend")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("errorcase.ini")
	rc, err := m.Rerun("1")
	if err != nil {
		t.Fatalf("Master start failed: %")
	}
	if rc != 0 {
		t.Fatalf("Master RC is not %d", 0)
	}

	if !containsInFile("CTM029I", masterLog) {
		t.Errorf("Message Code [%s] must be output, but it was not.", "CTM029I")
	}
}

func TestRerun_ErrorCase_AlreadyWarnEnd(t *testing.T) {
	defer util.SaveEvidence("rerun", "errorcase_normalend")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("errorcase.ini")
	rc, err := m.Rerun("2")
	if err != nil {
		t.Fatalf("Master start failed: %")
	}
	if rc != 0 {
		t.Fatalf("Master RC is not %d", 0)
	}

	if !containsInFile("CTM029I", masterLog) {
		t.Errorf("Message Code [%s] must be output, but it was not.", "CTM029I")
	}
}

func TestRerun_ErrorCase_NotExecuted(t *testing.T) {
	defer util.SaveEvidence("rerun", "errorcase_normalend")
	util.InitCutoRoot()
	util.DeployTestData("rerun")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("errorcase.ini")
	rc, err := m.Rerun("3")
	if err != nil {
		t.Fatalf("Master start failed: %")
	}
	if rc == 0 {
		t.Fatalf("Master RC must not be %d", 0)
	}

	if !containsInFile("Network[id = 3] not found.", masterLog) {
		t.Error("Expected error message was not output.")
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

func containsInFile(str string, logfile string) bool {
	file, err := os.Open(logfile)
	if err != nil {
		return false
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		if strings.Contains(s.Text(), str) {
			return true
		}
	}

	return false
}
