package ai

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

// MockProvider 本地模拟消除：从蒙版边缘向内传播邻域颜色（Demo 比整图模糊更自然）
type MockProvider struct{}

func (m *MockProvider) Name() string { return "mock" }

func (m *MockProvider) Inpaint(ctx context.Context, in InpaintInput) (*InpaintResult, error) {
	src, err := png.Decode(bytes.NewReader(in.ImagePNG))
	if err != nil {
		return nil, err
	}
	brush, err := BrushMaskToRGBA(in.MaskPNG)
	if err != nil {
		return nil, err
	}
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	brush = alignMaskSize(brush, w, h)

	out := image.NewRGBA(b)
	draw.Draw(out, b, src, b.Min, draw.Src)

	edit := make([]bool, w*h)
	editCount := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if pixelIsMask(brush.RGBAAt(x, y)) {
				edit[y*w+x] = true
				editCount++
			}
		}
	}
	if editCount == 0 {
		return encodeMock(out, "未识别到蒙版，请涂粗红色区域")
	}

	ratio := float64(editCount) / float64(w*h)
	if ratio > 0.35 {
		return encodeMock(out, "涂抹范围过大（超过35%），请只涂刀叉等小物体")
	}

	// 从蒙版边界向内逐层填充（不混入被删物体颜色）
	filled := propagateFromBorder(out, edit, w, h, b.Min.X, b.Min.Y)

	// 蒙版区域做轻度模糊，抹平传播纹理
	radius := 3
	if w > 600 || h > 600 {
		radius = 5
	}
	smooth := boxBlur(filled, radius)

	changed := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !edit[y*w+x] {
				continue
			}
			px, py := b.Min.X+x, b.Min.Y+y
			oc := out.RGBAAt(px, py)
			nc := smooth.RGBAAt(px, py)
			// 仅在最外 4px 与保留区做羽化，避免刀叉色回流
			edge := distToKeep(edit, x, y, w, h)
			if edge > 0 && edge < 4 {
				t := float64(edge) / 4
				nc = color.RGBA{
					R: lerpU8(oc.R, nc.R, t),
					G: lerpU8(oc.G, nc.G, t),
					B: lerpU8(oc.B, nc.B, t),
					A: 255,
				}
			}
			out.SetRGBA(px, py, nc)
			if nc != oc {
				changed++
			}
		}
	}

	msg := "mock 边缘传播填充完成（已修改 " + itoa(changed) + " 像素）"
	if ratio > 0.12 {
		msg += "。范围偏大时可能发糊，请只涂要删的物体"
	}
	return encodeMock(out, msg)
}

// propagateFromBorder 已知区=非蒙版，逐轮用 8 邻域已知像素均值填充蒙版
func propagateFromBorder(src *image.RGBA, edit []bool, w, h, ox, oy int) *image.RGBA {
	out := image.NewRGBA(src.Bounds())
	draw.Draw(out, src.Bounds(), src, src.Bounds().Min, draw.Src)

	known := make([]bool, w*h)
	for i := range known {
		known[i] = !edit[i]
	}

	maxPass := w + h
	for pass := 0; pass < maxPass; pass++ {
		progress := false
		next := append([]bool(nil), known...)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				idx := y*w + x
				if known[idx] || !edit[idx] {
					continue
				}
				var sr, sg, sb float64
				var n int
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						if dx == 0 && dy == 0 {
							continue
						}
						nx, ny := x+dx, y+dy
						if nx < 0 || ny < 0 || nx >= w || ny >= h {
							continue
						}
						nidx := ny*w + nx
						if !known[nidx] {
							continue
						}
						c := out.RGBAAt(ox+nx, oy+ny)
						sr += float64(c.R)
						sg += float64(c.G)
						sb += float64(c.B)
						n++
					}
				}
				if n == 0 {
					continue
				}
				out.SetRGBA(ox+x, oy+y, color.RGBA{
					R: uint8(sr / float64(n)),
					G: uint8(sg / float64(n)),
					B: uint8(sb / float64(n)),
					A: 255,
				})
				next[idx] = true
				progress = true
			}
		}
		known = next
		if !progress {
			break
		}
	}
	return out
}

// distToKeep 到最近非蒙版像素的步数
func distToKeep(edit []bool, x, y, w, h int) int {
	if !edit[y*w+x] {
		return 0
	}
	for d := 1; d <= 8; d++ {
		for dy := -d; dy <= d; dy++ {
			for dx := -d; dx <= d; dx++ {
				nx, ny := x+dx, y+dy
				if nx < 0 || ny < 0 || nx >= w || ny >= h {
					continue
				}
				if !edit[ny*w+nx] {
					return d
				}
			}
		}
	}
	return 8
}

func boxBlur(src *image.RGBA, radius int) *image.RGBA {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	tmp := image.NewRGBA(b)
	out := image.NewRGBA(b)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var sr, sg, sb float64
			var n int
			for dx := -radius; dx <= radius; dx++ {
				nx := x + dx
				if nx < 0 || nx >= w {
					continue
				}
				c := src.RGBAAt(b.Min.X+nx, b.Min.Y+y)
				sr += float64(c.R)
				sg += float64(c.G)
				sb += float64(c.B)
				n++
			}
			tmp.SetRGBA(b.Min.X+x, b.Min.Y+y, color.RGBA{
				R: uint8(sr / float64(n)),
				G: uint8(sg / float64(n)),
				B: uint8(sb / float64(n)),
				A: 255,
			})
		}
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var sr, sg, sb float64
			var n int
			for dy := -radius; dy <= radius; dy++ {
				ny := y + dy
				if ny < 0 || ny >= h {
					continue
				}
				c := tmp.RGBAAt(b.Min.X+x, b.Min.Y+ny)
				sr += float64(c.R)
				sg += float64(c.G)
				sb += float64(c.B)
				n++
			}
			out.SetRGBA(b.Min.X+x, b.Min.Y+y, color.RGBA{
				R: uint8(sr / float64(n)),
				G: uint8(sg / float64(n)),
				B: uint8(sb / float64(n)),
				A: 255,
			})
		}
	}
	return out
}

func lerpU8(a, b uint8, t float64) uint8 {
	return uint8(float64(a)*(1-t) + float64(b)*t)
}

func pixelIsMask(c color.RGBA) bool {
	return c.R > 180 || (c.A > 80 && c.R > c.B)
}

func alignMaskSize(brush *image.RGBA, w, h int) *image.RGBA {
	bw, bh := brush.Bounds().Dx(), brush.Bounds().Dy()
	if bw == w && bh == h {
		return brush
	}
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sx := x * bw / w
			sy := y * bh / h
			if sx >= bw {
				sx = bw - 1
			}
			if sy >= bh {
				sy = bh - 1
			}
			out.Set(x, y, brush.RGBAAt(sx, sy))
		}
	}
	return out
}

func encodeMock(out *image.RGBA, msg string) (*InpaintResult, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, out); err != nil {
		return nil, err
	}
	return &InpaintResult{ImagePNG: buf.Bytes(), Provider: "mock", Message: msg}, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var d [20]byte
	i := len(d)
	for n > 0 {
		i--
		d[i] = byte('0' + n%10)
		n /= 10
	}
	return string(d[i:])
}
