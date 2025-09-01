package domain

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseSingleData[Data any] struct {
	Code    int    `json:"code"`    // number
	Data    Data   `json:"data"`    // of data
	Message string `json:"message"` // string
}

type ResponseMultipleData[Data any] struct {
	Code    int    `json:"code"`    // number
	Data    []Data `json:"data"`    // list of data
	Message string `json:"message"` // string
}

type Empty struct{}
