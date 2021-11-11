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
	"strconv"
)

//	start listening the TCP/IP port and serve HTTP requests
func ListenAndServe(portNumber int32) error {

	fmt.Printf("> Listening on port %d\n", portNumber)

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

	requestLine, err := bufio.NewReader(clientConnection).ReadBytes('\n')

	if err != nil {
		clientConnection.Close()
		return err
	}

	log.Println("Client message:", string(requestLine[:len(requestLine)-1]))

	clientConnection.Write([]byte("200 OK"))
	clientConnection.Close()

	return nil
}
