package util

import (
	"encoding/json"
	"goweb/author-admin/server/pkg/e"
	"log"
)

type Meta struct {
	Page      `json:"page"`
	RequestID string `json:"requestid"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`

	// Meta存储如请求ID、分页等信息
	Meta Meta `json:"meta"`
}

func (resp *Response) ToMap(depth int) map[string]interface{} {
	m := make(map[string]interface{})
	err := StructToMapWithJSONKey(*resp, m, depth)
	if err != nil {
		log.Println(err)
		return nil
	}

	return m
}

func (resp *Response) ToJSON(depth int) []byte {
	m := resp.ToMap(depth)

	r, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return nil
	}

	return r
}

func NewResponse() *Response {
	return &Response{
		Code:    e.SUCCESS,
		Message: e.GetMsg(e.SUCCESS),
		Meta: Meta{
			RequestID: GenUUID(),
		},
	}
}

func NeWReponseWithCode(code int) *Response {
	resp := NewResponse()
	resp.Code = code
	resp.Message = e.GetMsg(code)
	return resp
}

func FailedResponseMap(code int) map[string]interface{} {
	resp := NeWReponseWithCode(code)
	return resp.ToMap(0)
}
