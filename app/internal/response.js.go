package internal

import "fmt"

type Response struct {
	Buffer []byte
}

func NewResponse() *Response {
	return &Response{Buffer: make([]byte, 0)}
}

func (r *Response) WriteStatusOk() *Response {
	r.WriteStatusLine("200", "OK")
	return r
}
func (r *Response) WriteStatusLine(code, statusText string) *Response {
	r.Buffer = append(r.Buffer, []byte(fmt.Sprintf("HTTP/1.1 %s %s\r\n", code, statusText))...)
	return r
}

func (r *Response) WriteHeader(key, val string) *Response {
	r.Buffer = append(r.Buffer, []byte(fmt.Sprintf("%s: %s\r\n", key, val))...)
	return r
}

func (r *Response) WriteHeadersEnd() *Response {
	r.Buffer = append(r.Buffer, []byte("\r\n")...)
	return r
}

func (r *Response) WriteBody(body string) *Response {
	r.Buffer = append(r.Buffer, []byte(body)...)
	return r
}
