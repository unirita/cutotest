package timezone

import (
	"fmt"
	"testing"
	"time"

	"github.com/unirita/cutotest/util/container"
)

func TestMain(m *testing.M) {
	os.Exit(realTestMain(m))
}

const imageName = "cuto/servant"

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

	return m.Run()
}

func TestNetworkTimestamp(t *testing.T) {

}
