//go:build ignore

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if filepath.Base(root) != "demo" {
		root = filepath.Join(root, "testdata", "demo")
	}
	left := filepath.Join(root, "left")
	right := filepath.Join(root, "right")
	_ = os.RemoveAll(left)
	_ = os.RemoveAll(right)

	write := func(path, content string) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			panic(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			panic(err)
		}
	}

	write(filepath.Join(left, "README.md"), "# Demo App\n\nVersion 1.0\n\nA sample project for diffr.\n")
	write(filepath.Join(right, "README.md"), "# Demo App\n\nVersion 2.0\n\nA sample project for diffr, now with extra features.\n")

	write(filepath.Join(left, "src", "app.go"), "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n")
	write(filepath.Join(right, "src", "app.go"), "package main\n\nfunc main() {\n\tprintln(\"hello, world\")\n}\n")

	util := "package main\n\nfunc add(a, b int) int {\n\treturn a + b\n}\n"
	write(filepath.Join(left, "src", "util.go"), util)
	write(filepath.Join(right, "src", "util.go"), util)

	write(filepath.Join(left, "deprecated.txt"), "This file only exists on the left tree.\n")
	write(filepath.Join(right, "changelog.txt"), "This file only exists on the right tree.\n- Added extra.txt\n")

	write(filepath.Join(left, "data", "config.json"), "{\n  \"port\": 8080\n}\n")
	write(filepath.Join(right, "data", "config.json"), "{\n  \"port\": 9090,\n  \"debug\": true\n}\n")

	write(filepath.Join(left, "data", "blob.bin"), "BIN\x00LEFT\xff\xfe")
	write(filepath.Join(right, "data", "blob.bin"), "BIN\x00RIGHT\xff\xfe")

	write(filepath.Join(left, "noise.log"), "old log line\n")
	write(filepath.Join(right, "noise.log"), "new log line\n")
	write(filepath.Join(left, ".diffrignore"), "*.log\n")
	write(filepath.Join(right, ".diffrignore"), "*.log\n")

	write(filepath.Join(left, "node_modules", "left-pad", "index.js"), "module.exports = function () { return 'old' }\n")
	write(filepath.Join(right, "node_modules", "left-pad", "index.js"), "module.exports = function () { return 'new' }\n")
	write(filepath.Join(left, ".git", "HEAD"), "ref: refs/heads/main\n")
	write(filepath.Join(right, ".git", "HEAD"), "ref: refs/heads/feature\n")

	writePNG(filepath.Join(left, "assets", "logo.png"), color.RGBA{R: 30, G: 90, B: 200, A: 255})
	writePNG(filepath.Join(right, "assets", "logo.png"), color.RGBA{R: 200, G: 40, B: 40, A: 255})

	if err := os.MkdirAll(filepath.Join(left, "docs"), 0o755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Join(right, "docs"), 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(left, "docs", "notes.pdf"), miniPDF("Draft notes"), 0o644); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(right, "docs", "notes.pdf"), miniPDF("Final notes"), 0o644); err != nil {
		panic(err)
	}
}

func writePNG(path string, mark color.Color) {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	bg := color.RGBA{R: 245, G: 245, B: 245, A: 255}
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, bg)
		}
	}
	for y := 16; y < 48; y++ {
		for x := 16; x < 48; x++ {
			img.Set(x, y, mark)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		panic(err)
	}
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func miniPDF(text string) []byte {
	stream := "BT /F1 18 Tf 50 700 Td (" + text + ") Tj ET\n"
	objs := []string{
		"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n",
		"2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n",
		"3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>\nendobj\n",
		"4 0 obj\n<< /Length " + strconv.Itoa(len(stream)) + " >>\nstream\n" + stream + "endstream\nendobj\n",
		"5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n",
	}

	header := "%PDF-1.4\n"
	offsets := make([]int, len(objs))
	pos := len(header)
	body := header
	for i, obj := range objs {
		offsets[i] = pos
		body += obj
		pos += len(obj)
	}

	xrefStart := pos
	xref := "xref\n0 6\n0000000000 65535 f \n"
	for _, off := range offsets {
		xref += fmt.Sprintf("%010d 00000 n \n", off)
	}
	trailer := fmt.Sprintf("trailer\n<< /Size 6 /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", xrefStart)
	return []byte(body + xref + trailer)
}
