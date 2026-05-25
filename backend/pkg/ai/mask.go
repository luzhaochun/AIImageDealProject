package ai

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// BuildOpenAIMask 将前端红色涂抹蒙版转为 OpenAI 格式：
// 透明(alpha=0) = 需要编辑/消除的区域；不透明 = 保留原图
func BuildOpenAIMask(src image.Image, brushMask *image.RGBA) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	brushMask = alignMaskSize(brushMask, w, h)

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	// 默认全不透明 = 保留整张图
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			out.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	isBrush := func(x, y int) bool {
		if x < 0 || y < 0 || x >= w || y >= h {
			return false
		}
		c := brushMask.RGBAAt(x, y)
		return pixelIsMask(c)
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if isBrush(x, y) {
				// 透明区域 = API 编辑区
				out.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, out); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
