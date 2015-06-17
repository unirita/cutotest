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
	assertJobRecord(t, conn, nid, "job1", "job1", 0)
	assertJobRecord(t, conn, nid, "job2", "job2.sh", 0)
	assertJobRecord(t, conn, nid, "param", "param.sh", 1)
	assertJobRecord(t, conn, nid, "env", "env.sh", 2)
	assertJobRecord(t, conn, nid, "work", "work.sh", 3)
	assertJobRecord(t, conn, nid, "path", "path.sh", 4)
}

func assertJoblog(t *testing.T, nid int) {
	workPath := filepath.Join(util.GetCutoRoot(), "jobscript", "otherpath")
	scriptPath := filepath.Join(workPath, "path.sh")

	assertHasNoEmptyJoblog(t, nid, "job1")
	assertHasNoEmptyJoblog(t, nid, "job2")
	assertJoblogContainsStr(t, nid, "param", "param1=test1", "param2=test2")
	assertJoblogContainsStr(t, nid, "env", "ENV1=TEST1", "ENV2=TEST2")
	assertJoblogContainsStr(t, nid, "work", workPath)
	assertJoblogContainsStr(t, nid, "path", scriptPath)
}
