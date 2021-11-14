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
	"regexp"
	"strconv"
	"strings"
)

type Server struct {
	Address   string
	Directory []DirectoryConfiguration
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

	log.Println("> Listening socket: " + server.Address)

	socketListening, err := net.Listen("tcp", server.Address)
	if err != nil {
		return err
	}
	defer socketListening.Close()

	//	wait for clients to connect to the server
	for {

		clientConnection, err := socketListening.Accept()
		if err != nil {
			return err
		}

		go server.handleHttpClient(clientConnection)
	}

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

	//	check if the resource can be served
	var responseContent string

	for i := 0; i < len(server.Directory); i++ {

		if strings.Index(request.Resource, server.Directory[i].Context) == 0 {

			responseContent = "<html><head /><body>Resouce: " + request.Resource + " being served from directory: " + server.Directory[i].DirectoryName + "</html>"
		}
	}

	//	if there's no content, the resource was not found
	if len(responseContent) == 0 {

		responseWriter.Write([]byte("HTTP/1.0 404 Resource not found\n"))
		responseWriter.Write([]byte("\n"))

		return errors.New("resource not found")
	}

	//	dummy response for now
	_, err := responseWriter.Write([]byte("HTTP/1.0 200 OK\n"))
	if err != nil {
		fmt.Printf("[debug] Error writing response content: %s\n", err)
	}
	responseWriter.Write([]byte("Content-Type: text/html\n"))
	responseWriter.Write([]byte("Content-Length: " + strconv.Itoa(len(responseContent)) + "\n"))
	responseWriter.Write([]byte("\n"))
	responseWriter.Write([]byte(responseContent))

	fmt.Printf("[debug] Response content: %s\n", responseContent)

	return nil
}
