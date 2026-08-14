package diffr

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Status int

const (
	StatusChanged Status = iota
	StatusAdded
	StatusRemoved
)

func (s Status) String() string {
	switch s {
	case StatusAdded:
		return "added"
	case StatusRemoved:
		return "removed"
	default:
		return "changed"
	}
}

type FileDiff struct {
	RelPath   string
	Status    Status
	Kind      Kind
	Unified   string
	Summary   string
	LeftPath  string
	RightPath string
	DiffPNG   []byte
	LeftHash  string
	RightHash string
}

type Stats struct {
	Changed   int
	Added     int
	Removed   int
	Identical int
	Duration  time.Duration
	Elapsed   string
}

type Result struct {
	Left     string
	Right    string
	LeftAbs  string
	RightAbs string
	LeftDir  bool
	RightDir bool
	Files    []FileDiff
	Stats    Stats
}

func (r *Result) HasChanges() bool {
	return r.Stats.Changed+r.Stats.Added+r.Stats.Removed > 0
}

type Options struct {
	Exclude          []string
	IgnoreFile       string
	NoDefaultExclude bool
	NoGitignore      bool
	IgnoreWhitespace bool
	IgnoreCase       bool
	Context          int
}

type filePair struct {
	rel   string
	left  string
	right string
}

func Compare(left, right string, opts Options) (*Result, error) {
	start := time.Now()
	if opts.Context <= 0 {
		opts.Context = 3
	}

	leftAbs, err := filepath.Abs(left)
	if err != nil {
		return nil, err
	}
	rightAbs, err := filepath.Abs(right)
	if err != nil {
		return nil, err
	}

	leftInfo, err := os.Stat(leftAbs)
	if err != nil {
		return nil, fmt.Errorf("left path: %w", err)
	}
	rightInfo, err := os.Stat(rightAbs)
	if err != nil {
		return nil, fmt.Errorf("right path: %w", err)
	}
	if leftInfo.IsDir() != rightInfo.IsDir() {
		return nil, fmt.Errorf("cannot compare a file with a directory")
	}

	ign, err := LoadIgnore(opts, leftAbs, rightAbs)
	if err != nil {
		return nil, err
	}

	res := &Result{
		Left:     left,
		Right:    right,
		LeftAbs:  leftAbs,
		RightAbs: rightAbs,
		LeftDir:  leftInfo.IsDir(),
		RightDir: rightInfo.IsDir(),
	}

	var pairs []filePair
	if !leftInfo.IsDir() {
		pairs = []filePair{{
			rel:   filepath.ToSlash(filepath.Base(leftAbs)),
			left:  leftAbs,
			right: rightAbs,
		}}
	} else {
		leftFiles, err := collectFiles(leftAbs, ign)
		if err != nil {
			return nil, fmt.Errorf("walking %s: %w", leftAbs, err)
		}
		rightFiles, err := collectFiles(rightAbs, ign)
		if err != nil {
			return nil, fmt.Errorf("walking %s: %w", rightAbs, err)
		}

		seen := make(map[string]struct{}, len(leftFiles)+len(rightFiles))
		var rels []string
		for rel := range leftFiles {
			seen[rel] = struct{}{}
			rels = append(rels, rel)
		}
		for rel := range rightFiles {
			if _, ok := seen[rel]; ok {
				continue
			}
			rels = append(rels, rel)
		}
		sort.Strings(rels)

		pairs = make([]filePair, 0, len(rels))
		for _, rel := range rels {
			pairs = append(pairs, filePair{
				rel:   rel,
				left:  leftFiles[rel],
				right: rightFiles[rel],
			})
		}
	}

	for _, p := range pairs {
		fd, identical, err := comparePair(p, opts)
		if err != nil {
			return nil, err
		}
		if identical {
			res.Stats.Identical++
			continue
		}
		switch fd.Status {
		case StatusAdded:
			res.Stats.Added++
		case StatusRemoved:
			res.Stats.Removed++
		default:
			res.Stats.Changed++
		}
		res.Files = append(res.Files, fd)
	}

	res.Stats.Duration = time.Since(start)
	res.Stats.Elapsed = res.Stats.Duration.Round(time.Millisecond).String()
	return res, nil
}

func collectFiles(root string, ign *Ignore) (map[string]string, error) {
	files := make(map[string]string)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			if ign.Match(rel) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if ign.Match(rel) {
			return nil
		}
		files[rel] = path
		return nil
	})
	return files, err
}

func comparePair(p filePair, opts Options) (FileDiff, bool, error) {
	fd := FileDiff{RelPath: p.rel, LeftPath: p.left, RightPath: p.right}
	switch {
	case p.left == "":
		fd.Status = StatusAdded
	case p.right == "":
		fd.Status = StatusRemoved
	default:
		fd.Status = StatusChanged
	}

	var leftKind, rightKind Kind
	var err error
	if p.left != "" {
		leftKind, err = sniffKind(p.left)
		if err != nil {
			return fd, false, err
		}
	}
	if p.right != "" {
		rightKind, err = sniffKind(p.right)
		if err != nil {
			return fd, false, err
		}
	}
	fd.Kind = mergeKind(leftKind, rightKind, p.left != "", p.right != "")

	switch fd.Kind {
	case KindImage:
		return compareImages(fd)
	case KindPDF:
		return comparePDFs(fd, opts)
	case KindBinary:
		return compareBinaries(fd)
	default:
		return compareText(fd, opts)
	}
}

func (r *Result) Patch() string {
	var b strings.Builder
	for _, f := range r.Files {
		if f.Kind == KindText || (f.Kind == KindPDF && f.Unified != "") {
			if f.Summary != "" && f.Kind != KindText {
				fmt.Fprintf(&b, "# %s\n", f.Summary)
			}
			if f.Unified != "" {
				b.WriteString(f.Unified)
				if !strings.HasSuffix(f.Unified, "\n") {
					b.WriteByte('\n')
				}
			}
			continue
		}

		b.WriteString(gitDiffPrefix(f.RelPath, f.Status))
		switch f.Status {
		case StatusAdded:
			fmt.Fprintf(&b, "Binary files /dev/null and b/%s differ\n", f.RelPath)
		case StatusRemoved:
			fmt.Fprintf(&b, "Binary files a/%s and /dev/null differ\n", f.RelPath)
		default:
			fmt.Fprintf(&b, "Binary files a/%s and b/%s differ\n", f.RelPath, f.RelPath)
		}
		if f.Summary != "" {
			fmt.Fprintf(&b, "# %s\n", f.Summary)
		}
	}
	return b.String()
}
