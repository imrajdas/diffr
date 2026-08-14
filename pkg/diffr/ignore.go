package diffr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// Built-in patterns applied unless --no-default-exclude is set.
var defaultExcludes = []string{
	".git",
	".git/",
	".svn",
	".svn/",
	".hg",
	".hg/",
	"node_modules",
	"node_modules/",
	"__pycache__",
	"__pycache__/",
	".DS_Store",
	"Thumbs.db",
}

type Ignore struct {
	gi *ignore.GitIgnore
}

func LoadIgnore(opts Options, left, right string) (*Ignore, error) {
	var lines []string
	if !opts.NoDefaultExclude {
		lines = append(lines, defaultExcludes...)
	}

	var files []string
	files = append(files, filepath.Join(".", ".diffrignore"))
	for _, root := range []string{left, right} {
		info, err := os.Stat(root)
		if err != nil || !info.IsDir() {
			continue
		}
		files = append(files, filepath.Join(root, ".diffrignore"))
	}
	if opts.IgnoreFile != "" {
		files = append(files, opts.IgnoreFile)
	}
	if !opts.NoGitignore {
		for _, root := range []string{left, right} {
			if gr := findGitRoot(root); gr != "" {
				files = append(files, filepath.Join(gr, ".gitignore"))
			}
		}
	}

	seen := make(map[string]struct{})
	for _, p := range files {
		abs, err := filepath.Abs(p)
		if err != nil {
			if opts.IgnoreFile != "" && p == opts.IgnoreFile {
				return nil, fmt.Errorf("ignore file: %w", err)
			}
			continue
		}
		if _, ok := seen[abs]; ok {
			continue
		}
		seen[abs] = struct{}{}

		data, err := os.ReadFile(p)
		if err != nil {
			if opts.IgnoreFile != "" && p == opts.IgnoreFile {
				return nil, fmt.Errorf("ignore file: %w", err)
			}
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			lines = append(lines, strings.TrimRight(line, "\r"))
		}
	}

	lines = append(lines, opts.Exclude...)
	if len(lines) == 0 {
		return &Ignore{}, nil
	}
	return &Ignore{gi: ignore.CompileIgnoreLines(lines...)}, nil
}

func (i *Ignore) Match(rel string) bool {
	if i == nil || i.gi == nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	if rel == "." || rel == "" {
		return false
	}
	if i.gi.MatchesPath(rel) {
		return true
	}
	if !strings.HasSuffix(rel, "/") && i.gi.MatchesPath(rel+"/") {
		return true
	}
	return false
}

func findGitRoot(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	dir := path
	if !info.IsDir() {
		dir = filepath.Dir(path)
	}
	for {
		git := filepath.Join(dir, ".git")
		if _, err := os.Stat(git); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
