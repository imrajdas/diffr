package diffr

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

type Kind int

const (
	KindText Kind = iota
	KindBinary
	KindImage
	KindPDF
)

func (k Kind) String() string {
	switch k {
	case KindBinary:
		return "binary"
	case KindImage:
		return "image"
	case KindPDF:
		return "pdf"
	default:
		return "text"
	}
}

const sniffSize = 8192

func sniffKind(path string) (Kind, error) {
	f, err := os.Open(path)
	if err != nil {
		return KindText, err
	}
	defer f.Close()

	buf := make([]byte, sniffSize)
	n, err := f.Read(buf)
	if n == 0 && err != nil && err != io.EOF {
		return KindText, err
	}
	head := buf[:n]

	if bytes.HasPrefix(head, []byte("%PDF")) {
		return KindPDF, nil
	}

	switch http.DetectContentType(head) {
	case "image/png", "image/jpeg", "image/gif", "image/webp", "image/bmp":
		return KindImage, nil
	}

	if bytes.IndexByte(head, 0) >= 0 {
		return KindBinary, nil
	}
	return KindText, nil
}

func mergeKind(left, right Kind, hasLeft, hasRight bool) Kind {
	if !hasLeft {
		return right
	}
	if !hasRight {
		return left
	}
	if left == right {
		return left
	}
	if left != KindText && right != KindText {
		return KindBinary
	}
	if left != KindText {
		return left
	}
	return right
}
