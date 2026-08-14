package diffr

import (
	"os"
	"strings"
)

const (
	ansiReset = "\x1b[0m"
	ansiRed   = "\x1b[31m"
	ansiGreen = "\x1b[32m"
	ansiCyan  = "\x1b[36m"
	ansiBold  = "\x1b[1m"
	ansiDim   = "\x1b[2m"
)

func stdoutIsTTY() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func colorizePatch(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		lines[i] = colorizeLine(line)
	}
	return strings.Join(lines, "\n")
}

func colorizeLine(line string) string {
	switch {
	case strings.HasPrefix(line, "+++"):
		return ansiGreen + line + ansiReset
	case strings.HasPrefix(line, "---"):
		return ansiRed + line + ansiReset
	case strings.HasPrefix(line, "@@"):
		return ansiCyan + line + ansiReset
	case strings.HasPrefix(line, "+"):
		return ansiGreen + line + ansiReset
	case strings.HasPrefix(line, "-"):
		return ansiRed + line + ansiReset
	case strings.HasPrefix(line, "diff "),
		strings.HasPrefix(line, "new file"),
		strings.HasPrefix(line, "deleted file"),
		strings.HasPrefix(line, "index "):
		return ansiBold + line + ansiReset
	case strings.HasPrefix(line, "#"), strings.HasPrefix(line, "Binary "):
		return ansiDim + line + ansiReset
	default:
		return line
	}
}
