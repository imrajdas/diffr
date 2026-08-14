package diffr

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

const maxTextBytes = 32 << 20

var errFileTooLarge = errors.New("file too large for text diff")

func compareText(fd FileDiff, opts Options) (FileDiff, bool, error) {
	left, err := readText(fd.LeftPath)
	if errors.Is(err, errFileTooLarge) {
		return compareBinaries(fd)
	}
	if err != nil {
		return fd, false, err
	}

	right, err := readText(fd.RightPath)
	if errors.Is(err, errFileTooLarge) {
		return compareBinaries(fd)
	}
	if err != nil {
		return fd, false, err
	}

	if equivalentText(left, right, opts) {
		return fd, true, nil
	}

	unified, err := makeUnified(fd.RelPath, left, right, fd.Status, opts.Context)
	if err != nil {
		return fd, false, err
	}
	fd.Unified = unified
	fd.Summary = "text file " + fd.Status.String()
	return fd, false, nil
}

func equivalentText(a, b string, opts Options) bool {
	if a == b {
		return true
	}
	if !opts.IgnoreCase && !opts.IgnoreWhitespace {
		return false
	}
	return normalizeText(a, opts) == normalizeText(b, opts)
}

func normalizeText(s string, opts Options) string {
	if !opts.IgnoreCase && !opts.IgnoreWhitespace {
		return s
	}
	lines := splitLines(s)
	var b strings.Builder
	for _, line := range lines {
		line = strings.TrimSuffix(line, "\n")
		if opts.IgnoreWhitespace {
			line = strings.Join(strings.Fields(line), "")
		}
		if opts.IgnoreCase {
			line = strings.ToLower(line)
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func makeUnified(rel, left, right string, status Status, context int) (string, error) {
	from, to := "a/"+rel, "b/"+rel
	if status == StatusAdded {
		from = "/dev/null"
	}
	if status == StatusRemoved {
		to = "/dev/null"
	}
	body, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        splitLines(left),
		B:        splitLines(right),
		FromFile: from,
		ToFile:   to,
		Context:  context,
	})
	if err != nil {
		return "", err
	}
	return gitDiffPrefix(rel, status) + body, nil
}

func gitDiffPrefix(rel string, status Status) string {
	var b strings.Builder
	fmt.Fprintf(&b, "diff --git a/%s b/%s\n", rel, rel)
	switch status {
	case StatusAdded:
		b.WriteString("new file mode 100644\n")
	case StatusRemoved:
		b.WriteString("deleted file mode 100644\n")
	}
	return b.String()
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	lines := difflib.SplitLines(s)
	if n := len(lines); n > 0 && lines[n-1] == "" {
		lines = lines[:n-1]
	}
	return lines
}

func readText(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.Size() > maxTextBytes {
		return "", errFileTooLarge
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
