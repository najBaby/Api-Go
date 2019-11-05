package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	W         http.ResponseWriter `json:"-"`
	Operation string              `json:"operation"`
	Data      interface{}         `json:"data"`
	Errors    []string            `json:"errors,omitempty"`
}

func NewResponse(w http.ResponseWriter) *Response {
	res := Response{
		Data:      nil,
		Operation: "",
		W:         w,
	}
	return &res
}

func (res *Response) ServeJSON() {
	b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	res.W.Write(b)
}
