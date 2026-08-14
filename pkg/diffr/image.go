package diffr

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	_ "image/gif"
	_ "image/jpeg"
)

const maxPixels = 20_000_000

func compareImages(fd FileDiff) (FileDiff, bool, error) {
	fd.Kind = KindImage

	if fd.LeftPath == "" || fd.RightPath == "" {
		if fd.Status == StatusAdded {
			fd.Summary = "image added"
		} else {
			fd.Summary = "image removed"
		}
		return fd, false, nil
	}

	left, err := decodeImage(fd.LeftPath)
	if err != nil {
		return compareBinaries(fd)
	}
	right, err := decodeImage(fd.RightPath)
	if err != nil {
		return compareBinaries(fd)
	}

	lb, rb := left.Bounds(), right.Bounds()
	if lb.Dx()*lb.Dy() > maxPixels || rb.Dx()*rb.Dy() > maxPixels {
		same, err := sameHash(fd.LeftPath, fd.RightPath)
		if err != nil {
			return fd, false, err
		}
		if same {
			return fd, true, nil
		}
		fd.Summary = fmt.Sprintf("images differ (%dx%d vs %dx%d; too large for pixel diff)", lb.Dx(), lb.Dy(), rb.Dx(), rb.Dy())
		return fd, false, nil
	}

	diffImg, changed, total := pixelDiff(left, right)
	if lb.Dx() == rb.Dx() && lb.Dy() == rb.Dy() && changed == 0 {
		return fd, true, nil
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, diffImg); err != nil {
		return fd, false, err
	}
	fd.DiffPNG = buf.Bytes()

	if lb.Dx() != rb.Dx() || lb.Dy() != rb.Dy() {
		fd.Summary = fmt.Sprintf("images differ in size: %dx%d vs %dx%d", lb.Dx(), lb.Dy(), rb.Dx(), rb.Dy())
		return fd, false, nil
	}

	pct := 100 * float64(changed) / float64(total)
	fd.Summary = fmt.Sprintf("images differ: %d/%d pixels (%.2f%%)", changed, total, pct)
	return fd, false, nil
}

func decodeImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func pixelDiff(a, b image.Image) (*image.RGBA, int, int) {
	ab, bb := a.Bounds(), b.Bounds()
	w := max(ab.Dx(), bb.Dx())
	h := max(ab.Dy(), bb.Dy())
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	changed := 0
	total := w * h

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			inA := x < ab.Dx() && y < ab.Dy()
			inB := x < bb.Dx() && y < bb.Dy()
			if !inA || !inB {
				changed++
				out.Set(x, y, color.RGBA{R: 37, G: 99, B: 235, A: 255})
				continue
			}
			r1, g1, b1, a1 := a.At(ab.Min.X+x, ab.Min.Y+y).RGBA()
			r2, g2, b2, a2 := b.At(bb.Min.X+x, bb.Min.Y+y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				changed++
				out.Set(x, y, color.RGBA{R: 220, G: 38, B: 38, A: 255})
				continue
			}
			gray := uint8(((r1 + g1 + b1) / 3) >> 8)
			out.Set(x, y, color.RGBA{R: gray, G: gray, B: gray, A: 255})
		}
	}
	return out, changed, total
}
