package ai

import (
	"bytes"
	"image"
	"image/png"
)

// PickOpenAISize 选择 GPT Image 支持的输出尺寸
func PickOpenAISize(w, h int) string {
	if w <= 0 || h <= 0 {
		return "1024x1024"
	}
	ratio := float64(w) / float64(h)
	if ratio > 1.15 {
		return "1536x1024"
	}
	if ratio < 0.87 {
		return "1024x1536"
	}
	return "1024x1024"
}

func parseSize(s string) (int, int) {
	switch s {
	case "1536x1024":
		return 1536, 1024
	case "1024x1536":
		return 1024, 1536
	default:
		return 1024, 1024
	}
}

// ResizePNG 将 PNG 缩放到目标宽高
func ResizePNG(pngData []byte, tw, th int) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return nil, err
	}
	scaled := scaleImage(img, tw, th)
	var buf bytes.Buffer
	if err := png.Encode(&buf, scaled); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func scaleImage(src image.Image, tw, th int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, tw, th))
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()
	if sw == 0 || sh == 0 {
		return dst
	}
	for y := 0; y < th; y++ {
		for x := 0; x < tw; x++ {
			sx := sb.Min.X + x*sw/tw
			sy := sb.Min.Y + y*sh/th
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}

// PrepareForOpenAIMaxLong 长边超过 maxLong 时等比缩小（加快 gpt-image-2 推理）
func PrepareForOpenAIMaxLong(imagePNG, maskPNG []byte, maxLong int) (imgOut, maskOut []byte, origW, origH int, err error) {
	src, err := png.Decode(bytes.NewReader(imagePNG))
	if err != nil {
		return nil, nil, 0, 0, err
	}
	b := src.Bounds()
	origW, origH = b.Dx(), b.Dy()
	long := origW
	if origH > long {
		long = origH
	}
	if long <= maxLong {
		return imagePNG, maskPNG, origW, origH, nil
	}
	scale := float64(maxLong) / float64(long)
	tw := int(float64(origW) * scale)
	th := int(float64(origH) * scale)
	if tw < 1 {
		tw = 1
	}
	if th < 1 {
		th = 1
	}
	imgOut, err = ResizePNG(imagePNG, tw, th)
	if err != nil {
		return nil, nil, 0, 0, err
	}
	if len(maskPNG) > 0 {
		maskOut, err = ResizePNG(maskPNG, tw, th)
		if err != nil {
			return nil, nil, 0, 0, err
		}
	}
	return imgOut, maskOut, origW, origH, nil
}

// PrepareForOpenAI 将原图与蒙版缩放到 API size
func PrepareForOpenAI(imagePNG, maskPNG []byte, sizeTag string) (imgOut, maskOut []byte, apiW, apiH, origW, origH int, err error) {
	src, err := png.Decode(bytes.NewReader(imagePNG))
	if err != nil {
		return nil, nil, 0, 0, 0, 0, err
	}
	b := src.Bounds()
	origW, origH = b.Dx(), b.Dy()
	apiW, apiH = parseSize(sizeTag)
	imgOut, err = ResizePNG(imagePNG, apiW, apiH)
	if err != nil {
		return nil, nil, 0, 0, 0, 0, err
	}
	if len(maskPNG) > 0 {
		maskOut, err = ResizePNG(maskPNG, apiW, apiH)
		if err != nil {
			return nil, nil, 0, 0, 0, 0, err
		}
	}
	return imgOut, maskOut, apiW, apiH, origW, origH, nil
}
