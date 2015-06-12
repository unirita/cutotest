package util

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ClearDir(path string) {
	os.RemoveAll(path)
	os.Mkdir(path, 0755)
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
