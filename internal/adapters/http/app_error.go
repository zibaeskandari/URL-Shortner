package http

type AppError struct {
	Code      string      `json:"code,omitempty"`
	Message   string      `json:"message"`           // human friendly
	Details   interface{} `json:"details,omitempty"` // field errors, extra info
	Status    int         `json:"-"`                 // HTTP status, not serialized
	RequestID string      `json:"request_id,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}
