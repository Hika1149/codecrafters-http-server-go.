package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func handleConnection(conn net.Conn) {

	req := make([]byte, 1024)
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

	res := "HTTP/1.1 200 OK\r\n\r\n"

	if path != "/" {
		res = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	_, err = conn.Write([]byte(res))
	if err != nil {
		fmt.Printf("conn write failed %v\n", err)
		return
	}
}

/**
request line
HttpMethod  RequestTarget HttpVersion

*/

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
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
