package utils

type StandardResponse struct {
	Status   string      `json:"status,omitempty"`
	Message  string      `json:"message,omitempty"`
	Errors   interface{} `json:"errors,omitempty"`
	Result   interface{} `json:"result,omitempty"`
	MetaData interface{} `json:"metadata,omitempty"`
}
