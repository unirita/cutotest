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
	m.Run("inst_test")
}
