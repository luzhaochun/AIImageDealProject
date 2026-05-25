package response

// Body 统一 API 响应（与 docs/05-API.md 对齐）
type Body struct {
	Data      interface{} `json:"data,omitempty"`
	Code      string      `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

func OK(data interface{}) Body {
	return Body{Data: data}
}

func Err(code, message string) Body {
	return Body{Code: code, Message: message}
}
