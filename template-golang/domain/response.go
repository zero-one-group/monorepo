package domain

type Response struct {
    Code int    `json:"code"`
    Status string `json:"status"`
	Message string    `json:"message"`
}

type ResponseMultipleData[D any] struct {
	RequestID string `json:"request_id"` // string
	Code      int    `json:"code"`       // number
	Status    string `json:"status"`     // string
	Data      []D    `json:"data"`       // list of data
	Message   string `json:"message"`    // string
}

type Empty struct{}

