package diffr

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writePNG(t *testing.T, path string, mark color.Color) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := range 8 {
		for x := range 8 {
			img.Set(x, y, color.RGBA{R: 20, G: 20, B: 20, A: 255})
		}
	}
	if mark != nil {
		img.Set(1, 1, mark)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

func byPath(res *Result) map[string]FileDiff {
	m := make(map[string]FileDiff, len(res.Files))
	for _, f := range res.Files {
		m[f.RelPath] = f
	}
	return m
}

func TestCompareReportsFilesOnlyInRight(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "shared.txt"), "hello\n")
	writeFile(t, filepath.Join(right, "shared.txt"), "hello\n")
	writeFile(t, filepath.Join(right, "only-right.txt"), "new\n")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	files := byPath(res)
	got, ok := files["only-right.txt"]
	if !ok {
		t.Fatalf("expected only-right.txt, got %+v", res.Files)
	}
	if got.Status != StatusAdded {
		t.Fatalf("status = %s, want added", got.Status)
	}
	if res.Stats.Added != 1 || res.Stats.Identical != 1 {
		t.Fatalf("stats = %+v", res.Stats)
	}
}

func TestCompareReportsFilesOnlyInLeft(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "gone.txt"), "bye\n")
	writeFile(t, filepath.Join(right, "kept.txt"), "ok\n")
	writeFile(t, filepath.Join(left, "kept.txt"), "ok\n")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	got, ok := byPath(res)["gone.txt"]
	if !ok {
		t.Fatalf("expected gone.txt, got %+v", res.Files)
	}
	if got.Status != StatusRemoved {
		t.Fatalf("status = %s, want removed", got.Status)
	}
}

func TestCompareTextChange(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "a.txt"), "one\n")
	writeFile(t, filepath.Join(right, "a.txt"), "two\n")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	got := byPath(res)["a.txt"]
	if got.Kind != KindText || got.Unified == "" {
		t.Fatalf("expected text unified diff, got %+v", got)
	}
	if !strings.Contains(got.Unified, "-one") || !strings.Contains(got.Unified, "+two") {
		t.Fatalf("unexpected diff:\n%s", got.Unified)
	}
}

func TestDefaultIgnoreSkipsGitAndNodeModules(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "keep.txt"), "a\n")
	writeFile(t, filepath.Join(right, "keep.txt"), "b\n")
	writeFile(t, filepath.Join(left, "node_modules", "pkg.js"), "old\n")
	writeFile(t, filepath.Join(right, "node_modules", "pkg.js"), "new\n")
	writeFile(t, filepath.Join(left, ".git", "config"), "old\n")
	writeFile(t, filepath.Join(right, ".git", "config"), "new\n")

	res, err := Compare(left, right, Options{})
	if err != nil {
		t.Fatal(err)
	}
	files := byPath(res)
	if _, ok := files["keep.txt"]; !ok {
		t.Fatal("expected keep.txt")
	}
	for _, p := range []string{"node_modules/pkg.js", ".git/config"} {
		if _, ok := files[p]; ok {
			t.Fatalf("did not expect %s", p)
		}
	}
}

func TestExcludeFlag(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "a.txt"), "a\n")
	writeFile(t, filepath.Join(right, "a.txt"), "b\n")
	writeFile(t, filepath.Join(left, "skip.log"), "a\n")
	writeFile(t, filepath.Join(right, "skip.log"), "b\n")

	res, err := Compare(left, right, Options{
		NoDefaultExclude: true,
		Exclude:          []string{"*.log"},
	})
	if err != nil {
		t.Fatal(err)
	}
	files := byPath(res)
	if _, ok := files["skip.log"]; ok {
		t.Fatal("skip.log should be excluded")
	}
	if _, ok := files["a.txt"]; !ok {
		t.Fatal("expected a.txt")
	}
}

func TestDiffrignore(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, ".diffrignore"), "secret.txt\n*.tmp\n")
	writeFile(t, filepath.Join(left, "secret.txt"), "a\n")
	writeFile(t, filepath.Join(right, "secret.txt"), "b\n")
	writeFile(t, filepath.Join(left, "x.tmp"), "a\n")
	writeFile(t, filepath.Join(right, "x.tmp"), "b\n")
	writeFile(t, filepath.Join(left, "visible.txt"), "a\n")
	writeFile(t, filepath.Join(right, "visible.txt"), "b\n")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	files := byPath(res)
	if _, ok := files["secret.txt"]; ok {
		t.Fatal("secret.txt should be ignored")
	}
	if _, ok := files["x.tmp"]; ok {
		t.Fatal("x.tmp should be ignored")
	}
	if _, ok := files["visible.txt"]; !ok {
		t.Fatal("expected visible.txt")
	}
}

func TestBinaryFilesAreNotDumpedAsText(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "a.bin"), "hello\x00old")
	writeFile(t, filepath.Join(right, "a.bin"), "hello\x00new")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	got := byPath(res)["a.bin"]
	if got.Kind != KindBinary {
		t.Fatalf("kind = %s, want binary", got.Kind)
	}
	if strings.Contains(got.Unified, "\x00") {
		t.Fatal("binary content should not appear in unified diff")
	}
	if !strings.Contains(got.Summary, "binary files differ") {
		t.Fatalf("summary = %q", got.Summary)
	}
	patch := res.Patch()
	if !strings.Contains(patch, "Binary files a/a.bin and b/a.bin differ") {
		t.Fatalf("patch = %s", patch)
	}
}

func TestImagePixelDiff(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writePNG(t, filepath.Join(left, "pic.png"), color.RGBA{R: 255, A: 255})
	writePNG(t, filepath.Join(right, "pic.png"), color.RGBA{G: 255, A: 255})

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	got := byPath(res)["pic.png"]
	if got.Kind != KindImage {
		t.Fatalf("kind = %s, want image", got.Kind)
	}
	if len(got.DiffPNG) == 0 {
		t.Fatal("expected pixel diff PNG")
	}
	if !strings.Contains(got.Summary, "pixels") {
		t.Fatalf("summary = %q", got.Summary)
	}
}

func TestIdenticalImagesAreSkipped(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writePNG(t, filepath.Join(left, "pic.png"), color.RGBA{R: 255, A: 255})
	writePNG(t, filepath.Join(right, "pic.png"), color.RGBA{R: 255, A: 255})

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Files) != 0 {
		t.Fatalf("expected no diffs, got %+v", res.Files)
	}
	if res.Stats.Identical != 1 {
		t.Fatalf("identical = %d", res.Stats.Identical)
	}
}

func TestPDFDetection(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "doc.pdf"), "%PDF-1.1 left-bytes")
	writeFile(t, filepath.Join(right, "doc.pdf"), "%PDF-1.1 right-bytes")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	got := byPath(res)["doc.pdf"]
	if got.Kind != KindPDF {
		t.Fatalf("kind = %s, want pdf", got.Kind)
	}
	if got.Status != StatusChanged {
		t.Fatalf("status = %s", got.Status)
	}
}

func TestCompareTwoFiles(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt")
	b := filepath.Join(dir, "b.txt")
	writeFile(t, a, "left\n")
	writeFile(t, b, "right\n")

	res, err := Compare(a, b, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Files) != 1 {
		t.Fatalf("got %+v", res.Files)
	}
	if res.Files[0].Kind != KindText {
		t.Fatalf("kind = %s", res.Files[0].Kind)
	}
}

func TestCompareFileAndDirectoryError(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a.txt")
	writeFile(t, file, "x\n")
	_, err := Compare(file, dir, Options{NoDefaultExclude: true})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPatchIncludesAddedRightFile(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "keep.txt"), "same\n")
	writeFile(t, filepath.Join(right, "keep.txt"), "same\n")
	writeFile(t, filepath.Join(right, "extra.txt"), "only here\n")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	patch := res.Patch()
	if !strings.Contains(patch, "+++ b/extra.txt") {
		t.Fatalf("patch missing added file:\n%s", patch)
	}
	if !strings.Contains(patch, "diff --git a/extra.txt b/extra.txt") {
		t.Fatalf("patch missing git header:\n%s", patch)
	}
	if !strings.Contains(patch, "new file mode 100644") {
		t.Fatalf("patch missing new file marker:\n%s", patch)
	}
}

func TestGitHeadersForAddedAndRemoved(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "gone.txt"), "bye\n")
	writeFile(t, filepath.Join(right, "new.txt"), "hi\n")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	files := byPath(res)
	added := files["new.txt"].Unified
	if !strings.Contains(added, "new file mode 100644") || !strings.Contains(added, "--- /dev/null") {
		t.Fatalf("added unified:\n%s", added)
	}
	removed := files["gone.txt"].Unified
	if !strings.Contains(removed, "deleted file mode 100644") || !strings.Contains(removed, "+++ /dev/null") {
		t.Fatalf("removed unified:\n%s", removed)
	}
}

func TestJSONReport(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "a.txt"), "old\n")
	writeFile(t, filepath.Join(right, "a.txt"), "new\n")
	writeFile(t, filepath.Join(right, "b.txt"), "added\n")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	out := res.JSON()
	if out.Stats.Changed != 1 || out.Stats.Added != 1 {
		t.Fatalf("stats = %+v", out.Stats)
	}
	if len(out.Files) != 2 {
		t.Fatalf("files = %+v", out.Files)
	}
	by := map[string]jsonFile{}
	for _, f := range out.Files {
		by[f.Path] = f
	}
	if by["a.txt"].Status != "changed" || by["a.txt"].Kind != "text" {
		t.Fatalf("a.txt = %+v", by["a.txt"])
	}
	if by["b.txt"].Status != "added" {
		t.Fatalf("b.txt = %+v", by["b.txt"])
	}
}

func TestIgnoreFileMissing(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "a.txt"), "a\n")
	writeFile(t, filepath.Join(right, "a.txt"), "a\n")

	_, err := Compare(left, right, Options{
		NoDefaultExclude: true,
		IgnoreFile:       filepath.Join(dir, "missing.ignore"),
	})
	if err == nil {
		t.Fatal("expected error for missing ignore file")
	}
}

func TestSniffKind(t *testing.T) {
	dir := t.TempDir()
	text := filepath.Join(dir, "a.txt")
	bin := filepath.Join(dir, "a.bin")
	pdf := filepath.Join(dir, "a.pdf")
	writeFile(t, text, "hello")
	writeFile(t, bin, "a\x00b")
	writeFile(t, pdf, "%PDF-1.4 rest")

	k, err := sniffKind(text)
	if err != nil || k != KindText {
		t.Fatalf("text: %v %v", k, err)
	}
	k, err = sniffKind(bin)
	if err != nil || k != KindBinary {
		t.Fatalf("binary: %v %v", k, err)
	}
	k, err = sniffKind(pdf)
	if err != nil || k != KindPDF {
		t.Fatalf("pdf: %v %v", k, err)
	}
}

func TestIgnoreWhitespace(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "a.txt"), "hello world\n")
	writeFile(t, filepath.Join(right, "a.txt"), "hello   world\n")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Files) != 1 {
		t.Fatalf("without -w expected a change, got %+v", res.Files)
	}

	res, err = Compare(left, right, Options{NoDefaultExclude: true, IgnoreWhitespace: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Files) != 0 || res.Stats.Identical != 1 {
		t.Fatalf("with -w expected identical, stats=%+v files=%+v", res.Stats, res.Files)
	}
}

func TestIgnoreCase(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "a.txt"), "Hello\n")
	writeFile(t, filepath.Join(right, "a.txt"), "hello\n")

	res, err := Compare(left, right, Options{NoDefaultExclude: true, IgnoreCase: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Files) != 0 || res.Stats.Identical != 1 {
		t.Fatalf("expected identical, stats=%+v files=%+v", res.Stats, res.Files)
	}
}

func TestContextLines(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, "a.txt"), "a\nb\nc\nOLD\nd\ne\nf\n")
	writeFile(t, filepath.Join(right, "a.txt"), "a\nb\nc\nNEW\nd\ne\nf\n")

	wide, err := Compare(left, right, Options{NoDefaultExclude: true, Context: 3})
	if err != nil {
		t.Fatal(err)
	}
	narrow, err := Compare(left, right, Options{NoDefaultExclude: true, Context: 1})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(wide.Files[0].Unified, "a\n") {
		t.Fatalf("context 3 should include distant line a:\n%s", wide.Files[0].Unified)
	}
	if strings.Contains(narrow.Files[0].Unified, "a\n") {
		t.Fatalf("context 1 should not include distant line a:\n%s", narrow.Files[0].Unified)
	}
}

func TestGitignoreFromGitRepo(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	writeFile(t, filepath.Join(left, ".git", "HEAD"), "ref: refs/heads/main\n")
	writeFile(t, filepath.Join(right, ".git", "HEAD"), "ref: refs/heads/main\n")
	writeFile(t, filepath.Join(left, ".gitignore"), "secret.txt\n")
	writeFile(t, filepath.Join(right, ".gitignore"), "secret.txt\n")
	writeFile(t, filepath.Join(left, "secret.txt"), "a\n")
	writeFile(t, filepath.Join(right, "secret.txt"), "b\n")
	writeFile(t, filepath.Join(left, "keep.txt"), "a\n")
	writeFile(t, filepath.Join(right, "keep.txt"), "b\n")

	res, err := Compare(left, right, Options{NoDefaultExclude: true})
	if err != nil {
		t.Fatal(err)
	}
	files := byPath(res)
	if _, ok := files["secret.txt"]; ok {
		t.Fatal("secret.txt should be gitignored")
	}
	if _, ok := files["keep.txt"]; !ok {
		t.Fatal("expected keep.txt")
	}

	res, err = Compare(left, right, Options{NoDefaultExclude: true, NoGitignore: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := byPath(res)["secret.txt"]; !ok {
		t.Fatal("secret.txt should appear with --no-gitignore")
	}
}

func TestColorizePatch(t *testing.T) {
	in := "diff --git a/f b/f\n--- a/f\n+++ b/f\n@@ -1 +1 @@\n-old\n+new\n# note\n"
	out := colorizePatch(in)
	if !strings.Contains(out, ansiRed+"-old"+ansiReset) {
		t.Fatalf("missing red deletion:\n%s", out)
	}
	if !strings.Contains(out, ansiGreen+"+new"+ansiReset) {
		t.Fatalf("missing green addition:\n%s", out)
	}
	if !strings.Contains(out, ansiCyan+"@@ -1 +1 @@"+ansiReset) {
		t.Fatalf("missing cyan hunk:\n%s", out)
	}
	if !strings.Contains(out, ansiBold+"diff --git a/f b/f"+ansiReset) {
		t.Fatalf("missing bold header:\n%s", out)
	}
}
