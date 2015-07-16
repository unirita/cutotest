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
	defer util.SaveEvidence("secondary", "noretry")
	util.InitCutoRoot()
	util.DeployTestData("secondary")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("master_noretry.ini")
	_, err := m.Run("nosecondary")
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

func TestRetry_NoSecondary(t *testing.T) {
	defer util.SaveEvidence("secondary", "nosecondary")
	util.InitCutoRoot()
	util.DeployTestData("secondary")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("master_withretry.ini")
	_, err := m.Run("nosecondary")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if !util.ContainsInFile(masterLog, "CTM025W") {
		t.Errorf("Job did not end abnormally.")
	}
	if !util.ContainsInFile(masterLog, "CTM027W") {
		t.Errorf("Job was not retried.")
	}
	if util.ContainsInFile(masterLog, "CTM028W") {
		t.Errorf("Executed job at secondary servant unexpectedly.")
	}
}

func TestRetry_WithSecondary(t *testing.T) {
	defer util.SaveEvidence("secondary", "withsecondary")
	util.InitCutoRoot()
	util.DeployTestData("secondary")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("master_withretry.ini")
	_, err := m.Run("withsecondary")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if !util.ContainsInFile(masterLog, "CTM024I") {
		t.Errorf("Job did not end normally.")
	}
	if !util.ContainsInFile(masterLog, "CTM027W") {
		t.Errorf("Job was not retried.")
	}
	if !util.ContainsInFile(masterLog, "CTM028W") {
		t.Errorf("Job was not executed at secondary servant.")
	}
}

func TestRetry_WithErrorSecondary(t *testing.T) {
	defer util.SaveEvidence("secondary", "witherrorsecondary")
	util.InitCutoRoot()
	util.DeployTestData("secondary")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("master_withretry.ini")
	_, err := m.Run("witherrorsecondary")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if !util.ContainsInFile(masterLog, "CTM025W") {
		t.Errorf("Job did not end abnormally.")
	}
	if util.CountInFile(masterLog, "CTM027W") != 2 {
		t.Errorf("Job was not retried twice.")
	}
	if !util.ContainsInFile(masterLog, "CTM028W") {
		t.Errorf("Job was not executed at secondary servant.")
	}
}

func TestRetry_WithSecondary_NoPrimaryError(t *testing.T) {
	defer util.SaveEvidence("secondary", "noprimaryerror")
	util.InitCutoRoot()
	util.DeployTestData("secondary")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("master_withretry.ini")
	_, err := m.Run("noprimaryerror")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if !util.ContainsInFile(masterLog, "CTM024I") {
		t.Errorf("Job did not end normally.")
	}
	if util.ContainsInFile(masterLog, "CTM027W") {
		t.Errorf("Retried job execution unexpectedly.")
	}
	if util.ContainsInFile(masterLog, "CTM028W") {
		t.Errorf("Executed job at secondary servant unexpectedly.")
	}
}

func TestMultiRetry(t *testing.T) {
	defer util.SaveEvidence("secondary", "multiretry")
	util.InitCutoRoot()
	util.DeployTestData("secondary")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("master_multiretry.ini")
	_, err := m.Run("witherrorsecondary")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if !util.ContainsInFile(masterLog, "CTM025W") {
		t.Errorf("Job did not end abnormally.")
	}
	if util.CountInFile(masterLog, "CTM027W") != 4 {
		t.Errorf("Job was not retried twice.")
	}
	if !util.ContainsInFile(masterLog, "CTM028W") {
		t.Errorf("Job was not executed at secondary servant.")
	}
}

func TestRuntimeError(t *testing.T) {
	defer util.SaveEvidence("secondary", "runtimeerror")
	util.InitCutoRoot()
	util.DeployTestData("secondary")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("master_withretry.ini")
	_, err := m.Run("runtimeerror")
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
