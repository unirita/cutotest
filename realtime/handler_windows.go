package realtime

import (
	"os"
	"time"
)

func makeBatchFileName(baseName string) string {
	return baseName + ".bat"
}

func waitProcessByPID(pid int, seconds time.Duration) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	proc.Wait()
}
