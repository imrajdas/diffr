package diffr

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

type pdfInfo struct {
	pages int
	text  string
	err   error
}

func readPDF(path string) pdfInfo {
	if path == "" {
		return pdfInfo{}
	}
	f, r, err := pdf.Open(path)
	if err != nil {
		return pdfInfo{err: err}
	}
	defer f.Close()

	info := pdfInfo{pages: r.NumPage()}
	reader, err := r.GetPlainText()
	if err != nil {
		info.err = err
		return info
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		info.err = err
		return info
	}
	info.text = buf.String()
	return info
}

func comparePDFs(fd FileDiff, opts Options) (FileDiff, bool, error) {
	fd.Kind = KindPDF
	left := readPDF(fd.LeftPath)
	right := readPDF(fd.RightPath)

	if (fd.LeftPath != "" && left.err != nil) || (fd.RightPath != "" && right.err != nil) {
		bd, ident, err := compareBinaries(fd)
		if err != nil {
			return fd, false, err
		}
		bd.Kind = KindPDF
		if ident {
			return bd, true, nil
		}
		var parts []string
		if left.err != nil {
			parts = append(parts, "could not parse left PDF: "+left.err.Error())
		}
		if right.err != nil {
			parts = append(parts, "could not parse right PDF: "+right.err.Error())
		}
		if bd.Summary != "" {
			parts = append(parts, bd.Summary)
		}
		bd.Summary = strings.Join(parts, "; ")
		return bd, false, nil
	}

	if left.pages == right.pages && equivalentText(left.text, right.text, opts) {
		if left.text == right.text && fd.LeftPath != "" && fd.RightPath != "" {
			same, err := sameHash(fd.LeftPath, fd.RightPath)
			if err != nil {
				return fd, false, err
			}
			if same {
				return fd, true, nil
			}
			fd.Summary = fmt.Sprintf("PDFs differ (same text, %d pages; file bytes differ)", left.pages)
			return fd, false, nil
		}
		return fd, true, nil
	}

	var summary []string
	switch fd.Status {
	case StatusAdded:
		summary = append(summary, fmt.Sprintf("PDF added (%d pages)", right.pages))
	case StatusRemoved:
		summary = append(summary, fmt.Sprintf("PDF removed (%d pages)", left.pages))
	default:
		if left.pages != right.pages {
			summary = append(summary, fmt.Sprintf("page count %d vs %d", left.pages, right.pages))
		}
		if left.text != right.text && !equivalentText(left.text, right.text, opts) {
			summary = append(summary, "text content differs")
		}
	}
	fd.Summary = "PDF: " + strings.Join(summary, "; ")

	if !equivalentText(left.text, right.text, opts) {
		unified, err := makeUnified(fd.RelPath, left.text, right.text, fd.Status, opts.Context)
		if err != nil {
			return fd, false, err
		}
		fd.Unified = unified
	}
	return fd, false, nil
}
