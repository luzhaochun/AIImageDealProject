package ai

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestBuildOpenAIMask_TransparentEditArea(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 4, 4))
	brush := image.NewRGBA(image.Rect(0, 0, 4, 4))
	brush.Set(1, 1, color.RGBA{255, 40, 40, 255})

	maskPNG, err := BuildOpenAIMask(src, brush)
	if err != nil {
		t.Fatal(err)
	}
	m, err := png.Decode(bytes.NewReader(maskPNG))
	if err != nil {
		t.Fatal(err)
	}
	_, _, _, aEdit := m.At(1, 1).RGBA()
	if aEdit != 0 {
		t.Fatalf("brush pixel alpha want 0 (edit), got %d", aEdit)
	}
	_, _, _, aKeep := m.At(0, 0).RGBA()
	if aKeep == 0 {
		t.Fatal("non-brush pixel should be opaque (keep)")
	}
}
