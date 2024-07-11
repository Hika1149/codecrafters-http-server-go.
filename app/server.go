package main

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/internal"
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

	fmt.Println("req: ", string(req))
	reqSections := strings.Split(string(req), "\r\n")

	/** find user agent header */
	var userAgent string
	for _, section := range reqSections {
		if strings.HasPrefix(section, "User-Agent: ") {
			userAgent = strings.TrimPrefix(section, "User-Agent: ")
		}
	}

	// get request line
	reqLine := reqSections[0]
	path := strings.Split(reqLine, " ")[1]

	fmt.Println("path: ", path)

	res := internal.NewResponse()

	if strings.HasPrefix(path, "/echo/") {
		echo := strings.Split(path, "/echo/")[1]
		res.WriteStatusOk().
			WriteHeader("Content-Type", "text/plain").
			WriteHeader("Content-Length", fmt.Sprintf("%v", len(echo))).
			WriteHeadersEnd().
			WriteBody(echo)
	} else if strings.HasPrefix(path, "/user-agent") {
		res.WriteStatusOk().
			WriteHeader("Content-Type", "text/plain").
			WriteHeader("Content-Length", fmt.Sprintf("%v", len(userAgent))).
			WriteHeadersEnd().
			WriteBody(userAgent)
	} else if path == "/" {
		res.WriteStatusOk().
			WriteHeadersEnd()
	} else {
		res.WriteStatusLine("404", "Not Found").
			WriteHeadersEnd()
	}

	_, err = conn.Write(res.Buffer)
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
