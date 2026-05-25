package ai

import (
	"bytes"
	"image"
	"image/png"
)

// BrushMaskHasEdit 红色画笔蒙版是否包含涂抹区域
func BrushMaskHasEdit(maskPNG []byte) bool {
	img, err := png.Decode(bytes.NewReader(maskPNG))
	if err != nil {
		return false
	}
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0x3000 {
				return true
			}
			r, _, _, _ := img.At(x, y).RGBA()
			if r > 0xc000 {
				return true
			}
		}
	}
	return false
}

func BrushMaskToRGBA(maskPNG []byte) (*image.RGBA, error) {
	img, err := png.Decode(bytes.NewReader(maskPNG))
	if err != nil {
		return nil, err
	}
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba, nil
	}
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba, nil
}
