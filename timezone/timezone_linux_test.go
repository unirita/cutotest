package timezone

import (
	"testing"
	"time"

	"github.com/unirita/cutotest/util"
)

const dateFormat = "20060102"

func runShow(t *testing.T, isUTC bool) string {
	now := time.Now()
	show := util.NewShow()
	show.UseConfig("master.ini")
	show.AddFrom(now.Format(dateFormat))
	show.AddTo(now.Format(dateFormat))
	if isUTC {
		show.AddUTCOption()
	}
	rc, err := show.Run()
	if err != nil {
		t.Fatalf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Fatalf("Show RC[%d] is not 0.", rc)
	}

	return show.Stdout
}

func extractNetworkResults(jobnets []*OutputJobNet) (*OutputJobNet, *OutputJobNet, *OutputJobNet) {
	var testnet, todaynet, yesterdaynet *OutputJobNet
	for _, jobnet := range jobnets {
		switch jobnet.Jobnetwork {
		case "timezone":
			testnet = jobnet
		case "today":
			todaynet = jobnet
		case "yesterday":
			yesterdaynet = jobnet
		}
	}

	return testnet, todaynet, yesterdaynet
}

func isTodayInUTC() bool {
	return time.Now().Hour() >= 9
}

func TestNetworkTimestamp(t *testing.T) {
	output := runShow(t, true)
	jsData, err := parseShowOutput(output)
	if err != nil {
		t.Log("Failed to parse output of show utility.")
		t.Fatalf("Reason: %s", err)
	}

	testnet, todaynet, yesterdaynet := extractNetworkResults(jsData.Jobnetworks)
	if testnet == nil {
		t.Errorf("Network[%s] result not found.", "timezone")
	}
	if isTodayInUTC() {
		if todaynet == nil {
			t.Errorf("Network[%s] result not found.", "today")
		}
	} else {
		if yesterdaynet != nil {
			t.Errorf("Network[%s] result must not found, but did it.", "yesterday")
		}
	}
}
