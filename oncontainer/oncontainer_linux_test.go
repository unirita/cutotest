package oncontainer

import (
	"fmt"
	"os"
	"testing"

	"github.com/unirita/cutotest/util"
	"github.com/unirita/cutotest/util/container"
	"github.com/unirita/cutotest/util/db"
)

func TestMain(m *testing.M) {
	os.Exit(realTestMain(m))
}

func realTestMain(m *testing.M) int {
	c := container.New("cuto/servant", "test_oncontainer")
	if err := c.Start(); err != nil {
		fmt.Printf("Could not run container: %s\n", err)
		return 1
	}
	defer c.Terminate()

	return m.Run()
}

func TestOnContainerJob_GetJoblog(t *testing.T) {
	defer util.SaveEvidence("oncontainer", "getjoblog")
	util.InitCutoRoot()
	util.DeployTestData("oncontainer")

	servant := util.NewServant()
	servant.UseConfig("servant.ini")
	if err := servant.Start(); err != nil {
		t.Fatalf("Servant start failed: %s\n", err)
	}
	defer servant.Kill()

	master := util.NewMaster()
	master.UseConfig("master.ini")
	rc, err := master.Run("joblog")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if rc != 0 {
		t.Fatalf("Master RC => %d, wants %d", rc, 0)
	}

	joblogs := util.FindJoblog("joblog", 2, "varout")
	if len(joblogs) != 1 {
		t.Fatalf("Number of joblog => %d, wants %d", len(joblogs), 1)
	}

	if !util.ContainsInFile(joblogs[0], "testparam") {
		t.Error("Joblog was not output correctly.")
	}
}

func TestOnContainerJob_RC(t *testing.T) {
	defer util.SaveEvidence("oncontainer", "remote")
	util.InitCutoRoot()
	util.DeployTestData("oncontainer")

	servant := util.NewServant()
	servant.UseConfig("servant.ini")
	if err := servant.Start(); err != nil {
		t.Fatalf("Servant start failed: %s\n", err)
	}
	defer servant.Kill()

	master := util.NewMaster()
	master.UseConfig("master.ini")
	rc, err := master.Run("rc")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if rc != 0 {
		t.Fatalf("Master RC => %d, wants %d", rc, 0)
	}

	conn, err := db.Open(util.GetDBDirPath())
	if err != nil {
		t.Fatalf("Database open failed: %s", err)
	}
	jobRecord, err := conn.SelectJob(2, "job1")
	if err != nil {
		t.Fatalf("Could not select job record: %s", err)
	}

	if jobRecord.RC != 123 {
		t.Errorf("RC of job => %d, wants %d", jobRecord.RC, 123)
	}
}

func TestOnContainerJob_Remote(t *testing.T) {
	defer util.SaveEvidence("oncontainer", "remote")
	util.InitCutoRoot()
	util.DeployTestData("oncontainer")

	servant := util.NewServant()
	servant.UseConfig("servant.ini")
	if err := servant.Start(); err != nil {
		t.Fatalf("Servant start failed: %s\n", err)
	}
	defer servant.Kill()

	master := util.NewMaster()
	master.UseConfig("master.ini")
	rc, err := master.Run("remote")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if rc != 0 {
		t.Fatalf("Master RC => %d, wants %d", rc, 0)
	}

	conn, err := db.Open(util.GetDBDirPath())
	if err != nil {
		t.Fatalf("Database open failed: %s", err)
	}
	jobRecord, err := conn.SelectJob(2, "job1")
	if err != nil {
		t.Fatalf("Could not select job record: %s", err)
	}

	if jobRecord.RC != 234 {
		t.Errorf("RC of job => %d, wants %d", jobRecord.RC, 234)
	}
}
