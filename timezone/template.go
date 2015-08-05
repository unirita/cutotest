package timezone

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/unirita/cutotest/util"
)

type hostParams struct {
	CntUTC string
	CntMST string
}

func complementJobDetail(hostUTC, hostMST string) error {
	params := new(hostParams)
	params.CntUTC(hostUTC)
	params.CntMST(hostMST)

	path := filepath.Join(util.GetCutoRoot(), "bpmn", "timezone.csv")
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
