package realtime

import (
	"time"
)

func makeBatchFileName(baseName string) string {
	return baseName + ".sh"
}

func waitProcessByPID(pid int, seconds time.Duration) {
	time.Sleep(time.Second * seconds)
}
