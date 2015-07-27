package realtime

import (
	"time"
)

func makeBatchFileName(baseName string) string {
	return baseName + ".sh"
}

func waitProcessByPID(pid int) {
	time.Sleep(time.Second * 5)
}
