package oncontainer

import (
	"testing"

	"github.com/unirita/cutotest/util/container"
)

func TestMain(m *testing.M) {
	os.Exit(realTestMain(m))
}

func realTestMain(m *testing.M) int {
	c := container.New("cuto/servant", "test/container")
	c.Start()
	defer c.Terminate()

	defer util.SaveEvidence("oncontainer")
	util.InitCutoRoot()
	util.DeployTestData("oncontainer")

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	return m.Run()
}

func TestOnContainerJob_Joblog(t *testing.T) {
	m := util.NewMaster()
	m.UseConfig("master.ini")
	rc, err = m.Run("joblog")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
}
