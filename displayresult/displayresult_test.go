package displayresult

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"encoding/json"
	"os/exec"
	"path/filepath"

	"github.com/unirita/cutotest/util"
)

const abnormal = 9

func createData() error {
	cur, _ := os.Getwd()

	if err := os.Chdir(filepath.Join(util.GetCutoRoot(), "data")); err != nil {
		return err
	}
	defer os.Chdir(cur)
	fmt.Println(filepath.Join(util.GetCutoRoot(), "data"))

	cmd := exec.Command("./initdata.sh")
	if err := cmd.Start(); err != nil {
		return err
	}
	err := cmd.Wait()
	if err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				rc := s.ExitStatus()
				if rc == 0 {
					return nil
				}
				return fmt.Errorf("initdata.sh RC=%v", rc)
			}
		}
	} else {
		return nil
	}
	return errors.New("createData() Unknown error.")
}

func TestMain(m *testing.M) {
	util.InitCutoRoot()
	util.DeployTestData("displayresult")
	util.ComplementConfig("master.ini")
	util.ComplementConfig("servant.ini")

	if err := createData(); err != nil {
		panic(err)
	}

	rc := m.Run()

	// Cannot use defer because of os.Exit.
	util.SaveEvidence("displayresult")

	os.Exit(rc)
}

// No.1 - No.3まで
func TestDisplayresult_Default(t *testing.T) {
	verifyJn := "Normal"
	verifyJobs := []string{"job1", "job2"}

	show := util.NewShow()
	show.UseConfig("master.ini")
	rc, err := show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	jsData := new(OutputRoot)
	err = json.Unmarshal([]byte(show.Stdout), jsData)
	if err != nil {
		t.Errorf("Json Parse failed. : %v", err)
	}
	if len(jsData.Jobnetworks) == 0 {
		t.Error("Getted data count 0.")
	}
	var chkJn *OutputJobNet
	now := time.Now()
	today := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	for _, jobnet := range jsData.Jobnetworks {
		if !strings.HasPrefix(jobnet.StartDate, today) {
			t.Errorf("Get Jobnet record's only StartDate item '%v', but getted '%v' record exist.", today, jobnet.StartDate)
		}
		if jobnet.Jobnetwork == verifyJn {
			chkJn = jobnet
		}
	}
	if len(chkJn.Jobs) != len(verifyJobs) {
		t.Errorf("%v job records must be even, but getted %v records", len(verifyJobs), len(chkJn.Jobs))
	}
	for _, job := range chkJn.Jobs {
		if job.Jobname != "job1" && job.Jobname != "job2" {
			t.Errorf("Got invalid job record. : jobname(%v)", job.Jobname)
		}
	}
	t.Log("No.1 - No.3 PASS.")
}

// No.5
func TestDisplayresult_JobnetSetting(t *testing.T) {
	verifyJn := "Running"

	show := util.NewShow()
	show.UseConfig("master.ini")
	show.AddJobnet(verifyJn)
	rc, err := show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	jsData := new(OutputRoot)
	err = json.Unmarshal([]byte(show.Stdout), jsData)
	if err != nil {
		t.Errorf("Json Parse failed. : %v", err)
	}
	if len(jsData.Jobnetworks) != 1 {
		t.Errorf("1 jobnet records must be even, but getted %v records", len(jsData.Jobnetworks))
	}
	jn := jsData.Jobnetworks[0]
	if jn.Jobnetwork != verifyJn {
		t.Errorf("Get jobnet name must be '%v', but getted by '%v'", verifyJn, jn.Jobnetwork)
	}
	t.Log("No.5 PASS.")
}

// No.6,No.7
func TestDisplayresult_PeriodSetting(t *testing.T) {
	// No.6
	now := time.Now()
	to := fmt.Sprintf("%04d%02d%02d", now.Year(), now.Month(), now.Day())
	day4ago := now.AddDate(0, 0, -4)
	from := fmt.Sprintf("%04d%02d%02d", day4ago.Year(), day4ago.Month(), day4ago.Day())

	show := util.NewShow()
	show.UseConfig("master.ini")
	show.AddFrom(from)
	show.AddTo(to)
	rc, err := show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	jsData := new(OutputRoot)
	err = json.Unmarshal([]byte(show.Stdout), jsData)
	if err != nil {
		t.Errorf("Json Parse failed. : %v", err)
	}
	if len(jsData.Jobnetworks) != 20 {
		t.Errorf("20 jobnet records must be even, but getted %v records", len(jsData.Jobnetworks))
	}
	t.Log("No.6 PASS.")

	// No.7
	yesterday := now.AddDate(0, 0, -1)
	to = fmt.Sprintf("%04d%02d%02d", yesterday.Year(), yesterday.Month(), yesterday.Day())
	day3ago := now.AddDate(0, 0, -3)
	from = fmt.Sprintf("%04d%02d%02d", day3ago.Year(), day3ago.Month(), day3ago.Day())

	show = util.NewShow()
	show.UseConfig("master.ini")
	show.AddFrom(from)
	show.AddTo(to)
	rc, err = show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	jsData = new(OutputRoot)
	err = json.Unmarshal([]byte(show.Stdout), jsData)
	if err != nil {
		t.Errorf("Json Parse failed. : %v", err)
	}
	if len(jsData.Jobnetworks) != 12 {
		t.Errorf("12 jobnet records must be even, but getted %v records", len(jsData.Jobnetworks))
	}
	t.Log("No.7 PASS.")
}

// No.8 - No.10
func TestDisplayresult_StatusSetting(t *testing.T) {
	// No.8
	status := "normal"
	verifyJn := "Normal"

	show := util.NewShow()
	show.UseConfig("master.ini")
	show.AddStatus(status)
	rc, err := show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	jsData := new(OutputRoot)
	err = json.Unmarshal([]byte(show.Stdout), jsData)
	if err != nil {
		t.Errorf("Json Parse failed. : %v", err)
	}
	if len(jsData.Jobnetworks) != 1 {
		t.Errorf("1 jobnet records must be even, but getted %v records", len(jsData.Jobnetworks))
	}
	jn := jsData.Jobnetworks[0]
	if jn.Jobnetwork != verifyJn {
		t.Errorf("Get jobnet name must be '%v', but getted by '%v'", verifyJn, jn.Jobnetwork)
	}
	t.Log("No.8 PASS.")

	// No.9
	status = "running"
	verifyJn = "Running"

	show = util.NewShow()
	show.UseConfig("master.ini")
	show.AddStatus(status)
	rc, err = show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	jsData = new(OutputRoot)
	err = json.Unmarshal([]byte(show.Stdout), jsData)
	if err != nil {
		t.Errorf("Json Parse failed. : %v", err)
	}
	if len(jsData.Jobnetworks) != 1 {
		t.Errorf("1 jobnet records must be even, but getted %v records", len(jsData.Jobnetworks))
	}
	jn = jsData.Jobnetworks[0]
	if jn.Jobnetwork != verifyJn {
		t.Errorf("Get jobnet name must be '%v', but getted by '%v'", verifyJn, jn.Jobnetwork)
	}
	t.Log("No.9 PASS.")

	// No.10
	status = "abnormal"
	verifyJn = "Abnormal"

	show = util.NewShow()
	show.UseConfig("master.ini")
	show.AddStatus(status)
	rc, err = show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	jsData = new(OutputRoot)
	err = json.Unmarshal([]byte(show.Stdout), jsData)
	if err != nil {
		t.Errorf("Json Parse failed. : %v", err)
	}
	if len(jsData.Jobnetworks) != 1 {
		t.Errorf("1 jobnet records must be even, but getted %v records", len(jsData.Jobnetworks))
	}
	jn = jsData.Jobnetworks[0]
	if jn.Jobnetwork != verifyJn {
		t.Errorf("Get jobnet name must be '%v', but getted by '%v'", verifyJn, jn.Jobnetwork)
	}
	t.Log("No.10 PASS.")
}

// No.11 - No.12
func TestDisplayresult_FormatSetting(t *testing.T) {
	// No.11
	format := "json"

	show := util.NewShow()
	show.UseConfig("master.ini")
	show.AddFormat(format)
	rc, err := show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	jsData := new(OutputRoot)
	err = json.Unmarshal([]byte(show.Stdout), jsData)
	if err != nil {
		t.Errorf("Json Parse failed. : %v", err)
	}
	if len(jsData.Jobnetworks) != 4 {
		t.Errorf("4 jobnet records must be even, but getted %v records", len(jsData.Jobnetworks))
	}
	now := time.Now()
	today := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	for _, jobnet := range jsData.Jobnetworks {
		if !strings.HasPrefix(jobnet.StartDate, today) {
			t.Errorf("Get Jobnet record's only StartDate item '%v', but getted '%v' record exist.", today, jobnet.StartDate)
		}
	}
	t.Log("No.11 PASS.")

	// No.12
	format = "csv"

	show = util.NewShow()
	show.UseConfig("master.ini")
	show.AddFormat(format)
	rc, err = show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	jsData = new(OutputRoot)
	err = json.Unmarshal([]byte(show.Stdout), jsData)
	if err == nil {
		t.Errorf("Get output format must be json, but successed json parse.")
	}
	if strings.Count(show.Stdout, "\n") != 11 {
		t.Errorf("11 jobnet records must be even, but getted %v records", strings.Count(show.Stdout, "\n"))
	}
	t.Log("No.12 PASS.")
}

// No.13
func TestDisplayresult_MultiSetting(t *testing.T) {
	// No.13
	status := "abnormal"

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	to := fmt.Sprintf("%04d%02d%02d", yesterday.Year(), yesterday.Month(), yesterday.Day())
	day3ago := now.AddDate(0, 0, -3)
	from := fmt.Sprintf("%04d%02d%02d", day3ago.Year(), day3ago.Month(), day3ago.Day())

	show := util.NewShow()
	show.UseConfig("master.ini")
	show.AddStatus(status)
	show.AddFrom(from)
	show.AddTo(to)
	rc, err := show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	jsData := new(OutputRoot)
	err = json.Unmarshal([]byte(show.Stdout), jsData)
	if err != nil {
		t.Errorf("Json Parse failed. : %v", err)
	}
	if len(jsData.Jobnetworks) != 3 {
		t.Errorf("3 jobnet records must be even, but getted %v records", len(jsData.Jobnetworks))
	}
	for _, jobnet := range jsData.Jobnetworks {
		if jobnet.Status != abnormal {
			t.Error("Exist invalid status of jobnet records.")
		}
	}
	t.Log("No.13 PASS.")
}

// No.14
func TestDisplayresult_Version(t *testing.T) {
	// No.14 Versioin info.
	show := util.NewShow()
	rc, err := show.Version()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	t.Log("No.14 PASS.")
}

// No.15
func TestDisplayresult_Help(t *testing.T) {
	// No.15 Versioin info.
	show := util.NewShow()
	rc, err := show.Help()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 0 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 0.", rc)
	}
	t.Log("No.15 PASS.")
}

// No.16-
func TestDisplayresult_ErrorCase(t *testing.T) {
	// No.16 invalid config
	show := util.NewShow()
	show.UseConfig("xxx")
	rc, err := show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 8 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 8.", rc)
	}
	t.Log("No.16 PASS.")

	// No.17 0 record
	show = util.NewShow()
	show.UseConfig("master.ini")
	show.AddTo("20100101")
	rc, err = show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 4 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 4.", rc)
	}
	t.Log("No.17 PASS.")

	// No.18 Invalid DB
	show = util.NewShow()
	show.UseConfig("error.ini")
	rc, err = show.Run()
	if err != nil {
		t.Errorf("Show Run failed : %v", err)
	}
	if rc != 12 {
		t.Logf("Show stdout: %s", show.Stdout)
		t.Logf("Show stderr: %s", show.Stderr)
		t.Errorf("Show RC[%d] is not 12.", rc)
	}
	t.Log("No.18 PASS.")
}
