package util

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ClearDir(path string) {
	os.RemoveAll(path)
	os.MkdirAll(path, 0755)
}

func CopyFile(srcPath string, targetPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer target.Close()

	r := bufio.NewReader(src)
	w := bufio.NewWriter(target)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		w.Write(buf[:n])
	}
	return w.Flush()
}

func CopyDir(srcDir string, targetDir string) error {
	os.MkdirAll(targetDir, 0755)
	fis, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		srcFilePath := filepath.Join(srcDir, fi.Name())
		targetFilePath := filepath.Join(targetDir, fi.Name())

		if fi.IsDir() {
			err := CopyDir(srcFilePath, targetFilePath)
			if err != nil {
				return err
			}
		} else {
			err := CopyFile(srcFilePath, targetFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func PermitExecRecursive(path string) error {
	rootfi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !rootfi.IsDir() {
		return os.Chmod(path, 0755)
	}

	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		childPath := filepath.Join(path, fi.Name())
		PermitExecRecursive(childPath)
	}

	return nil
}

func IsExistNoEmptyFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Size() > 0
}

func ContainsInFile(path string, substr string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		if strings.Contains(s.Text(), substr) {
			return true
		}
	}

	return false
}

func CountInFile(path string, substr string) int {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()

	count := 0
	s := bufio.NewScanner(file)
	for s.Scan() {
		if strings.Contains(s.Text(), substr) {
			count++
		}
	}

	return count
}

func IsPatternExistInFile(path string, pattern string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	r, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	s := bufio.NewScanner(file)
	for s.Scan() {
		if r.MatchString(s.Text()) {
			return true
		}
	}

	return false
}
