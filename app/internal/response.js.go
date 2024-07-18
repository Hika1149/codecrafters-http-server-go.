package internal

import "fmt"

/*
Response
	Status line
		HttpVersion StatusCode ReasonPhrase
	Header
		HeaderName: HeaderValue
	Body
*/

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
func (r *Response) SetContentType(str string) *Response {
	r.WriteHeader("Content-Type", str)
	return r
}
func (r *Response) SetContentLength(str string) *Response {
	r.WriteHeader("Content-Length", str)
	return r
}
func (r *Response) SetContentEncoding(str string) *Response {
	r.WriteHeader("Content-Encoding", str)
	return r
}
