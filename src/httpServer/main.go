////////////////////////////////////////////////////////////////////////////////
//	main.go  -  Nov-9-2021  -  aldebap
//
//	HTTP Server Protocol
////////////////////////////////////////////////////////////////////////////////

package httpServer

import (
	"fmt"
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

		fmt.Print("> New client from " + clientConnection.RemoteAddr().String() + " connected\n")
	}

	return nil
}
