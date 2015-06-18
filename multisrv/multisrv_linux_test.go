package multisrv

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/unirita/cutotest/util"
	"github.com/unirita/cutotest/util/container"
)

type hostParams struct {
	Containers [25]string
}

const imageName = "cuto/servant"

func Test255Job(t *testing.T) {
	util.InitCutoRoot()
	util.DeployTestData("multisrv")
	util.ComplementConfig("master.ini")
	util.ComplementConfig("servant.ini")

	s := util.NewServant()
	s.SetConfig("servant.ini")
	if err := s.Start(); err != nil {
		t.Fatalf("Servant start failed: %s", err)
	}
	defer s.Kill()

	params := new(hostParams)
	for i := 0; i < 25; i++ {
		name := fmt.Sprintf("TestContainer%d", i+1)
		cont := container.New(imageName, name)
		err := cont.Start()
		if err != nil {
			t.Fatalf("Failed to start container[%s]", name)
		}
		defer cont.Terminate()
		params.Containers[i], err = cont.IPAddress()
		if err != nil {
			t.Fatalf("Failed to get container[%s] IPAddress", name)
		}
	}

	if err := complement255JobDetail(params); err != nil {
		t.Fatalf("Failed to complete job detail csv file: %s", err)
	}

	m := util.NewMaster()
	m.SetConfig("master.ini")

	rc, err := m.SyntaxCheck("255Job")
	if err != nil {
		t.Fatalf("Master bpmn syntax check failed: %s", err)
	}
	if rc != 0 {
		t.Errorf("Master RC[%d] is not 0.", rc)
	}

	rc, err = m.Run("255Job")
	if err != nil {
		t.Fatalf("Master run failed: %s", err)
	}
	if rc != 0 {
		t.Errorf("Master RC[%d] is not 0.", rc)
	}
	logPath := util.GetLogPath("master.log")
	if util.HasLogError(logPath) {
		t.Errorf("There is error log in [%s]", logPath)
	}
}

func complement255JobDetail(params *hostParams) error {
	path := filepath.Join(util.GetCutoRoot(), "bpmn", "255Job.csv")
	tpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tpl.Execute(file, params); err != nil {
		return err
	}

	return nil
}
