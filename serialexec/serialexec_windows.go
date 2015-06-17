package serialexec

import (
	"path/filepath"
	"testing"

	"github.com/unirita/cutotest/util"
	"github.com/unirita/cutotest/util/db"
)

func assertDB(t *testing.T, nid int) {
	conn, err := db.Open(util.GetDBDirPath())
	if err != nil {
		t.Fatalf("DB file open failed: %v", err)
	}
	defer conn.Close()

	assertNetworkRecord(t, conn, 1)
	assertJobRecord(t, conn, nid, "job1", "job1.js", 0)
	assertJobRecord(t, conn, nid, "job2", "job2.vbs", 0)
	assertJobRecord(t, conn, nid, "job3", "job3.jar", 0)
	assertJobRecord(t, conn, nid, "job4", "job4.ps1", 0)
	assertJobRecord(t, conn, nid, "job5", "job5.bat", 0)
	assertJobRecord(t, conn, nid, "job6", "job6.exe", 0)
	assertJobRecord(t, conn, nid, "param", "param.bat", 1)
	assertJobRecord(t, conn, nid, "env", "env.bat", 2)
	assertJobRecord(t, conn, nid, "work", "work.bat", 3)
	assertJobRecord(t, conn, nid, "path", "path.bat", 4)
}

func assertJoblog(t *testing.T, nid int) {
	workPath := filepath.Join(util.GetCutoRoot(), "jobscript", "otherpath")
	scriptPath := filepath.Join(workPath, "path.bat")

	assertHasNoEmptyJoblog(t, nid, "job1")
	assertHasNoEmptyJoblog(t, nid, "job2")
	assertHasNoEmptyJoblog(t, nid, "job3")
	assertHasNoEmptyJoblog(t, nid, "job4")
	assertHasNoEmptyJoblog(t, nid, "job5")
	assertHasNoEmptyJoblog(t, nid, "job6")
	assertJoblogContainsStr(t, nid, "param", "param1=test1", "param2=test2")
	assertJoblogContainsStr(t, nid, "env", "ENV1=TEST1", "ENV2=TEST2")
	assertJoblogContainsStr(t, nid, "work", workPath)
	assertJoblogContainsStr(t, nid, "path", scriptPath)
}
