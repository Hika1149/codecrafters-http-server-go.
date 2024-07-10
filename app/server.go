package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

/*
*
Request

	Request line
		HttpMethod  RequestTarget HttpVersion

Response

	Status line
		HttpVersion StatusCode ReasonPhrase
	Header
		HeaderName: HeaderValue
	Body
*/
func handleConnection(conn net.Conn) {

	req := make([]byte, 1024*5)
	_, err := conn.Read(req)
	if err != nil {
		fmt.Printf("conn read failed %v\n", err)
		return
	}
	reqSections := strings.Split(string(req), "\r\n")

	// get request line
	reqLine := reqSections[0]
	path := strings.Split(reqLine, " ")[1]

	fmt.Println("path: ", path)

	res := make([]byte, 0)

	if strings.HasPrefix(path, "/echo/") {
		echo := strings.Split(path, "/echo/")[1]
		res = append(res, []byte("HTTP/1.1 200 OK\r\n")...)
		res = append(res, []byte("Content-Type: text/plain\r\n")...)
		res = append(res, []byte(fmt.Sprintf("Content-Length: %v\r\n", len(echo)))...)
		res = append(res, []byte("\r\n")...) // End of headers
		res = append(res, []byte(fmt.Sprintf("%s\r\n", echo))...)
	} else if path == "/" {
		res = append(res, []byte("HTTP/1.1 200 OK\r\n\r\n")...)
	} else {
		res = append(res, []byte("HTTP/1.1 404 Not Found\r\n\r\n")...)
	}

	_, err = conn.Write(res)
	if err != nil {
		fmt.Printf("conn write failed %v\n", err)
		return
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}

}
