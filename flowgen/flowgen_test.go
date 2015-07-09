package flowgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/unirita/cmptxt/comparer"
	"github.com/unirita/cutotest/util"
)

func TestMain(m *testing.M) {
	os.Exit(realTestMain(m))
}

func realTestMain(m *testing.M) int {
	util.InitCutoRoot()
	util.DeployTestData("flowgen")
	util.ComplementConfig("master.ini")
	util.ComplementConfig("servant.ini")

	s := util.NewServant()
	s.Start()
	defer s.Kill()

	code := m.Run()
	return code
}

func absFlowPath(filename string) string {
	return filepath.Join(util.GetCutoRoot(), "bpmn", filename)
}

func assertFlow(t *testing.T, name string) {
	actual, err := os.Open(absFlowPath(name + ".bpmn"))
	if err != nil {
		t.Errorf("Unexpected file open error occured: %s", err)
	}
	expected, err := os.Open(absFlowPath(name + "_expected.bpmn"))
	if err != nil {
		t.Errorf("Unexpected file open error occured: %s", err)
	}

	c := comparer.New()
	if !c.Compare(expected, actual) {
		t.Errorf("Generated flow file content is unexpected.")
		t.Logf("Actual file: %s", actual.Name())
		t.Logf("Expected file: %s", expected.Name())
	}
}

func TestParameter_NoParam(t *testing.T) {
	defer util.SaveEvidence("flowgen", "parameter", "no_param")

	f := util.NewFlowgen()
	rc, err := f.Run()
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 1 {
		t.Errorf("RC[%d] must be 1.", rc)
	}
	if !strings.HasPrefix(f.Stdout, "Usage") {
		t.Errorf("Usage was not displayed.")
	}
}

func TestParameter_TooManyParams(t *testing.T) {
	defer util.SaveEvidence("flowgen", "parameter", "many_params")

	f := util.NewFlowgen()
	rc, err := f.Run("param1", "param2")
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 1 {
		t.Errorf("RC[%d] must be 1.", rc)
	}
	if !strings.HasPrefix(f.Stdout, "Usage") {
		t.Errorf("Usage was not displayed.")
	}
}

func TestFileName_OnePeriod(t *testing.T) {
	defer util.SaveEvidence("flowgen", "filename", "one_period")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("test1.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("RC[%d] must be 0.", rc)
	}

	_, err = os.Stat(absFlowPath("test1.bpmn"))
	if os.IsNotExist(err) {
		t.Errorf("BPMN file was not output or file name wrong.")
	}
}

func TestFileName_TwoPeriods(t *testing.T) {
	defer util.SaveEvidence("flowgen", "filename", "two_priods")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("test2.sub.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("RC[%d] must be 0.", rc)
	}

	_, err = os.Stat(absFlowPath("test2.sub.bpmn"))
	if os.IsNotExist(err) {
		t.Errorf("BPMN file was not output or file name wrong.")
	}
}

func TestFileName_NoPeriod(t *testing.T) {
	defer util.SaveEvidence("flowgen", "filename", "no_priod")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("test3"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("RC[%d] must be 0.", rc)
	}

	_, err = os.Stat(absFlowPath("test3.bpmn"))
	if os.IsNotExist(err) {
		t.Errorf("BPMN file was not output or file name wrong.")
	}
}

func TestRelativePath(t *testing.T) {
	defer util.SaveEvidence("flowgen", "relative_path")

	os.Chdir(util.GetCutoRoot())
	f := util.NewFlowgen()
	rc, err := f.Run(filepath.Join("bpmn", "test4.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("RC[%d] must be 0.", rc)
	}

	_, err = os.Stat(absFlowPath("test4.bpmn"))
	if os.IsNotExist(err) {
		t.Errorf("BPMN file was not output or file name wrong.")
	}
}

func TestRun_Solo(t *testing.T) {
	defer util.SaveEvidence("flowgen", "run", "solo")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("j.flow"))
	if err != nil {
		t.Fatalf("Unexpexted flowgen error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("flowgen RC[%d] must be 0.", rc)
	}

	assertFlow(t, "j")

	m := util.NewMaster()
	m.UseConfig("master.ini")
	rc, err = m.Run("j")
	if err != nil {
		t.Fatalf("Unexpexted master error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("master RC[%d] must be 0.", rc)
	}
}

func TestRun_Serial(t *testing.T) {
	defer util.SaveEvidence("flowgen", "run", "serial")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("jj.flow"))
	if err != nil {
		t.Fatalf("Unexpexted flowgen error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("flowgen RC[%d] must be 0.", rc)
	}

	assertFlow(t, "jj")

	m := util.NewMaster()
	m.UseConfig("master.ini")
	rc, err = m.Run("jj")
	if err != nil {
		t.Fatalf("Unexpexted master error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("master RC[%d] must be 0.", rc)
	}
}

func TestRun_Parallel(t *testing.T) {
	defer util.SaveEvidence("flowgen", "run", "parallel")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("p.flow"))
	if err != nil {
		t.Fatalf("Unexpexted flowgen error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("flowgen RC[%d] must be 0.", rc)
	}

	assertFlow(t, "p")

	m := util.NewMaster()
	m.UseConfig("master.ini")
	rc, err = m.Run("p")
	if err != nil {
		t.Fatalf("Unexpexted master error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("master RC[%d] must be 0.", rc)
	}
}

func TestRun_Alternate(t *testing.T) {
	defer util.SaveEvidence("flowgen", "run", "alternate")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("jpj.flow"))
	if err != nil {
		t.Fatalf("Unexpexted flowgen error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("flowgen RC[%d] must be 0.", rc)
	}

	assertFlow(t, "jpj")

	m := util.NewMaster()
	m.UseConfig("master.ini")
	rc, err = m.Run("jpj")
	if err != nil {
		t.Fatalf("Unexpexted master error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("master RC[%d] must be 0.", rc)
	}
}

func TestRun_TwoParallels(t *testing.T) {
	defer util.SaveEvidence("flowgen", "run", "two_parallels")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("pp.flow"))
	if err != nil {
		t.Fatalf("Unexpexted flowgen error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("flowgen RC[%d] must be 0.", rc)
	}

	assertFlow(t, "pp")

	m := util.NewMaster()
	m.UseConfig("master.ini")
	rc, err = m.Run("pp")
	if err != nil {
		t.Fatalf("Unexpexted master error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("master RC[%d] must be 0.", rc)
	}
}

func TestRun_WithBlank(t *testing.T) {
	defer util.SaveEvidence("flowgen", "run", "with_blank")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("with_blank.flow"))
	if err != nil {
		t.Fatalf("Unexpexted flowgen error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("flowgen RC[%d] must be 0.", rc)
	}

	assertFlow(t, "with_blank")

	m := util.NewMaster()
	m.UseConfig("master.ini")
	rc, err = m.Run("with_blank")
	if err != nil {
		t.Fatalf("Unexpexted master error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("master RC[%d] must be 0.", rc)
	}
}

func TestError_IrregalName_Backslash(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "backslash")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("backslash.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_Slash(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "slash")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("slash.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_Colon(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "colon")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("colon.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_Asterisk(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "asterisk")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("asterisk.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_Question(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "question")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("question.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_LT(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "lt")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("lt.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_GT(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "gt")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("gt.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_Doller(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "doller")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("doller.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_Amp(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "amp")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("amp.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_BracketLeft(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "bracket_left")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("bracket_left.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_BracketRight(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "bracket_right")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("bracket_right.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_Comma(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "comma")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("comma.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_IrregalName_Hyphen(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "hyphen")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("hyphen.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_Syntax_NoJob(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "nojob")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("nojob.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_Syntax_EmptyBranch(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "empty_branch")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("empty_branch.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_Syntax_EmptyPath(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "irregal_name", "empty_path")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("empty_path.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestError_FileNotExists(t *testing.T) {
	defer util.SaveEvidence("flowgen", "error", "file_not_exists")

	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("file_not_exists.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 2 {
		t.Errorf("flowgen RC[%d] must be 2.", rc)
	}
}

func TestAlreadyExists(t *testing.T) {
	defer util.SaveEvidence("flowgen", "already_exists")

	// first time.
	f := util.NewFlowgen()
	rc, err := f.Run(absFlowPath("already_exists.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("RC[%d] must be 0.", rc)
	}

	_, err = os.Stat(absFlowPath("already_exists.bpmn"))
	if os.IsNotExist(err) {
		t.Errorf("BPMN file was not output or file name wrong.")
	}

	_, err = os.Stat(absFlowPath("already_exists.bpmn.bk"))
	if os.IsNotExist(err) {
		t.Errorf("Backup failed.")
	}
	if !util.ContainsInFile(absFlowPath("already_exists.bpmn.bk"), "dummy") {
		t.Errorf("Content of backup file is unexpected.")
	}

	// second time.
	rc, err = f.Run(absFlowPath("already_exists.flow"))
	if err != nil {
		t.Fatalf("Unexpexted error occured: %s", err)
	}
	if rc != 0 {
		t.Errorf("RC[%d] must be 0.", rc)
	}

	_, err = os.Stat(absFlowPath("already_exists.bpmn"))
	if os.IsNotExist(err) {
		t.Errorf("BPMN file was not output or file name wrong.")
	}

	_, err = os.Stat(absFlowPath("already_exists.bpmn.bk"))
	if os.IsNotExist(err) {
		t.Errorf("Backup failed.")
	}
	if util.ContainsInFile(absFlowPath("already_exists.bpmn.bk"), "dummy") {
		t.Errorf("Content of backup file is unexpected.")
	}
}
