package ai

import (
	"strconv"
	"strings"
)

// StudioMode Canvas 工作室 AI 能力
type StudioMode string

const (
	ModeErase      StudioMode = "erase"
	ModeUpscale    StudioMode = "upscale"
	ModeRepair     StudioMode = "repair"
	ModeBackground StudioMode = "background"
	ModeOutpaint   StudioMode = "outpaint"
	ModeWatermark  StudioMode = "watermark"
)

var validStudioModes = map[StudioMode]bool{
	ModeErase: true, ModeUpscale: true, ModeRepair: true,
	ModeBackground: true, ModeOutpaint: true, ModeWatermark: true,
}

func ParseStudioMode(s string) (StudioMode, bool) {
	m := StudioMode(strings.TrimSpace(s))
	if m == "" {
		m = ModeErase
	}
	return m, validStudioModes[m]
}

// StudioModeRequiresMask 是否必须涂抹蒙版
func StudioModeRequiresMask(m StudioMode) bool {
	return m == ModeErase
}

// StudioModeInfo 供 API / 前端展示
type StudioModeInfo struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	RequiresMask bool   `json:"requires_mask"`
	MaskHint     string `json:"mask_hint"`
	DefaultPrompt string `json:"default_prompt"`
}

func AllStudioModes() []StudioModeInfo {
	out := make([]StudioModeInfo, 0, len(studioModeMeta))
	for _, m := range []StudioMode{
		ModeErase, ModeUpscale, ModeRepair, ModeBackground, ModeOutpaint, ModeWatermark,
	} {
		meta := studioModeMeta[m]
		out = append(out, StudioModeInfo{
			ID:            string(m),
			Label:         meta.label,
			RequiresMask:  meta.requiresMask,
			MaskHint:      meta.maskHint,
			DefaultPrompt: meta.defaultPrompt,
		})
	}
	return out
}

type studioMeta struct {
	label, maskHint, defaultPrompt string
	requiresMask                   bool
}

var studioModeMeta = map[StudioMode]studioMeta{
	ModeErase: {
		label: "消除", requiresMask: true,
		maskHint: "必须用红色画笔涂抹要删除的区域",
		defaultPrompt: "移除蒙版区域内的物体，用与周围环境一致的自然背景无缝填充，保持真实照片质感。",
	},
	ModeUpscale: {
		label: "高清增强", requiresMask: false,
		maskHint: "整图增强；可选涂蒙版仅增强局部",
		defaultPrompt: "提升整张图片的清晰度与细节，减少模糊与噪点，保持色彩自然、无过度锐化。",
	},
	ModeRepair: {
		label: "修复", requiresMask: false,
		maskHint: "可整图修复；涂蒙版则只修复局部（划痕、破损、噪点）",
		defaultPrompt: "修复图片中的瑕疵、划痕、破损与噪点，使画面干净自然，保持原有构图与风格。",
	},
	ModeBackground: {
		label: "换背景", requiresMask: false,
		maskHint: "可整图换背景；涂蒙版可保留主体仅替换背景区域",
		defaultPrompt: "替换为简洁自然的背景，主体边缘清晰，光照与透视一致，整体像真实拍摄。",
	},
	ModeOutpaint: {
		label: "扩图", requiresMask: false,
		maskHint: "整图扩展画布；涂蒙版可指定扩展方向或区域",
		defaultPrompt: "在保持原图主体不变的前提下自然扩展画面边缘，延伸场景与环境，透视与风格一致。",
	},
	ModeWatermark: {
		label: "去水印", requiresMask: false,
		maskHint: "涂蒙版标出水印位置效果更好；也可整图自动查找去除",
		defaultPrompt: "去除图片中的水印、文字叠加与 logo，用周围纹理自然填充，不留明显修补痕迹。",
	},
}

// BuildStudioPrompt 合并模式默认说明与用户输入
func BuildStudioPrompt(mode StudioMode, userPrompt string) string {
	meta := studioModeMeta[mode]
	user := strings.TrimSpace(userPrompt)
	if user == "" {
		return meta.defaultPrompt
	}
	return meta.defaultPrompt + "\n\nAdditional requirements from user: " + user
}

// StudioResultMessage 成功提示文案
func StudioResultMessage(mode StudioMode, quality string, w, h int) string {
	label := studioModeMeta[mode].label
	return "GPT Image " + label + "完成 (quality=" + quality + ", " + strconv.Itoa(w) + "x" + strconv.Itoa(h) + ")"
}
