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

	hostUTC, err := cntUTC.IPAddress()
	if err != nil {
		fmt.Printf("Could not get container hostname.")
		return 1
	}
	hostMST, err := cntMST.IPAddress()
	if err != nil {
		fmt.Printf("Could not get container hostname.")
		return 1
	}
	complementJobDetail(hostUTC, hostMST)
	complementDB()

	servant := util.NewServant()
	servant.UseConfig("servant.ini")
	if err := servant.Start(); err != nil {
		fmt.Printf("Servant start failed: %s\n", err)
		return 1
	}
	defer servant.Kill()

	master := util.NewMaster()
	master.UseConfig("master.ini")
	if _, err := master.Run("timezone"); err != nil {
		fmt.Printf("Master run failed: %s", err)
		return 1
	}

	return m.Run()
}
