package realtime

import (
	"net/http/httptest"
	"os"
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
	assertIsNotNetworkEnd(t)

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
	assertIsNotNetworkEnd(t)

	waitProcessByPID(res.PID, 5)
	assertIsJob1to3Successed(t)
}

func TestJSONOnly_WithJobDetail(t *testing.T) {
	defer util.SaveEvidence("realtime", "json_only", "with_job_detail")
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
	defer util.SaveEvidence("realtime", "with_csv", "default_csv")
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
	defer util.SaveEvidence("realtime", "with_csv", "named_csv")
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

func TestAbnormal_InvalidURL(t *testing.T) {
	defer util.SaveEvidence("realtime", "abnormal", "invalid_url")
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
	rc, err := r.Run("wrongschema://wrong.co.jp:wrongport")
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if rc != 1 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 1)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertFailedRealtimeOutput(t, res, 2)
	assertNotRemainNetworkFile(t)
	assertIsNotNetworkStart(t)
}

func TestAbnormal_ServerNotExists(t *testing.T) {
	defer util.SaveEvidence("realtime", "abnormal", "server_not_exists")
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
	rc, err := r.Run("http://notexists")
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if rc != 1 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 1)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertFailedRealtimeOutput(t, res, 2)
	assertNotRemainNetworkFile(t)
	assertIsNotNetworkStart(t)
}

func TestAbnormal_FlowError(t *testing.T) {
	defer util.SaveEvidence("realtime", "abnormal", "flow_error")
	util.InitCutoRoot()
	util.DeployTestData("realtime")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	ts := httptest.NewServer(outputErrorFlow())
	defer ts.Close()

	r := util.NewRealtime()
	rc, err := r.Run(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if rc != 1 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 1)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertFailedRealtimeOutput(t, res, 2)
	assertNotRemainNetworkFile(t)
	assertIsNotNetworkStart(t)
}

func TestAbnormal_ConfigNotFound(t *testing.T) {
	defer util.SaveEvidence("realtime", "abnormal", "config_not_found")
	util.InitCutoRoot()
	util.DeployTestData("realtime")

	masterConfig := filepath.Join(util.GetCutoRoot(), "bin", "master.ini")
	os.Remove(masterConfig)

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
	if rc != 1 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 1)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertFailedRealtimeOutput(t, res, 2)
	assertNotRemainNetworkFile(t)
	assertIsNotNetworkStart(t)
}

func TestAbnormal_ErrorInMaster(t *testing.T) {
	defer util.SaveEvidence("realtime", "abnormal", "error_in_master")
	util.InitCutoRoot()
	util.DeployTestData("realtime")

	masterConfig := filepath.Join(util.GetCutoRoot(), "bin", "master.ini")
	masterErrorConfig := filepath.Join(util.GetCutoRoot(), "bin", "master_error.ini")
	os.Remove(masterConfig)
	os.Rename(masterErrorConfig, masterConfig)

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
	if rc != 1 {
		t.Log(r.Stdout)
		t.Fatalf("rc => %d, want %d", rc, 1)
	}

	res := parseRealtimeResult(t, r.Stdout)
	assertFailedRealtimeOutput(t, res, 1)
	assertNotRemainNetworkFile(t)
	assertIsNotNetworkStart(t)
}
