package serialexec

import (
	"path/filepath"
	"testing"

	"github.com/unirita/cutotest/util"
	"github.com/unirita/cutotest/util/db"
)

func TestSerialExecution(t *testing.T) {
	util.InitCutoRoot()
	util.DeployTestData("serialexec")
	util.ComplementConfig("master.ini")
	util.ComplementConfig("servant.ini")

	s := util.NewServant()
	s.SetConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.SetConfig("master.ini")

	rc, err := m.SyntaxCheck("Serial")
	if err != nil {
		t.Fatalf("Master bpmn syntax check failed: %s", err)
	}
	if rc != 0 {
		t.Errorf("Master RC[%d] is not 0.", rc)
	}

	rc, err = m.Run("Serial")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if rc != 0 {
		t.Errorf("Master RC[%d] is not 0.", rc)
	}
	if util.HasLogError(m.ConfigPath) {
		t.Errorf("There is error log in [%s]", m.ConfigPath)
	}

	conn, err := db.Open(util.GetDBDirPath())
	if err != nil {
		t.Fatalf("DB file open failed: %v", err)
	}
	defer conn.Close()

	network, err := conn.SelectJobNetwork(1)
	if err != nil {
		t.Fatalf("Can't read network record: %v", err)
	}
	if network.Name != "Serial" {
		t.Errorf("Unexpected JOBNETWORK.JOBNETWORK[%s]", network.Name)
	}
	if network.Start == "" {
		t.Error("JOBNETWORK.STARTDATE is empty.")
	}
	if network.End == "" {
		t.Error("JOBNETWORK.STARTDATE is empty.")
	}
	if network.Status != 1 {
		t.Errorf("Unexpected JOBNETWORK.STATUS[%d]", network.Status)
	}

	if util.IsWindows() {
		testDBForWindows(t, conn)
	} else {
		testDBForLinux(t, conn)
	}

	jobname := "job1"
	joblogs := util.FindJoblog("joblog", 1, jobname)
	if len(joblogs) != 1 {
		t.Fatalf("%s has no joblog or multi joblogs.", jobname)
	}
	if !util.IsExistNoEmptyFile(joblogs[0]) {
		t.Errorf("Joblog[%s] is not exist or empty.", joblogs[0])
	}

	jobname = "job2"
	joblogs = util.FindJoblog("joblog", 1, jobname)
	if len(joblogs) != 1 {
		t.Fatalf("%s has no joblog or multi joblogs.", jobname)
	}
	if !util.IsExistNoEmptyFile(joblogs[0]) {
		t.Errorf("Joblog[%s] is not exist or empty.", joblogs[0])
	}

	// job3-6 is not exists on linux test data.
	if util.IsWindows() {
		jobname = "job3"
		joblogs = util.FindJoblog("joblog", 1, jobname)
		if len(joblogs) != 1 {
			t.Fatalf("%s has no joblog or multi joblogs.", jobname)
		}
		if !util.IsExistNoEmptyFile(joblogs[0]) {
			t.Errorf("Joblog[%s] is not exist or empty.", joblogs[0])
		}

		jobname = "job4"
		joblogs = util.FindJoblog("joblog", 1, jobname)
		if len(joblogs) != 1 {
			t.Fatalf("%s has no joblog or multi joblogs.", jobname)
		}
		if !util.IsExistNoEmptyFile(joblogs[0]) {
			t.Errorf("Joblog[%s] is not exist or empty.", joblogs[0])
		}

		jobname = "job5"
		joblogs = util.FindJoblog("joblog", 1, jobname)
		if len(joblogs) != 1 {
			t.Fatalf("%s has no joblog or multi joblogs.", jobname)
		}
		if !util.IsExistNoEmptyFile(joblogs[0]) {
			t.Errorf("Joblog[%s] is not exist or empty.", joblogs[0])
		}

		jobname = "job6"
		joblogs = util.FindJoblog("joblog", 1, jobname)
		if len(joblogs) != 1 {
			t.Fatalf("%s has no joblog or multi joblogs.", jobname)
		}
		if !util.IsExistNoEmptyFile(joblogs[0]) {
			t.Errorf("Joblog[%s] is not exist or empty.", joblogs[0])
		}
	}

	jobname = "param"
	joblogs = util.FindJoblog("joblog", 1, jobname)
	if len(joblogs) != 1 {
		t.Fatalf("%s has no joblog or multi joblogs.", jobname)
	}
	if !util.ContainsInFile(joblogs[0], "param1=test1") {
		t.Errorf("Joblog[%s] is not expect output.", joblogs[0])
	}
	if !util.ContainsInFile(joblogs[0], "param2=test2") {
		t.Errorf("Joblog[%s] is not expect output.", joblogs[0])
	}

	jobname = "env"
	joblogs = util.FindJoblog("joblog", 1, jobname)
	if len(joblogs) != 1 {
		t.Fatalf("%s has no joblog or multi joblogs.", jobname)
	}
	if !util.ContainsInFile(joblogs[0], "ENV1=TEST1") {
		t.Errorf("Joblog[%s] is not expect output.", joblogs[0])
	}
	if !util.ContainsInFile(joblogs[0], "ENV2=TEST2") {
		t.Errorf("Joblog[%s] is not expect output.", joblogs[0])
	}

	jobname = "work"
	joblogs = util.FindJoblog("joblog", 1, jobname)
	workPath := filepath.Join(util.GetCutoRoot(), "jobscript", "otherpath")
	if len(joblogs) != 1 {
		t.Fatalf("%s has no joblog or multi joblogs.", jobname)
	}
	if !util.ContainsInFile(joblogs[0], workPath) {
		t.Errorf("Joblog[%s] is not expect output.", joblogs[0])
	}

	jobname = "path"
	joblogs = util.FindJoblog("joblog", 1, jobname)
	var scriptPath string
	if util.IsWindows() {
		scriptPath = filepath.Join(workPath, "path.bat")
	} else {
		scriptPath = filepath.Join(workPath, "path.sh")
	}
	if len(joblogs) != 1 {
		t.Fatalf("%s has no joblog or multi joblogs.", jobname)
	}
	if !util.ContainsInFile(joblogs[0], scriptPath) {
		t.Errorf("Joblog[%s] is not expect output.", joblogs[0])
	}
}

func testDBForWindows(t *testing.T, conn *db.Connection) {
	jid := "job1"
	job, err := conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "job1.js" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 0 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "job2"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "job2.vbs" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 0 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "job3"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "job3.jar" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 0 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "job4"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "job4.ps1" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 0 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "job5"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "job5.bat" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 0 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "job6"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "job6.exe" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 0 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "param"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "param.bat" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "env"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "env.bat" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 2 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "work"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "work.bat" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 3 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "path"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "path.bat" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 4 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}
}

func testDBForLinux(t *testing.T, conn *db.Connection) {
	jid := "job1"
	job, err := conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "job1" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 0 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "job2"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "job2.sh" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 0 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "param"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "param.sh" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "env"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "env.sh" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 2 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "work"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "work.sh" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 3 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}

	jid = "path"
	job, err = conn.SelectJob(1, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != "path.sh" {
		t.Errorf("Unexpected JOB.JOBNAME[%s] on JOBID=%s", job.Name, jid)
	}
	if job.Start == "" {
		t.Errorf("JOB.STARTDATE is empty on JOBID=%s", jid)
	}
	if job.End == "" {
		t.Errorf("JOB.ENDDATE is empty on JOBID=%s", jid)
	}
	if job.Status != 1 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.Status, jid)
	}
	if job.RC != 4 {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}
}
