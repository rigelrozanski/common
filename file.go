package common

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Credit: https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write

// ReadLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// WriteLines writes the lines to the given file.
func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func CreateEmptyFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

// replaces all strings in file
func ReplaceAllStringInFile(path, origStr, newStr string) error {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	output := bytes.Replace(input, []byte(origStr), []byte(newStr), -1)

	return ioutil.WriteFile(path, output, 0666)
}

// test if the file or folder exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
// credit: https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// move a file
func Move(src, dst string) error {
	return os.Rename(src, dst)
}

//___________________________________________________________________________________________-
// Directory Processing

// OperateOnDir - loop through files in the path and perform the Operation
func OperateOnDir(path string, op Operation) {
	filepath.Walk(path, visitFunc(op))
}

type Operation func(path string) error // nolint

func visitFunc(op Operation) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		return op(path)
	}
}

//___________________________________________________________________________________________-

// get the relative path to current loco
func GetRelPath(absPath, file string) (string, error) {
	curPath, err := filepath.Abs("")
	if err != nil {
		return "", err
	}

	goPath, _ := os.LookupEnv("GOPATH")

	relBoardsPath, err := filepath.Rel(curPath, path.Join(goPath,
		absPath))

	//create the boards directory if it doesn't exist
	os.Mkdir(relBoardsPath, os.ModePerm)

	relWbPath := path.Join(relBoardsPath, file)

	return relWbPath, err
}
