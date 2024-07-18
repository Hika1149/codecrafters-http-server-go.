package main

import (
	"flag"
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/internal"
	"net"
	"os"
	"strings"
)

var directory = flag.String("directory", "", "")

var SupportedCompressions = []string{"gzip"}

func GetFileContent(res *internal.Response, path string) {
	filename := strings.TrimPrefix(path, "/files/")
	file, err := os.Open(fmt.Sprintf("%v%v", *directory, filename))
	if err != nil {
		fmt.Printf("open file failed err=%v\n", err)
		res.WriteStatusLine("404", "Not Found").WriteHeadersEnd()
	} else {
		stat, err := file.Stat()
		if err != nil {
			fmt.Printf("file stat failed  err=%v\n", err)
			return
		}
		content := make([]byte, 1024*4)
		_, err = file.Read(content)
		if err != nil {
			fmt.Printf("file Read failed  err=%v\n", err)
			return
		}
		res.WriteStatusOk().
			SetContentType("application/octet-stream").
			SetContentLength(fmt.Sprintf("%v", stat.Size())).
			WriteHeadersEnd().
			WriteBody(string(content))
	}
}

func PostWriteFile(req *internal.Request, res *internal.Response, path string) {
	filename := strings.TrimPrefix(path, "/files/")
	file, err := os.Create(fmt.Sprintf("%v%v", *directory, filename))
	defer file.Close()
	if err != nil {
		fmt.Printf("create file failed err=%v\n", err)
		res.WriteStatusLine("500", "Internal Server Error").WriteHeadersEnd()
		return
	}
	_, err = file.Write(req.Body)
	if err != nil {
		fmt.Printf("write file failed err=%v\n", err)
		res.WriteStatusLine("500", "Internal Server Error").WriteHeadersEnd()
		return
	}

	res.WriteStatusLine("201", "Created").WriteHeadersEnd()

}
func handleConnection(conn net.Conn) {

	request := internal.ConnToRequest(conn)

	// get request line
	path := request.Path
	userAgent := request.UserAgent

	// get compress scheme that client supports
	contentEncoding := request.Headers["Accept-Encoding"]
	var compressMethod string
	for _, c := range SupportedCompressions {
		if contentEncoding == c {
			compressMethod = c
			break
		}
	}

	fmt.Println("path: ", path)
	res := internal.NewResponse()

	if strings.HasPrefix(path, "/echo/") {
		echo := strings.Split(path, "/echo/")[1]

		res.WriteStatusOk().
			SetContentType("text/plain").
			SetContentLength(fmt.Sprintf("%v", len(echo)))

		if compressMethod != "" {
			res.SetContentEncoding(compressMethod)
		}

		res.WriteHeadersEnd().
			WriteBody(echo)

	} else if strings.HasPrefix(path, "/user-agent") {
		res.WriteStatusOk().
			SetContentType("text/plain").
			SetContentLength(fmt.Sprintf("%v", len(userAgent))).
			WriteHeadersEnd().
			WriteBody(userAgent)
	} else if strings.HasPrefix(path, "/files") {
		if request.Method == "POST" {
			PostWriteFile(request, res, path)
		} else {
			GetFileContent(res, path)
		}

	} else if path == "/" {
		res.WriteStatusOk().
			WriteHeadersEnd()
	} else {
		res.WriteStatusLine("404", "Not Found").
			WriteHeadersEnd()
	}

	_, err := conn.Write(res.Buffer)
	if err != nil {
		fmt.Printf("conn write failed %v\n", err)
		return
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	flag.Parse()

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
