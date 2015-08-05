package timezone

import (
	"testing"
	"time"

	"github.com/unirita/cutotest/util"
)

const dateFormat = "20060102"
const timeFormat = "2006-01-02 15:04:05.000"

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

func findNetworkResult(jobnets []*OutputJobNet, name string) *OutputJobNet {
	for _, jobnet := range jobnets {
		if jobnet.Jobnetwork == name {
			return jobnet
		}
	}
	return nil
}

func isTodayInUTC() bool {
	return time.Now().Hour() >= 9
}

func isTimeUTC(timeStr string) bool {
	now := time.Now().UTC()
	target, err := time.ParseInLocation(timeFormat, timeStr, time.UTC)
	if err != nil {
		return false
	}
	diff := now.Sub(target)
	return diff.Seconds() < 300 && -300 < diff.Seconds()
}

func TestNetworkTimestamp_UTC(t *testing.T) {
	output := runShow(t, true)
	jsData, err := parseShowOutput(output)
	if err != nil {
		t.Log("Failed to parse output of show utility.")
		t.Fatalf("Reason: %s", err)
	}

	testnet := findNetworkResult(jsData.Jobnetworks, "timezone")
	if testnet == nil {
		t.Errorf("Network result not found.", "timezone")
	}
	if !isTimeUTC(testnet.StartDate) {
		t.Error("Network startdate is not utc.")
	}
	if !isTimeUTC(testnet.EndDate) {
		t.Error("Network enddate is not utc.")
	}
	if !isTimeUTC(testnet.CreateDate) {
		t.Error("Network record createdate is not utc.")
	}
	if !isTimeUTC(testnet.UpdateDate) {
		t.Error("Network record updatedate is not utc.")
	}
}

func TestJobTimestamp_UTC(t *testing.T) {
	output := runShow(t, true)
	jsData, err := parseShowOutput(output)
	if err != nil {
		t.Log("Failed to parse output of show utility.")
		t.Fatalf("Reason: %s", err)
	}

	testnet := findNetworkResult(jsData.Jobnetworks, "timezone")
	if testnet == nil {
		t.Errorf("Network result not found.", "timezone")
	}
	if len(testnet.Jobs) != 3 {
		t.Fatalf("len(testnet.Jobs) => %d, wants %d.", len(testnet.Jobs), 3)
	}
}

func TestDateBorder_UTC(t *testing.T) {
	output := runShow(t, true)
	jsData, err := parseShowOutput(output)
	if err != nil {
		t.Log("Failed to parse output of show utility.")
		t.Fatalf("Reason: %s", err)
	}
	if isTodayInUTC() {
		todaynet := findNetworkResult(jsData.Jobnetworks, "today")
		if todaynet == nil {
			t.Errorf("Network[%s] result not found.", "today")
		}
	} else {
		yesterdaynet := findNetworkResult(jsData.Jobnetworks, "yesterday")
		if yesterdaynet != nil {
			t.Errorf("Network[%s] result must not found, but did it.", "yesterday")
		}
	}
}
