package serialexec

import (
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

	assertDB(t, 1)
	assertJoblog(t, 1)
}

func assertNetworkRecord(t *testing.T, conn *db.Connection, nid int) {
	network, err := conn.SelectJobNetwork(nid)
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
}

func assertJobRecord(t *testing.T, conn *db.Connection, nid int, jid string, jobname string, rc int) {
	job, err := conn.SelectJob(nid, jid)
	if err != nil {
		t.Fatalf("Can't read %s record: %v", jid, err)
	}
	if job.Name != jobname {
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
	if job.RC != rc {
		t.Errorf("Unexpected JOB.STATUS[%d] on JOBID=%s", job.RC, jid)
	}
}

func assertHasNoEmptyJoblog(t *testing.T, nid int, jobname string) {
	joblogs := util.FindJoblog("joblog", nid, jobname)
	if len(joblogs) != 1 {
		t.Fatalf("%s has no joblog or multi joblogs.", jobname)
	}
	if !util.IsExistNoEmptyFile(joblogs[0]) {
		t.Errorf("Joblog[%s] is not exist or empty.", joblogs[0])
	}
}

func assertJoblogContainsStr(t *testing.T, nid int, jobname string, substrs ...string) {
	joblogs := util.FindJoblog("joblog", nid, jobname)
	if len(joblogs) != 1 {
		t.Fatalf("%s has no joblog or multi joblogs.", jobname)
	}
	for _, substr := range substrs {
		if !util.ContainsInFile(joblogs[0], substr) {
			t.Errorf("Joblog[%s] is not expect output.", joblogs[0])
		}
	}
}
