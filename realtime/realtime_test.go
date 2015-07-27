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

	proc, err := os.FindProcess(res.PID)
	if err == nil {
		proc.Wait()
	}
	assertIsJob1to3Successed(t)
}
