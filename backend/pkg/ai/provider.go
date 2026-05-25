package ai

import "context"

type InpaintInput struct {
	ImagePNG []byte
	MaskPNG  []byte // 可选；空则整图编辑
	Prompt   string
	Mode     StudioMode
}

type InpaintResult struct {
	ImagePNG []byte
	Provider string `json:"provider"`
	Message  string `json:"message,omitempty"`
}

type Provider interface {
	Inpaint(ctx context.Context, in InpaintInput) (*InpaintResult, error)
	Name() string
}

func NewProvider(mode, apiKey, model, quality string, maxLong int) Provider {
	if mode == "openai" && apiKey != "" {
		if model == "" {
			model = "gpt-image-2"
		}
		return &OpenAIProvider{APIKey: apiKey, Model: model, Quality: quality, MaxLong: maxLong}
	}
	return &MockProvider{}
}
