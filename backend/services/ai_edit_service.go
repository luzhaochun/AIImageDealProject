package services

import (
	"bytes"
	"context"
	"image/png"

	"backend/pkg/apperr"
	"os"
	"path/filepath"
	"time"

	"backend/pkg/ai"
	"backend/pkg/storage"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
)

type AiEditService struct{}

func newOpenAIProvider(key, model string) *ai.OpenAIProvider {
	if model == "" {
		model = "gpt-image-2"
	}
	quality, _ := beego.AppConfig.String("openai_image_quality")
	if quality == "" {
		quality = "medium"
	}
	maxLong, _ := beego.AppConfig.Int("openai_image_max_long")
	if maxLong <= 0 {
		maxLong = 1024
	}
	return &ai.OpenAIProvider{APIKey: key, Model: model, Quality: quality, MaxLong: maxLong}
}

func openAIProviderOpts() (quality string, maxLong int) {
	quality, _ = beego.AppConfig.String("openai_image_quality")
	if quality == "" {
		quality = "medium"
	}
	maxLong, _ = beego.AppConfig.Int("openai_image_max_long")
	if maxLong <= 0 {
		maxLong = 1024
	}
	return
}

type AiInpaintResponse struct {
	ResultURL string `json:"result_url"`
	Provider  string `json:"provider"`
	Message   string `json:"message,omitempty"`
}

// StudioEdit Canvas 工作室：多模式 GPT Image（OpenAI）
func (s *AiEditService) StudioEdit(imageData, maskData []byte, prompt, modeStr string) (*AiInpaintResponse, error) {
	mode, ok := ai.ParseStudioMode(modeStr)
	if !ok {
		return nil, apperr.New("VALIDATION_ERROR", "无效的 mode", apperr.ErrValidation)
	}

	hasMask := ai.BrushMaskHasEdit(maskData)
	if ai.StudioModeRequiresMask(mode) && !hasMask {
		meta := ai.AllStudioModes()
		for _, m := range meta {
			if m.ID == string(mode) {
				return nil, apperr.New("VALIDATION_ERROR", m.MaskHint, apperr.ErrValidation)
			}
		}
		return nil, apperr.New("VALIDATION_ERROR", "请用画笔涂抹蒙版区域", apperr.ErrValidation)
	}

	key, _ := beego.AppConfig.String("openai_api_key")
	if key == "" {
		return nil, apperr.New(
			"CONFIG_ERROR",
			"请配置 GPT Image：在 backend/conf/app.conf 设置 ai_edit_mode=openai 并填写 openai_api_key",
			apperr.ErrValidation,
		)
	}
	model, _ := beego.AppConfig.String("openai_image_model")
	if model == "" {
		model = "gpt-image-2"
	}

	var maskForAPI []byte
	if hasMask {
		src, err := png.Decode(bytes.NewReader(imageData))
		if err != nil {
			return nil, err
		}
		rgba, err := ai.BrushMaskToRGBA(maskData)
		if err != nil {
			return nil, err
		}
		openaiMask, err := ai.BuildOpenAIMask(src, rgba)
		if err != nil {
			return nil, err
		}
		maskForAPI = openaiMask
	}

	fullPrompt := ai.BuildStudioPrompt(mode, prompt)
	provider := newOpenAIProvider(key, model)
	res, err := provider.Inpaint(context.Background(), ai.InpaintInput{
		ImagePNG: imageData,
		MaskPNG:  maskForAPI,
		Prompt:   fullPrompt,
		Mode:     mode,
	})
	if err != nil {
		return nil, err
	}
	return s.saveInpaintResult(res)
}

// StudioInpaint 兼容旧接口，等同 erase 模式
func (s *AiEditService) StudioInpaint(imageData, maskData []byte, prompt string) (*AiInpaintResponse, error) {
	return s.StudioEdit(imageData, maskData, prompt, string(ai.ModeErase))
}

func (s *AiEditService) Inpaint(imageData, maskData []byte, prompt string) (*AiInpaintResponse, error) {
	if !ai.BrushMaskHasEdit(maskData) {
		return nil, apperr.New("VALIDATION_ERROR", "请用画笔涂抹要消除的区域", apperr.ErrValidation)
	}

	mode, _ := beego.AppConfig.String("ai_edit_mode")
	key, _ := beego.AppConfig.String("openai_api_key")
	model, _ := beego.AppConfig.String("openai_image_model")
	quality, maxLong := openAIProviderOpts()
	provider := ai.NewProvider(mode, key, model, quality, maxLong)

	maskForAPI := maskData
	if mode == "openai" && key != "" {
		src, err := png.Decode(bytes.NewReader(imageData))
		if err != nil {
			return nil, err
		}
		rgba, err := ai.BrushMaskToRGBA(maskData)
		if err != nil {
			return nil, err
		}
		openaiMask, err := ai.BuildOpenAIMask(src, rgba)
		if err != nil {
			return nil, err
		}
		maskForAPI = openaiMask
	}

	res, err := provider.Inpaint(context.Background(), ai.InpaintInput{
		ImagePNG: imageData,
		MaskPNG:  maskForAPI,
		Prompt:   prompt,
	})
	if err != nil {
		return nil, err
	}
	return s.saveInpaintResult(res)
}

func (s *AiEditService) saveInpaintResult(res *ai.InpaintResult) (*AiInpaintResponse, error) {
	id := uuid.New().String()
	rel := filepath.Join("ai-demo", id+".png")
	full := storage.FullPath(rel)
	if err := storage.EnsureDir(full); err != nil {
		return nil, err
	}
	if err := os.WriteFile(full, res.ImagePNG, 0644); err != nil {
		return nil, err
	}
	return &AiInpaintResponse{
		ResultURL: "/static/" + rel + "?t=" + time.Now().Format("150405"),
		Provider:  res.Provider,
		Message:   res.Message,
	}, nil
}
