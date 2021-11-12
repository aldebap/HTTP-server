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
	"strings"
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

	fmt.Print("> New client from " + clientConnection.RemoteAddr().String() + " connected\n")

	//	read the start line
	requestLine, err := bufio.NewReader(clientConnection).ReadBytes('\n')
	if err != nil {
		clientConnection.Close()
		return err
	}

	startLineRegEx, _ := regexp.Compile(`^(\S+)\s+(\S+)\s+(\S.*)$`)
	tokens := startLineRegEx.FindAllStringSubmatch(strings.TrimSuffix(string(requestLine), "\n"), -1)

	if nil != tokens && len(tokens[0]) == 4 {

		fmt.Printf("Operation: %s - Resource: %s - Protocol: %s\n", tokens[0][1], tokens[0][2], tokens[0][3])
	} else {

		fmt.Printf("Bad request: %s\n", string(requestLine))

		clientConnection.Write([]byte("400 Bad Request"))
		clientConnection.Close()

		return nil
	}

	//	read the request headers

	clientConnection.Write([]byte("200 OK"))
	clientConnection.Write([]byte("Content-Type: text/html"))
	clientConnection.Write([]byte("Content-Length: 13"))
	clientConnection.Write([]byte(""))
	clientConnection.Write([]byte("<html></html>"))
	clientConnection.Close()

	return nil
}
