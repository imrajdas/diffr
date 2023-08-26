package diffr

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/pmezard/go-difflib/difflib"
)

func readFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bufferedReader := bufio.NewReader(file)
	var contentBuilder strings.Builder
	bufferSize := 4096
	buffer := make([]byte, bufferSize)

	for {
		n, err := bufferedReader.Read(buffer)
		if err != nil {
			break
		}
		contentBuilder.Write(buffer[:n])
	}

	return contentBuilder.String(), nil
}

func compareFiles(file1, file2 string) (string, error) {
	content1, err := readFileContent(file1)
	if err != nil {
		return "", err
	}

	content2, err := readFileContent(file2)
	if err != nil {
		return "", err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(content1),
		B:        difflib.SplitLines(content2),
		FromFile: file1,
		ToFile:   file2,
		Context:  3,
	}

	diffs, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", err
	}

	return diffs, nil
}

func CompareDirectories(dir1, dir2 string, diffChan chan<- string, errorChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	filepath.Walk(dir1, func(path1 string, info os.FileInfo, err error) error {
		if err != nil {
			errorChan <- fmt.Errorf("error accessing %s: %s", path1, err)
			return nil
		}

		relPath, err := filepath.Rel(dir1, path1)
		if err != nil {
			errorChan <- fmt.Errorf("error getting relative path of %s: %s", path1, err)
			return nil
		}

		path2 := filepath.Join(dir2, relPath)

		if info.IsDir() {
			return nil
		}

		if _, err := os.Stat(path2); err == nil {
			diff, err := compareFiles(path1, path2)
			if err != nil {
				errorChan <- fmt.Errorf("error comparing files %s and %s: %s", path1, path2, err)
				return nil
			}

			if diff != "" {
				diffChan <- fmt.Sprintf("Differences in file: %s\n%s", relPath, diff)
			}
		} else if os.IsNotExist(err) {
			var nullDevice string
			if runtime.GOOS == "windows" {
				nullDevice = "NUL"
			} else {
				nullDevice = "/dev/null"
			}

			diff, err := compareFiles(path1, nullDevice)
			if err != nil {
				errorChan <- fmt.Errorf("error comparing files %s and %s: %s", path1, nullDevice, err)
				return nil
			}

			if diff != "" {
				diffChan <- fmt.Sprintf("Differences in file: %s (present in %s but not in %s)\n%s", relPath, dir1, dir2, diff)
			}
		} else {
			errorChan <- fmt.Errorf("error accessing %s: %s", path2, err)
		}

		return nil
	})
}
