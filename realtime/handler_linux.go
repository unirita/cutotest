package realtime

import (
	"time"
)

func makeBatchFileName(baseName string) string {
	return baseName + ".sh"
}

func waitProcessByPID(pid, seconds int) {
	time.Sleep(time.Second * seconds)
}
