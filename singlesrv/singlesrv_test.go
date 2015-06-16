package singlesrv

import (
	"fmt"
	"testing"

	"github.com/unirita/cutotest/util"
	"github.com/unirita/cutotest/util/db"
)

func Test255Job(t *testing.T) {
	util.InitCutoRoot()
	util.DeployTestData("singlesrv")
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

	rc, err := m.SyntaxCheck("255JobSingle")
	if err != nil {
		t.Fatalf("Master bpmn syntax check failed: %s", err)
	}
	if rc != 0 {
		t.Errorf("Master RC[%d] is not 0.", rc)
	}

	rc, err = m.Run("255JobSingle")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if rc != 0 {
		t.Errorf("Master RC[%d] is not 0.", rc)
	}
	if util.HasLogError(m.ConfigPath) {
		t.Errorf("There is error log in [%s]", m.ConfigPath)
	}

	// Check database.
	conn, err := db.Open(util.GetDBDirPath())
	if err != nil {
		t.Fatalf("DB file open failed: %v", err)
	}
	defer conn.Close()

	network, err := conn.SelectJobNetwork(1)
	if err != nil {
		t.Fatalf("Can't read network record: %v", err)
	}
	if network.Name != "255JobSingle" {
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

	count, err := conn.CountJobs(1)
	if err != nil {
		t.Fatalf("Unexpected DB error occured: %v", err)
	}
	if count != 511 {
		t.Errorf("Number of Job record[%d] is not expected value.", count)
	}

	subquery := "SELECT STARTDATE FROM JOB WHERE ID = 1 AND JOBNAME = 'afterbranch' LIMIT 1"
	cond := fmt.Sprintf("ID = 1 AND ENDDATE > (%s)", subquery)
	jobs, err := conn.SelectJobsByCond(cond)
	if err != nil {
		t.Fatalf("Unexpected DB error occured: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatal("Job which executed after branch must be only one.")
	}
	if jobs[0].Name != "afterbranch" {
		t.Errorf("Unexpected Job[%s] executed after branch.", jobs[0].JID)
	}
}
