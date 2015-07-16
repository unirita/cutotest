package secondary

import (
	"testing"

	"github.com/unirita/cutotest/util"
)

var masterLog string = util.GetLogPath("master.log")

func TestLegacy(t *testing.T) {
	defer util.SaveEvidence("secondary", "legacy")
	util.InitCutoRoot()
	util.DeployTestData("secondary")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("master_legacy.ini")
	_, err := m.Run("legacy")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if !util.ContainsInFile(masterLog, "CTM025W") {
		t.Errorf("Job did not end abnormally.")
	}
	if util.ContainsInFile(masterLog, "CTM027W") {
		t.Errorf("Retried job execution unexpectedly.")
	}
	if util.ContainsInFile(masterLog, "CTM028W") {
		t.Errorf("Executed job at secondary servant unexpectedly.")
	}
}

func TestNoRetry(t *testing.T) {

}

func TestRetry_NoSecondary(t *testing.T) {

}

func TestRetry_WithSecondary(t *testing.T) {

}

func TestRetry_WithErrorSecondary(t *testing.T) {

}

func TestRetry_WithSecondary_NoPrimaryError(t *testing.T) {

}
