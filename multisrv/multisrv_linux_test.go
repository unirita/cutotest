package multisrv

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/unirita/cutotest/util"
	"github.com/unirita/cutotest/util/container"
	"github.com/unirita/cutotest/util/db"
)

type hostParams struct {
	Containers [25]string
}

const imageName = "cuto/servant"

func Test255Job(t *testing.T) {
	defer util.SaveEvidence("multisrv")
	util.InitCutoRoot()
	util.DeployTestData("multisrv")
	util.ComplementConfig("master.ini")
	util.ComplementConfig("servant.ini")

	s := util.NewServant()
	s.SetConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	params := new(hostParams)
	for i := 0; i < 25; i++ {
		name := fmt.Sprintf("TestContainer%d", i+1)
		cont := container.New(imageName, name)
		err := cont.Start()
		if err != nil {
			t.Fatalf("Failed to start container[%s]", name)
		}
		defer cont.Terminate()
		params.Containers[i], err = cont.IPAddress()
		if err != nil {
			t.Fatalf("Failed to get container[%s] IPAddress", name)
		}
	}

	if err := complement255JobDetail(params); err != nil {
		t.Fatalf("Failed to complete job detail csv file: %s", err)
	}

	m := util.NewMaster()
	m.SetConfig("master.ini")

	rc, err := m.SyntaxCheck("255Job")
	if err != nil {
		t.Fatalf("Master bpmn syntax check failed: %s", err)
	}
	if rc != 0 {
		t.Errorf("Master RC[%d] is not 0.", rc)
	}

	rc, err = m.Run("255Job")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if rc != 0 {
		t.Errorf("Master RC[%d] is not 0.", rc)
	}

	// Check master log.
	logPath := util.GetLogPath("master.log")
	if util.HasLogError(logPath) {
		t.Errorf("There is error log in [%s]", logPath)
	}
	if !util.IsPatternExistInFile(logPath, params.Containers[0]) {
		t.Errorf("Node name was not recorded in [%s]", logPath)
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
	if network.Name != "255Job" {
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

	subquery := "SELECT STARTDATE FROM JOB WHERE ID = 1 AND JOBNAME = 'usevar' LIMIT 1"
	cond := fmt.Sprintf("ID = 1 AND ENDDATE > (%s)", subquery)
	jobs, err := conn.SelectJobsByCond(cond)
	if err != nil {
		t.Fatalf("Unexpected DB error occured: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatal("Job which executed after branch must be only one.")
	}
	if jobs[0].Name != "usevar" {
		t.Errorf("Unexpected Job[%s] executed after branch.", jobs[0].JID)
	}

	// Check joblog.
	joblogs := util.FindJoblog("joblog", 1, "usevar")
	if len(joblogs) != 1 {
		t.Fatalf("usevar has no joblog or multi joblogs.")
	}
	joblog := joblogs[0]

	if !util.IsPatternExistInFile(joblog, "^FLOWID=1") {
		t.Errorf("Job[usevar] did not output correct FLOWID.")
	}
	if !util.IsPatternExistInFile(joblog, "^FLOWSD=.+") {
		t.Errorf("Job[usevar] did not output correct FLOWSD.")
	}
	if !util.IsPatternExistInFile(joblog, "^SSROOT=.+") {
		t.Errorf("Job[usevar] did not output correct SSROOT.")
	}
	if !util.IsPatternExistInFile(joblog, "^MEPATH=.+") {
		t.Errorf("Job[usevar] did not output correct MEPATH.")
	}
	if !util.IsPatternExistInFile(joblog, "^SEPATH=.+") {
		t.Errorf("Job[usevar] did not output correct SEPATH.")
	}
	if !util.IsPatternExistInFile(joblog, "^JOUT=cap0011out") {
		t.Errorf("Job[usevar] did not output correct JOUT.")
	}
	if !util.IsPatternExistInFile(joblog, "^JRC=21") {
		t.Errorf("Job[usevar] did not output correct JRC.")
	}
	if !util.IsPatternExistInFile(joblog, "^JSD=.+") {
		t.Errorf("Job[usevar] did not output correct JSD.")
	}
	if !util.IsPatternExistInFile(joblog, "^JED=.+") {
		t.Errorf("Job[usevar] did not output correct JED.")
	}
}

func complement255JobDetail(params *hostParams) error {
	path := filepath.Join(util.GetCutoRoot(), "bpmn", "255Job.csv")
	tpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tpl.Execute(file, params); err != nil {
		return err
	}

	return nil
}
