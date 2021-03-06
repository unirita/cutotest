package timezone

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
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
	diffSec := now.Sub(target).Seconds()
	return diffSec < 300 && -300 < diffSec
}

func isTimeLocal(timeStr string) bool {
	now := time.Now()
	target, err := time.ParseInLocation(timeFormat, timeStr, time.Local)
	if err != nil {
		return false
	}
	diffSec := now.Sub(target).Seconds()
	return diffSec < 300 && -300 < diffSec
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
		t.Fatalf("Network result not found.", "timezone")
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

func TestNetworkTimestamp_Local(t *testing.T) {
	output := runShow(t, false)
	jsData, err := parseShowOutput(output)
	if err != nil {
		t.Log("Failed to parse output of show utility.")
		t.Fatalf("Reason: %s", err)
	}

	testnet := findNetworkResult(jsData.Jobnetworks, "timezone")
	if testnet == nil {
		t.Fatalf("Network result not found.", "timezone")
	}
	if !isTimeLocal(testnet.StartDate) {
		t.Error("Network startdate is not local timezone.")
	}
	if !isTimeLocal(testnet.EndDate) {
		t.Error("Network enddate is not local timezone.")
	}
	if !isTimeLocal(testnet.CreateDate) {
		t.Error("Network record createdate is not local timezone.")
	}
	if !isTimeLocal(testnet.UpdateDate) {
		t.Error("Network record updatedate is not local timezone.")
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
	for _, job := range testnet.Jobs {
		if !isTimeUTC(job.StartDate) {
			t.Errorf("Job[%s] startdate is not utc.", job.Jobname)
		}
		if !isTimeUTC(job.EndDate) {
			t.Errorf("Job[%s] enddate is not utc.", job.Jobname)
		}
		if !isTimeUTC(job.CreateDate) {
			t.Errorf("Job[%s] record createdate is not utc.", job.Jobname)
		}
		if !isTimeUTC(job.UpdateDate) {
			t.Errorf("Job[%s] record updatedate is not utc.", job.Jobname)
		}
	}
}

func TestJobTimestamp_Local(t *testing.T) {
	output := runShow(t, false)
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
	for _, job := range testnet.Jobs {
		if !isTimeLocal(job.StartDate) {
			t.Errorf("Job[%s] startdate is not local timezone.", job.Jobname)
		}
		if !isTimeLocal(job.EndDate) {
			t.Errorf("Job[%s] enddate is not local timezone.", job.Jobname)
		}
		if !isTimeLocal(job.CreateDate) {
			t.Errorf("Job[%s] record createdate is not local timezone.", job.Jobname)
		}
		if !isTimeLocal(job.UpdateDate) {
			t.Errorf("Job[%s] record updatedate is not local timezone.", job.Jobname)
		}
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
		if todaynet != nil {
			t.Errorf("Network[%s] result must not found, but did it.", "today")
		}
	} else {
		yesterdaynet := findNetworkResult(jsData.Jobnetworks, "yesterday")
		if yesterdaynet == nil {
			t.Errorf("Network[%s] result not found.", "yesterday")
		}
	}
}

func TestDateBorder_Local(t *testing.T) {
	output := runShow(t, false)
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

func TestJoblogTimestamp_Filename(t *testing.T) {
	joblogs := util.FindJoblog("joblog", 3, "receive")
	if len(joblogs) != 1 {
		t.Fatalf("len(joblogs) => %d, want %d.", len(joblogs), 1)
	}

	parts := strings.Split(joblogs[0], ".")
	timeStr := parts[len(parts)-3]

	now := time.Now()
	format := "20060102150405"
	target, err := time.ParseInLocation(format, timeStr, time.Local)
	if err != nil {
		t.Fatalf("Failed to parse joblog filename[%s] as time format.", joblogs[0])
	}
	diffSec := now.Sub(target).Seconds()
	if diffSec >= 300 || -300 >= diffSec {
		t.Errorf("Timestamp in joblog[%s] is not local timezone.")
	}
}

func TestJoblogTimestamp_Parameter(t *testing.T) {
	joblogs := util.FindJoblog("joblog", 3, "receive")
	if len(joblogs) != 1 {
		t.Fatalf("len(joblogs) => %d, want %d.", len(joblogs), 1)
	}

	file, err := os.Open(joblogs[0])
	if err != nil {
		t.Fatalf("Could not open joblog[%s].", joblogs[0])
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isScanned := false
	for scanner.Scan() {
		isScanned = true
		line := scanner.Text()
		if !isTimeLocal(line) {
			t.Errorf("Parameter[%s] is not local timestamp.", line)
		}
	}
	if !isScanned {
		t.Errorf("Joblog[%s] is empty.", joblogs[0])
	}
}

func TestLogTimestamp_Master(t *testing.T) {
	path := filepath.Join(util.GetCutoRoot(), "log", "master.log")
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Could not open master log[%s].", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		t.Fatalf("Master log[%s] is empty.", path)
	}
	firstLine := scanner.Text()
	timestamp := firstLine[:len(timeFormat)]
	if !isTimeLocal(timestamp) {
		t.Error("Timestamp in master log is not local timezone.")
		t.Log("Line example:")
		t.Log(firstLine)
	}
}

func TestLogTimestamp_Servant(t *testing.T) {
	path := filepath.Join(util.GetCutoRoot(), "log", "servant.log")
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Could not open servant log[%s].", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		t.Fatalf("Master log[%s] is empty.", path)
	}
	firstLine := scanner.Text()
	timestamp := firstLine[:len(timeFormat)]
	if !isTimeLocal(timestamp) {
		t.Error("Timestamp in servant log is not local timezone.")
		t.Log("Line example:")
		t.Log(firstLine)
	}
}
