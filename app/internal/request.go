package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
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
	return &Request{}
}

func scanDoubleNewLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	//if i := bytes.Index(data, []byte("\r\n")); i >= 0 {
	//	return i + 2, data[0:i], nil
	//}

	if i := bytes.Index(data, []byte("\r\n\r\n")); i >= 0 {
		return i + 4, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return

}

func ConnToRequest(conn net.Conn) *Request {
	/** read by line */
	scanner := bufio.NewScanner(conn)

	scanner.Split(bufio.ScanLines)
	//scanner.Split(scanDoubleNewLine)

	req := NewRequest()

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

		// reached the end of the headers
		//if strings.HasSuffix(line, "\r\n\r\n") {
		//	break
		//}

	}

	fmt.Println("scan completed")
	body := make([]byte, 1024*4)

	_, err := conn.Read(body)
	if err != nil {
		fmt.Printf("conn read failed %v\n", err)
	}

	fmt.Println("read text body: ", body)
	return req
}
