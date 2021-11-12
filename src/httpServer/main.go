////////////////////////////////////////////////////////////////////////////////
//	main.go  -  Nov-9-2021  -  aldebap
//
//	HTTP Server Protocol
////////////////////////////////////////////////////////////////////////////////

package httpServer

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
)

//	start listening the TCP/IP port and serve HTTP requests
func ListenAndServe(portNumber int32) error {

	log.Println("> Listening on port " + strconv.Itoa(int(portNumber)))

	socketListening, err := net.Listen("tcp", "localhost:"+strconv.Itoa(int(portNumber)))
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

		go handleHttpRequest(clientConnection)
	}

	return nil
}

//	handle HTTP request
func handleHttpRequest(clientConnection net.Conn) error {

	log.Print("> New client from " + clientConnection.RemoteAddr().String() + " connected")

	//	get a reader for the request
	requestReader := bufio.NewReader(clientConnection)

	//	read the start line
	requestLine, _, err := requestReader.ReadLine()
	if err != nil {
		clientConnection.Close()
		return err
	}

	//	TODO: this regexp compilation should be made only once for better performance
	startLineRegEx, _ := regexp.Compile(`^(\S+)\s+(\S+)\s+(\S.*)$`)
	tokens := startLineRegEx.FindAllStringSubmatch(string(requestLine), -1)

	if nil != tokens && len(tokens[0]) == 4 {

		fmt.Printf("Operation: %s - Resource: %s - Protocol: %s\n", tokens[0][1], tokens[0][2], tokens[0][3])
	} else {

		fmt.Printf("Bad request: %s\n", string(requestLine))

		clientConnection.Write([]byte("400 Bad Request"))
		clientConnection.Close()

		return nil
	}

	//	read the request headers
	requestHeaderRegEx, _ := regexp.Compile(`^(\S+):\s+(\S.*)$`)

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

			fmt.Printf("Header: %s --> %s\n", tokens[0][1], tokens[0][2])
		} else {

			fmt.Printf("Bad request: %s\n", string(requestLine))

			clientConnection.Write([]byte("400 Bad Request"))
			clientConnection.Close()

			return nil
		}
	}

	clientConnection.Write([]byte("HTTP/1.0 200 OK\n"))
	clientConnection.Write([]byte("Content-Type: text/html\n"))
	clientConnection.Write([]byte("Content-Length: 14\n"))
	clientConnection.Write([]byte("\n"))
	clientConnection.Write([]byte("<html></html>\n"))
	clientConnection.Close()

	return nil
}
