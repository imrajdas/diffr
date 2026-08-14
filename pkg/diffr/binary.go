package diffr

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func compareBinaries(fd FileDiff) (FileDiff, bool, error) {
	fd.Kind = KindBinary

	var leftSize, rightSize int64
	if fd.LeftPath != "" {
		sum, n, err := hashFile(fd.LeftPath)
		if err != nil {
			return fd, false, err
		}
		fd.LeftHash = sum
		leftSize = n
	}
	if fd.RightPath != "" {
		sum, n, err := hashFile(fd.RightPath)
		if err != nil {
			return fd, false, err
		}
		fd.RightHash = sum
		rightSize = n
	}

	if fd.LeftPath != "" && fd.RightPath != "" && fd.LeftHash == fd.RightHash {
		return fd, true, nil
	}

	switch fd.Status {
	case StatusAdded:
		fd.Summary = fmt.Sprintf("binary file added (%s)", formatSize(rightSize))
	case StatusRemoved:
		fd.Summary = fmt.Sprintf("binary file removed (%s)", formatSize(leftSize))
	default:
		fd.Summary = fmt.Sprintf("binary files differ (%s vs %s)", formatSize(leftSize), formatSize(rightSize))
	}
	return fd, false, nil
}

func hashFile(path string) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}

func sameHash(a, b string) (bool, error) {
	ha, _, err := hashFile(a)
	if err != nil {
		return false, err
	}
	hb, _, err := hashFile(b)
	if err != nil {
		return false, err
	}
	return ha == hb, nil
}

func formatSize(n int64) string {
	switch {
	case n < 1024:
		return fmt.Sprintf("%d B", n)
	case n < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(n)/1024)
	default:
		return fmt.Sprintf("%.1f MB", float64(n)/(1024*1024))
	}
}
