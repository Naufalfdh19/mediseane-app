package wrapper

type WrapperModel struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func ResponseData(data interface{}, message string, err interface{}) WrapperModel {
	result := WrapperModel{
		Data:    data,
		Message: message,
		Error:   err,
	}
	return result
}
