package regression

import (
	"testing"

	"github.com/unirita/cutotest/util"
)

func TestRegression(t *testing.T) {
	util.InitCutoRoot()
	util.DeployTestData("regression")
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
	rc, err := m.Run("inst_test")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if rc != 0 {
		t.Logf("Master stderr: %s", m.Stderr)
		t.Errorf("Master RC[%d] is not 0.", rc)
	}
	if util.HasLogError(s.ConfigPath) {
		t.Errorf("There is error log in [%s]", s.ConfigPath)
	}
}
