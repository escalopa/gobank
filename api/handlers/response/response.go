package response

type JSON struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty" swaggerignore:"true"`
	Error   *Error      `json:"error,omitempty" swaggerignore:"true"`
}

type Error struct {
	Error string `json:"error,omitempty"`
}

func Success(data interface{}) JSON {
	return JSON{Success: true, Data: data}
}

func Err(err error) JSON {
	return JSON{Success: false, Error: &Error{Error: err.Error()}}
}
