package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

type OpenAIProvider struct {
	APIKey  string
	Model   string
	Quality string // low | medium | high，medium 更快
	MaxLong int    // 送 API 前长边上限，0=不限制
}

func (o *OpenAIProvider) Name() string { return "openai" }

func (o *OpenAIProvider) Inpaint(ctx context.Context, in InpaintInput) (*InpaintResult, error) {
	if o.APIKey == "" {
		return nil, fmt.Errorf("openai: 未配置 openai_api_key")
	}

	prompt := strings.TrimSpace(in.Prompt)
	if prompt == "" {
		prompt = "Remove the masked objects and fill the area with a natural, seamless, photorealistic background that matches the surrounding scene."
	}

	src, err := png.Decode(bytes.NewReader(in.ImagePNG))
	if err != nil {
		return nil, fmt.Errorf("openai: 无效图片: %w", err)
	}
	b := src.Bounds()
	origW, origH := b.Dx(), b.Dy()

	quality := o.Quality
	if quality == "" {
		quality = "medium"
	}
	maxLong := o.MaxLong
	if maxLong <= 0 {
		maxLong = 1024
	}

	var imgForAPI, maskForAPI []byte
	sizeTag := "auto"
	if strings.HasPrefix(o.Model, "gpt-image-2") {
		var err error
		imgForAPI, maskForAPI, _, _, err = PrepareForOpenAIMaxLong(in.ImagePNG, in.MaskPNG, maxLong)
		if err != nil {
			return nil, err
		}
	} else {
		sizeTag = PickOpenAISize(origW, origH)
		var err error
		imgForAPI, maskForAPI, _, _, _, _, err = PrepareForOpenAI(in.ImagePNG, in.MaskPNG, sizeTag)
		if err != nil {
			return nil, err
		}
	}

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	_ = w.WriteField("model", o.Model)
	_ = w.WriteField("prompt", prompt)
	_ = w.WriteField("quality", quality)
	_ = w.WriteField("size", sizeTag)

	if err := writeMultipartPNG(w, "image", "image.png", imgForAPI); err != nil {
		return nil, err
	}
	if len(maskForAPI) > 0 {
		if err := writeMultipartPNG(w, "mask", "mask.png", maskForAPI); err != nil {
			return nil, err
		}
	}
	_ = w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/images/edits", &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+o.APIKey)
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 180 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai 请求失败: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, friendlyOpenAIHTTPError(resp.StatusCode, raw)
	}

	imgPNG, err := parseOpenAIImageResponse(raw)
	if err != nil {
		return nil, err
	}

	// 缩放回画布原始尺寸，便于 PSD 图层对齐
	outW, outH := origW, origH
	if outW > 0 && outH > 0 {
		decoded, decErr := png.Decode(bytes.NewReader(imgPNG))
		if decErr == nil {
			b2 := decoded.Bounds()
			if b2.Dx() != outW || b2.Dy() != outH {
				imgPNG, err = ResizePNG(imgPNG, outW, outH)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &InpaintResult{
		ImagePNG: imgPNG,
		Provider: "openai:" + o.Model,
		Message:  fmt.Sprintf("GPT Image 消除完成 (quality=%s, %dx%d)", quality, outW, outH),
	}, nil
}

func parseOpenAIImageResponse(raw []byte) ([]byte, error) {
	var parsed struct {
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
		Data []struct {
			B64JSON string `json:"b64_json"`
			URL     string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return nil, fmt.Errorf("openai: %s", parsed.Error.Message)
	}
	if len(parsed.Data) == 0 {
		return nil, fmt.Errorf("openai: 响应无图片数据")
	}
	if parsed.Data[0].B64JSON != "" {
		return base64.StdEncoding.DecodeString(parsed.Data[0].B64JSON)
	}
	if parsed.Data[0].URL != "" {
		r2, err := http.Get(parsed.Data[0].URL)
		if err != nil {
			return nil, err
		}
		defer r2.Body.Close()
		return io.ReadAll(r2.Body)
	}
	return nil, fmt.Errorf("openai: 无 b64_json 或 url")
}

// writeMultipartPNG OpenAI 要求文件 part 带 image/png，不能用默认 octet-stream
func writeMultipartPNG(w *multipart.Writer, field, filename string, data []byte) error {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, field, filename))
	h.Set("Content-Type", "image/png")
	part, err := w.CreatePart(h)
	if err != nil {
		return err
	}
	_, err = part.Write(data)
	return err
}

func truncateErr(raw []byte) string {
	s := string(raw)
	if len(s) > 400 {
		return s[:400] + "..."
	}
	return s
}
