package shared

type StandardResponse struct {
	Status   string      `json:"status,omitempty"`
	Errors   []string    `json:"errors,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	MetaData interface{} `json:"metadata,omitempty"`
}
