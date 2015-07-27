package realtime

import (
	"os"
)

func makeBatchFileName(baseName string) string {
	return baseName + ".bat"
}

func waitProcessByPID(pid int) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	proc.Wait()
}
