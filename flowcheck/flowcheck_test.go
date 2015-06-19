package flowcheck

import (
	"strings"
	"testing"

	"github.com/unirita/cutotest/util"
)

type testCases struct {
	name string
	msg  string
}

var errorCases []testCases = getErrorCases()
var normalCases []testCases = getNormalCases()

func getErrorCases() []testCases {
	return []testCases{
		{"NoServiceTask", "ProcessFlow is empty."},
		{"NoStartEvent", "StartEvent element is required, and must be unique."},
		{"NoEndEvent", "EndEvent element is required, and must be unique."},
		{"MultiStartEvent", "StartEvent element is required, and must be unique."},
		{"MultiEndEvent", "EndEvent element is required, and must be unique."},
		{"ForbiddenJobName1", "Job name[job1\\.bat] includes forbidden character."},
		{"ForbiddenJobName2", "Job name[job1/.bat] includes forbidden character."},
		{"ForbiddenJobName3", "Job name[job1:.bat] includes forbidden character."},
		{"ForbiddenJobName4", "Job name[job1*.bat] includes forbidden character."},
		{"ForbiddenJobName5", "Job name[job1?.bat] includes forbidden character."},
		{"ForbiddenJobName6", "Job name[job1\".bat] includes forbidden character."},
		{"ForbiddenJobName7", "Job name[job1<.bat] includes forbidden character."},
		{"ForbiddenJobName8", "Job name[job1>.bat] includes forbidden character."},
		{"ForbiddenJobName9", "Job name[job1|.bat] includes forbidden character."},
		{"ForbiddenJobName10", "Job name[job1$.bat] includes forbidden character."},
		{"ForbiddenJobName11", "Job name[job1&.bat] includes forbidden character."},
		{"StartWithoutStartEvent", "There is no element which connects with startEvent."},
		{"EndWithoutEndEvent", "There is no element which connects with endEvent."},
		{"Isolation", "Isolated element is detected."},
		{"DuplicateID", "Element[id = job] duplicated."},
		{"BranchWithoutGateway", "ServiceTask cannot connect with over 1 element."},
		{"MergeWithoutGateway", "Element[id = job3] duplicated."},
		{"EndBeforeMerge", "EndEvent cannot connect with over 1 element."},
		{"NestedBranch", "Cannot nest branches."},
		{"NotMerge", "There is a sequenceFlow which refers imaginary element[id = scripttask1]."},
	}
}

func getNormalCases() []testCases {
	return []testCases{
		{"WithExtraTags", "CTM020I"},
	}
}

func TestFlowcheck(t *testing.T) {
	defer util.SaveEvidence("flowcheck")
	util.InitCutoRoot()
	util.DeployTestData("flowcheck")
	util.ComplementConfig("master.ini")
	util.ComplementConfig("servant.ini")

	s := util.NewServant()
	s.SetConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	// ここからmaster起動
	m := util.NewMaster()
	m.SetConfig("master.ini")
	// 異常ケース
	for _, test := range errorCases {
		isOK := true
		// シミュレーション実行
		rc, err := m.SyntaxCheck(test.name)
		if err != nil {
			t.Errorf("Master SyntaxCheck failed: %s", err)
			isOK = false
		}
		if rc != 1 {
			t.Logf("Master stdout: %s", m.Stdout)
			t.Logf("Master stderr: %s", m.Stderr)
			t.Errorf("Master RC[%d] is not 1.", rc)
			isOK = false
		}
		if !strings.Contains(m.Stdout, test.msg) {
			t.Errorf("Invalid stdout message. - %v", m.Stdout)
			isOK = false
		}
		// 実行
		rc, err = m.Run(test.name)
		if err != nil {
			t.Fatalf("Master run failed: %s", err)
			isOK = false
		}
		if rc != 1 {
			t.Logf("Master stdout: %s", m.Stdout)
			t.Logf("Master stderr: %s", m.Stderr)
			t.Errorf("Master RC[%d] is not 1.", rc)
			isOK = false
		}
		if !strings.Contains(m.Stdout, test.msg) {
			t.Errorf("Invalid stdout message. - %v", m.Stdout)
			isOK = false
		}
		if isOK {
			t.Logf("%v testcase OK.", test.name)
		}
	}
	// 正常ケース
	for _, test := range normalCases {
		isOK := true
		// シミュレーション実行
		rc, err := m.SyntaxCheck(test.name)
		if err != nil {
			t.Errorf("Master SyntaxCheck failed: %s", err)
			isOK = false
		}
		if rc != 0 {
			t.Logf("Master stdout: %s", m.Stdout)
			t.Logf("Master stderr: %s", m.Stderr)
			t.Errorf("Master RC[%d] is not 0.", rc)
			isOK = false
		}
		if !strings.Contains(m.Stdout, test.msg) {
			t.Errorf("Invalid stdout message. - %v", m.Stdout)
			isOK = false
		}
		// 実行
		rc, err = m.Run(test.name)
		if err != nil {
			t.Errorf("Master run failed: %s", err)
			isOK = false
		}
		if rc != 0 {
			t.Logf("Master stdout: %s", m.Stdout)
			t.Logf("Master stderr: %s", m.Stderr)
			t.Errorf("Master RC[%d] is not 0.", rc)
			isOK = false
		}
		if isOK {
			t.Logf("%v testcase OK.", test.name)
		}
	}
	// ログファイルの確認
	logPath := util.GetLogPath("servant.log")
	if util.HasLogError(logPath) {
		t.Errorf("There is error log in [%s]", logPath)
	}
}
