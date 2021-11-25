////////////////////////////////////////////////////////////////////////////////
//	main.go  -  Nov-9-2021  -  aldebap
//
//	HTTP Server Protocol
////////////////////////////////////////////////////////////////////////////////

package httpServer

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RequestHandler struct {
	Context         string
	DirectoryConfig DirectoryConfiguration
}

type Server struct {
	Address string
	Handler []RequestHandler
	Stop    bool
}

type Request struct {
	Method   string
	Resource string
	Protocol string

	Host           string
	UserAgent      string
	Accept         string
	AcceptEncoding string
}

type Response struct {
}

//	HTTP methods
const Get = "GET"
const Post = "POST"
const Put = "PUT"
const Patch = "PATCH"
const Delete = "DELETE"

//	HTTP response status
const StatusOk = 200
const StatusBadRequest = 400
const StatusNotFound = 404
const StatusMethodNotAllowed = 405

//	standard request headers
const HostHeader = "Host"
const UserAgentHeader = "User-Agent"
const AcceptHeader = "Accept"
const AcceptEncodingHeader = "Accept-Encoding"

//	regexs required to parse HTTP requests
var startLineRegEx *regexp.Regexp
var requestHeaderRegEx *regexp.Regexp

//	initialize a HTTP Server
func (server *Server) Initialize(portNumber int32) error {

	server.Address = "localhost:" + strconv.Itoa(int(portNumber))

	//	compile all regexs required to parse HTTP requests
	startLineRegEx, _ = regexp.Compile(`^(\S+)\s+(\S+)\s+(\S.*)$`)
	requestHeaderRegEx, _ = regexp.Compile(`^(\S+):\s+(\S.*)$`)

	return nil
}

//	start listening the TCP/IP port and serve HTTP requests
func (server *Server) ListenAndServe() error {

	log.Println("> Listening TCP socket: " + server.Address)

	server.Stop = false

	serverAddress, err := net.ResolveTCPAddr("tcp", server.Address)
	if err != nil {
		return err
	}

	socketListening, err := net.ListenTCP("tcp", serverAddress)
	if err != nil {
		return err
	}
	defer socketListening.Close()

	//	wait for clients to connect to the server
	for {
		if server.Stop {
			break
		}

		//	set a timeout to make it possible to check if a stop method was invoked
		err = socketListening.SetDeadline(time.Now().Add(100 * time.Millisecond))
		if err != nil {
			return err
		}

		clientConnection, err := socketListening.Accept()
		if err == nil {

			go server.handleHttpClient(clientConnection)
			continue
		}

		//	check if a timeout was reached, or an error really occurred
		if !errors.Is(err, os.ErrDeadlineExceeded) {

			return err
		}
	}

	return nil
}

//	stop listening the TCP/IP port and serving HTTP requests
func (server *Server) StopServer() error {

	log.Println("> Stop server: " + server.Address)

	return nil
}

//	handle HTTP request
func (server *Server) handleHttpClient(clientConnection net.Conn) error {

	log.Print("> New client from " + clientConnection.RemoteAddr().String() + " connected")

	var request Request

	//	get a reader for the request
	requestReader := bufio.NewReader(clientConnection)

	//	read the start line
	requestLine, _, err := requestReader.ReadLine()
	if err != nil {
		clientConnection.Close()
		return err
	}

	tokens := startLineRegEx.FindAllStringSubmatch(string(requestLine), -1)

	if nil != tokens && len(tokens[0]) == 4 {

		request.Method = tokens[0][1]
		request.Resource = tokens[0][2]
		request.Protocol = tokens[0][3]
	} else {

		log.Printf("Bad request: %s\n", string(requestLine))

		clientConnection.Write([]byte("400 Bad Request"))
		clientConnection.Close()

		return nil
	}

	//	read the request headers
	for {

		requestLine, _, err := requestReader.ReadLine()
		if err != nil {
			clientConnection.Close()
			return err
		}

		//	the first empty line indicates the end of headers
		if len(requestLine) == 0 {
			break
		}

		tokens = requestHeaderRegEx.FindAllStringSubmatch(string(requestLine), -1)

		if nil != tokens && len(tokens[0]) == 3 {

			//	check if it's a standard header
			if HostHeader == tokens[0][1] {

				request.Host = tokens[0][2]
			} else if UserAgentHeader == tokens[0][1] {

				request.UserAgent = tokens[0][2]
			} else if AcceptHeader == tokens[0][1] {

				request.Accept = tokens[0][2]
			} else if AcceptEncodingHeader == tokens[0][1] {

				request.AcceptEncoding = tokens[0][2]
			}
		} else {

			log.Printf("Bad request: %s\n", string(requestLine))

			clientConnection.Write([]byte("400 Bad Request"))
			clientConnection.Close()

			return nil
		}
	}

	responseWriter := bufio.NewWriter(clientConnection)
	server.handleHttpRequest(request, responseWriter)
	responseWriter.Flush()
	clientConnection.Close()

	return nil
}

func (server *Server) handleHttpRequest(request Request, responseWriter *bufio.Writer) error {

	fmt.Printf("[debug] Request: %s\n", request)

	var responseCode int
	var responseContent []byte
	var err error

	switch request.Method {

	case Get:
		//	check if the resource can be served
		for i := 0; i < len(server.Handler); i++ {

			if strings.Index(request.Resource, server.Handler[i].Context) == 0 {

				if len(server.Handler[i].Context)+1 > len(request.Resource) {

					responseContent, err = handleFilesFromDirectory(server.Handler[i].Context, "", server.Handler[i].DirectoryConfig)
				} else {

					responseContent, err = handleFilesFromDirectory(server.Handler[i].Context, request.Resource[len(server.Handler[i].Context)+1:], server.Handler[i].DirectoryConfig)
				}
				if err == nil {
					break
				}
			}
		}
		//	if there's no content, the resource was not found
		if len(responseContent) > 0 {

			responseCode = StatusOk
		} else {

			responseCode = StatusNotFound
		}

	case Post:
		responseCode = StatusMethodNotAllowed
	case Put:
		responseCode = StatusMethodNotAllowed
	case Patch:
		responseCode = StatusMethodNotAllowed
	case Delete:
		responseCode = StatusMethodNotAllowed
	}

	//	output the response
	switch responseCode {

	case StatusOk:
		_, err = responseWriter.Write([]byte("HTTP/1.0 200 OK\n"))
		if err != nil {
			fmt.Printf("[debug] Error writing response content: %s\n", err)
		}
		responseWriter.Write([]byte("Content-Type: text/html\n"))
		responseWriter.Write([]byte("Content-Length: " + strconv.Itoa(len(responseContent)) + "\n"))
		responseWriter.Write([]byte("\n"))
		responseWriter.Write(responseContent)

	case StatusNotFound:
		responseWriter.Write([]byte("HTTP/1.0 404 Resource not found\n"))
		responseWriter.Write([]byte("\n"))
		return errors.New("resource not found")

	case StatusMethodNotAllowed:
		responseWriter.Write([]byte("HTTP/1.0 405 Method not allowed\n"))
		responseWriter.Write([]byte("\n"))
		return errors.New("Method not allowed")
	}

	return nil
}
