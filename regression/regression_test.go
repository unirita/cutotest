package regression

import (
	"testing"

	"github.com/unirita/cutotest/util"
)

func TestRegression(t *testing.T) {
	util.InitCutoRoot()
	util.DeployTestData("regression")
	t.Log(util.ComplementConfig("master.ini"))
	//m := util.NewMaster()
	//s := util.NewServant()
}
