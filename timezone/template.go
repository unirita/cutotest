package timezone

import (
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/unirita/cutotest/util"
)

type hostParams struct {
	CntUTC string
	CntMST string
}

func complementJobDetail(hostUTC, hostMST string) error {
	params := new(hostParams)
	params.CntUTC = hostUTC
	params.CntMST = hostMST

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

type dateParams struct {
	Ysday string
}

func complementDB() error {
	yesterday := time.Now().AddDate(0, 0, -1)

	params := new(dateParams)
	params.Ysday = yesterday.Format("2006-01-02")

	path := filepath.Join(util.GetCutoRoot(), "data", "cuto.sqlite")
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
