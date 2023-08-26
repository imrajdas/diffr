package diffr

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/pmezard/go-difflib/difflib"
)

func compareFiles(file1, file2 string) (string, error) {
	content1, err := ioutil.ReadFile(file1)
	if err != nil {
		return "", err
	}

	content2, err := ioutil.ReadFile(file2)
	if err != nil {
		return "", err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(content1)),
		B:        difflib.SplitLines(string(content2)),
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
			diff, err := compareFiles(path1, "/dev/null")
			if err != nil {
				errorChan <- fmt.Errorf("error comparing files %s and /dev/null: %s", path1, err)
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
