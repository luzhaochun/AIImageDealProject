package ai

import (
	"encoding/json"
	"fmt"
)

type openAIErrorBody struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

func friendlyOpenAIHTTPError(status int, raw []byte) error {
	var body openAIErrorBody
	_ = json.Unmarshal(raw, &body)
	msg := body.Error.Message
	code := body.Error.Code

	switch code {
	case "billing_hard_limit_reached":
		return fmt.Errorf("OpenAI 账户已达计费硬上限：请到 https://platform.openai.com/settings/organization/billing 提高 Monthly budget 或充值后再试")
	case "insufficient_quota":
		return fmt.Errorf("OpenAI 配额不足：请检查账户余额与用量限制")
	case "rate_limit_exceeded":
		return fmt.Errorf("OpenAI 请求过于频繁，请稍后再试")
	}

	if msg != "" {
		return fmt.Errorf("OpenAI 错误 (%d): %s", status, msg)
	}
	return fmt.Errorf("OpenAI 错误 (%d): %s", status, string(raw))
}
