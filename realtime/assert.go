package realtime

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/unirita/cutotest/util"
	"github.com/unirita/cutotest/util/db"
)

func assertSuccessRealtimeOutput(t *testing.T, res *result) {
	if res.Status != 0 {
		t.Errorf("output status => %d, want %d", res.Status, 0)
	}
	if res.Message != "Success." {
		t.Errorf("output message => %s, want %s", res.Message, "Success.")
	}
	if res.PID == 0 {
		t.Errorf("output pid must not be %d, but it is", 0)
	}
	if res.Network.Instance == 0 {
		t.Errorf("output instance ID must not be %d, but it is", 0)
	}
	if !regexp.MustCompile(`realtime_\d{14}`).MatchString(res.Network.Name) {
		t.Errorf("output network name does not match valid pattern.")
	}
}

func assertNotRemainNetworkFile(t *testing.T) {
	jobnetDir := filepath.Join(util.GetCutoRoot(), "bpmn")
	fis, err := ioutil.ReadDir(jobnetDir)
	if err != nil {
		t.Fatalf("Unexpected read directory error occured: %s", err)
	}

	matcher := regexp.MustCompile(`.+\.(bpmn|csv)`)
	for _, fi := range fis {
		if matcher.MatchString(fi.Name()) {
			t.Errorf("Network file remains: %s", fi.Name())
		}
	}
}

func assertIsNotNetworkEnds(t *testing.T) {
	conn, err := db.Open(util.GetDBDirPath())
	if err != nil {
		t.Fatalf("DB file open failed: %v", err)
	}
	defer conn.Close()

	nwks, err := conn.SelectJobNetworksByCond("ID=1 AND STATUS<>0")
	if err != nil {
		t.Fatalf("Unexpected DB error occured: %s", err)
	}
	if len(nwks) != 0 {
		t.Errorf("Network must not ends, but it does.")
	}
}

func assertIsJob1to3Successed(t *testing.T) {
	conn, err := db.Open(util.GetDBDirPath())
	if err != nil {
		t.Fatalf("DB file open failed: %v", err)
	}
	defer conn.Close()

	jobs1, err := conn.SelectJobsByCond("JOBNAME like 'job1.%' AND STATUS=1")
	if err != nil {
		t.Fatalf("Unexpected DB error occured: %s", err)
	}
	if len(jobs1) != 1 {
		t.Errorf("job1 must end successfully once, but dose not.")
	}

	jobs2, err := conn.SelectJobsByCond("JOBNAME like 'job2.%' AND STATUS=1")
	if err != nil {
		t.Fatalf("Unexpected DB error occured: %s", err)
	}
	if len(jobs2) != 1 {
		t.Errorf("job2 must end successfully once, but dose not.")
	}

	jobs3, err := conn.SelectJobsByCond("JOBNAME like 'job3.%' AND STATUS=1")
	if err != nil {
		t.Fatalf("Unexpected DB error occured: %s", err)
	}
	if len(jobs3) != 1 {
		t.Errorf("job3 must end successfully once, but dose not.")
	}
}
