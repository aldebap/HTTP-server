////////////////////////////////////////////////////////////////////////////////
//	main.go  -  Nov-8-2021  -  aldebap
//
//	HTTP Web Server entry point
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////////
//	Start the HTTP web server Application
////////////////////////////////////////////////////////////////////////////////

func main() {

	//	splash screen
	fmt.Printf(">>> Starting web server\n\n")

	//	start the HTTP web server
	webServerApp := App{}

	err := webServerApp.Initialize("config/httpServer-config.json")
	if nil != err {

		fmt.Printf("[error] %s\n", err)
		os.Exit(-1)
	} else {

		err = webServerApp.Run()
		if nil != err {

			fmt.Printf("[error] %s\n", err)
		}
	}
}
