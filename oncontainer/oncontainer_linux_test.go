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

	servant := util.NewServant()
	servant.UseConfig("servant.ini")
	if err := servant.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer servant.Kill()

	return m.Run()
}

func TestOnContainerJob_Joblog(t *testing.T) {
	master := util.NewMaster()
	master.UseConfig("master.ini")
	rc, err = master.Run("joblog")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
}
