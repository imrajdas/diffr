package diffr

import (
	"html/template"
	"strings"
	"testing"

	"github.com/imrajdas/diffr/static"
)

func TestPageDataOmitsPDFFromTextDiff(t *testing.T) {
	res := &Result{
		Files: []FileDiff{
			{RelPath: "a.txt", Kind: KindText, Status: StatusChanged, Unified: "diff --git a/a.txt b/a.txt\n"},
			{RelPath: "docs/notes.pdf", Kind: KindPDF, Status: StatusChanged, Unified: "diff --git a/docs/notes.pdf b/docs/notes.pdf\n-Draft\n+Final\n", Summary: "PDF: text content differs"},
		},
		Stats: Stats{Changed: 2},
	}
	pd := pageData(res)
	if strings.Contains(pd.Diff, "notes.pdf") {
		t.Fatalf("PDF should not be in the main text diff:\n%s", pd.Diff)
	}
	if !strings.Contains(pd.Diff, "a.txt") {
		t.Fatal("text file should remain in the main diff")
	}
	if len(pd.PDFs) != 1 || pd.PDFs[0].Unified == "" {
		t.Fatalf("PDF text diff should stay on the PDF row: %+v", pd.PDFs)
	}
}

func TestHTMLTemplateParses(t *testing.T) {
	if _, err := template.New("html").Parse(static.HTML); err != nil {
		t.Fatal(err)
	}
}
