package timezone

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/unirita/cutotest/util"
	"github.com/unirita/cutotest/util/container"
)

const imageName = "cuto/servant"

func TestMain(m *testing.M) {
	os.Exit(realTestMain(m))
}

func realTestMain(m *testing.M) int {
	if _, offset := time.Now().Zone(); offset/3600 != 9 {
		fmt.Println("Timezone must be +0900.")
		return 1
	}
	cntUTC := container.New(imageName, "cntutc")
	cntUTC.SetTimezone("UTC")
	cntUTC.Start()
	defer cntUTC.Terminate()

	cntMST := container.New(imageName, "cntmst")
	cntMST.SetTimezone("MST")
	cntMST.Start()
	defer cntMST.Terminate()

	util.InitCutoRoot()
	util.DeployTestData("timezone")
	complementJobDetail(cntUTC.IPAddress(), cntMST.IPAddress())

	s := util.NewServant()
	s.UseConfig("servant.ini")
	if err := s.Start(); err != nil {
		fmt.Printf("Servant start failed: %s\n", err)
		return 1
	}
	defer s.Kill()

	m := util.NewMaster()
	m.UseConfig("master.ini")
	_, err = m.Run("timezone")
	if err != nil {
		fmt.Printf("Master run failed: %s", err)
		return 1
	}

	return m.Run()
}

func TestNetworkTimestamp(t *testing.T) {

}
