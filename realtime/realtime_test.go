package realtime

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/unirita/cutotest/util"
)

func TestJSONOnly_Serial(t *testing.T) {
	defer util.SaveEvidence("realtime", "json_only", "serial")
	util.InitCutoRoot()
	util.DeployTestData("realtime")

	jobnetDir := filepath.Join(util.GetCutoRoot(), "bpmn")
	util.ClearDir(jobnetDir)

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	ts := httptest.NewServer(outputSerialFlow())
	defer ts.Close()

	r := util.NewRealtime()
	rc, err := r.Run(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if rc != 0 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 0)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertSuccessRealtimeOutput(t, res)
	assertNotRemainNetworkFile(t)
	assertIsNotNetworkEnds(t)

	waitProcessByPID(res.PID, 5)
	assertIsJob1to3Successed(t)
}

func TestJSONOnly_Parallel(t *testing.T) {
	defer util.SaveEvidence("realtime", "json_only", "parallel")
	util.InitCutoRoot()
	util.DeployTestData("realtime")

	jobnetDir := filepath.Join(util.GetCutoRoot(), "bpmn")
	util.ClearDir(jobnetDir)

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	ts := httptest.NewServer(outputParallelFlow())
	defer ts.Close()

	r := util.NewRealtime()
	rc, err := r.Run(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if rc != 0 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 0)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertSuccessRealtimeOutput(t, res)
	assertNotRemainNetworkFile(t)
	assertIsNotNetworkEnds(t)

	waitProcessByPID(res.PID, 5)
	assertIsJob1to3Successed(t)
}

func TestJSONOnly_WithJobDetail(t *testing.T) {
	defer util.SaveEvidence("realtime", "json_only", "withjobdetail")
	util.InitCutoRoot()
	util.DeployTestData("realtime")

	jobnetDir := filepath.Join(util.GetCutoRoot(), "bpmn")
	util.ClearDir(jobnetDir)

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	ts := httptest.NewServer(outputWithJobDetail("job1"))
	defer ts.Close()

	r := util.NewRealtime()
	rc, err := r.Run(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if rc != 0 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 0)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertSuccessRealtimeOutput(t, res)
	assertNotRemainNetworkFile(t)

	waitProcessByPID(res.PID, 1)
	assertTestjobSuccessed(t)
	assertJoblogExists(t, "job1")
}

func TestWithCSV_DefaultCSV(t *testing.T) {
	defer util.SaveEvidence("realtime", "with_csv", "defaultcsv")
	util.InitCutoRoot()
	util.DeployTestData("realtime")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	ts := httptest.NewServer(outputTestJobOnly())
	defer ts.Close()

	r := util.NewRealtime()
	rc, err := r.Run(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if rc != 0 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 0)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertSuccessRealtimeOutput(t, res)
	assertNotRemainNetworkFile(t)

	waitProcessByPID(res.PID, 1)
	assertTestjobSuccessed(t)
	assertJoblogExists(t, "job1")
}

func TestWithCSV_NamedCSV(t *testing.T) {
	defer util.SaveEvidence("realtime", "with_csv", "namedcsv")
	util.InitCutoRoot()
	util.DeployTestData("realtime")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	ts := httptest.NewServer(outputTestJobOnly())
	defer ts.Close()

	r := util.NewRealtime()
	rc, err := r.Run("-n", "test", ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if rc != 0 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 0)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertSuccessNamedRealtimeOutput(t, res, "test")
	assertNotRemainNetworkFile(t)

	waitProcessByPID(res.PID, 1)
	assertTestjobSuccessed(t)
	assertJoblogExists(t, "job2")
}

func TestWithCSV_Priority(t *testing.T) {
	defer util.SaveEvidence("realtime", "with_csv", "priority")
	util.InitCutoRoot()
	util.DeployTestData("realtime")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	ts := httptest.NewServer(outputWithJobDetail("job3"))
	defer ts.Close()

	r := util.NewRealtime()
	rc, err := r.Run("-n", "test", ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if rc != 0 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 0)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertSuccessNamedRealtimeOutput(t, res, "test")
	assertNotRemainNetworkFile(t)

	waitProcessByPID(res.PID, 5)
	assertTestjobSuccessed(t)
	assertJoblogExists(t, "job3")
}
