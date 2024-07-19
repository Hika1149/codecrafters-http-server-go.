package internal

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

/*
Request
	Request line
		HttpMethod RequestTarget HttpVersion
	Headers
	Body
*/

type Request struct {
	Method    string
	Path      string
	UserAgent string

	Headers map[string]string

	Body []byte
}

func NewRequest() *Request {
	return &Request{

		Headers: make(map[string]string),
	}
}

/**

bufio.NewScanner(conn) cause reading loop
curl sends its request but does not close its side of the TCP connection afterward as is typical for most real-world clients
and especially web browsers (they are not obliged to do that), and waits for response.
Since your scanning code does not ever see the curl's (writing) side of the connection closed, it never receives an io.EOF "error" and hence never stops its reading loop

*/

func ConnToRequest(conn net.Conn) *Request {

	req := NewRequest()
	/** read by line */
	connByte := make([]byte, 1024*4)
	_, err := conn.Read(connByte)
	if err != nil {
		fmt.Printf("conn read failed %v\n", err)
		return req
	}
	rb := strings.NewReader(string(connByte))

	scanner := bufio.NewScanner(rb)
	scanner.Split(bufio.ScanLines)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		fmt.Println("scan text: ", scanner.Text())
		if i == 0 {
			req.Method = strings.Split(line, " ")[0]
			req.Path = strings.Split(line, " ")[1]
			continue
		}
		// reached the end of the headers
		if line == "" {
			break
		}

		if strings.HasPrefix(line, "User-Agent: ") {
			req.UserAgent = strings.TrimPrefix(line, "User-Agent: ")
		}

		req.Headers[strings.Split(line, ": ")[0]] = strings.Split(line, ": ")[1]

	}
	// set req body
	body := ""
	for scanner.Scan() {
		body += scanner.Text()
	}

	contentLength, _ := strconv.Atoi(req.Headers["Content-Length"])
	//fmt.Println("scan content length ", contentLength)
	req.Body = []byte(body)[:contentLength]
	//fmt.Println("scan body: ", body, len(req.Body))
	//fmt.Println("scan completed")

	return req
}
