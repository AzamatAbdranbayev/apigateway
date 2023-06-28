package pkg

import (
	"encoding/json"
	"time"
)

type Response struct {
	Status      bool        `json:"status"`
	Errors      interface{} `json:"errors"`
	Data        interface{} `json:"data"`
	TmRequest   string      `json:"tm_req"`
	TmRequestSt time.Time   `json:"-"`
}
type ErrorItem struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func InitResp() *Response {
	return &Response{
		TmRequestSt: time.Now(),
	}
}
func (r *Response) SetError(code int, mess string) *Response {
	r.Errors = ErrorItem{
		Code:    code,
		Message: mess,
	}

	return r
}

func (r *Response) SetValue(val interface{}) *Response {
	r.Data = val
	return r
}
func (r *Response) FormResponse() *Response {

	if r.Errors == nil {
		r.Status = true
	} else {
		r.Status = false
	}
	return r
}

func (r *Response) Json() []byte {
	r.TmRequest = time.Now().Sub(r.TmRequestSt).String()
	if bts, err := json.Marshal(r); err == nil {
		return bts
	}
	return []byte{}
}
