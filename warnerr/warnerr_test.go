package warnerr

import (
	"fmt"
	"testing"

	"github.com/unirita/cutotest/util"
	"github.com/unirita/cutotest/util/db"
)

func TestWarnings(t *testing.T) {
	defer util.SaveEvidence("warnerr_warn")
	util.InitCutoRoot()
	util.DeployTestData("warnerr")
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
	rc, err := m.Run("Warn")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if rc != 0 {
		t.Logf("Master stderr: %s", m.Stderr)
		t.Errorf("Master RC[%d] is not 0.", rc)
	}
	logPath := util.GetLogPath("master.log")
	if util.HasLogError(logPath) {
		t.Errorf("There is error log in [%s]", logPath)
	}

	conn, err := db.Open(util.GetDBDirPath())
	if err != nil {
		t.Fatalf("DB file open failed: %v", err)
	}
	defer conn.Close()

	assertNetworkRecord(t, conn, 1, "Warn", 2)
	assertJobStatus(t, conn, 1, "nrcmax", 1)
	assertJobStatus(t, conn, 1, "wrcmin", 2)
	assertJobStatus(t, conn, 1, "wrcmax", 2)
	assertJobStatus(t, conn, 1, "wstdout", 2)
	assertJobStatus(t, conn, 1, "wstderr", 2)

	mustExecJobs, err := conn.SelectJobsByCond("ID=1 AND JOBNAME='mustexec'")
	if err != nil {
		t.Fatalf("Unexpected DB error occured: %s", err)
	}
	assertJobIDExists(t, mustExecJobs, "$C$2")
	assertJobIDExists(t, mustExecJobs, "$C$3")
	assertJobIDExists(t, mustExecJobs, "$C$4")
	assertJobIDExists(t, mustExecJobs, "$D$1")
	assertJobIDExists(t, mustExecJobs, "$F$1")
}

func TestErrors(t *testing.T) {
	defer util.SaveEvidence("warnerr_err")
	util.InitCutoRoot()
	util.DeployTestData("warnerr")
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
	rc, err := m.Run("Error")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if rc != 1 {
		t.Logf("Master stderr: %s", m.Stderr)
		t.Errorf("Master RC[%d] is not 1.", rc)
	}
	logPath := util.GetLogPath("master.log")
	if !util.HasLogError(logPath) {
		t.Errorf("There is no error log in [%s]", logPath)
	}

	conn, err := db.Open(util.GetDBDirPath())
	if err != nil {
		t.Fatalf("DB file open failed: %v", err)
	}
	defer conn.Close()

	assertNetworkRecord(t, conn, 1, "Error", 9)
	assertJobStatus(t, conn, 1, "ercmin", 9)
	assertJobStatus(t, conn, 1, "estdout", 9)
	assertJobStatus(t, conn, 1, "estderr", 9)
	assertJobStatus(t, conn, 1, "aftererror", 1)

	noExecJobs, err := conn.SelectJobsByCond("ID=1 AND JOBNAME='noexec'")
	if err != nil {
		t.Fatalf("Unexpected DB error occured: %s", err)
	}
	if len(noExecJobs) != 0 {
		t.Errorf("Unexpected Job was executed.")
	}
}

func assertNetworkRecord(t *testing.T, conn *db.Connection, nid int, name string, status int) {
	network, err := conn.SelectJobNetwork(nid)
	if err != nil {
		t.Fatalf("Can't read network record: %v", err)
	}
	if network.Name != name {
		t.Errorf("Unexpected JOBNETWORK.JOBNETWORK[%s]", network.Name)
	}
	if network.Status != status {
		t.Errorf("Unexpected JOBNETWORK.STATUS[%d]", network.Status)
	}
}

func assertJobStatus(t *testing.T, conn *db.Connection, nid int, jobname string, status int) {
	cond := fmt.Sprintf("ID=%d AND JOBNAME='%s'", nid, jobname)
	jobs, err := conn.SelectJobsByCond(cond)
	if err != nil {
		t.Fatalf("Unexpected DB error occured: %s", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("JOBNAME[%s] must be only one.", jobname)
	}
	if jobs[0].Status != status {
		t.Errorf("Unexpected JOB.STATUS[%d]", jobs[0].Status)
	}
}

func assertJobIDExists(t *testing.T, jobs []*db.Job, jid string) {
	exists := false
	for _, job := range jobs {
		if job.JID == jid {
			exists = true
			break
		}
	}
	if !exists {
		t.Errorf("JOB ID[%s] not found.", jid)
	}
}
