////////////////////////////////////////////////////////////////////////////////
//	serveFiles.go  -  Nov-13-2021  -  aldebap
//
//	HTTP Server Protocol
////////////////////////////////////////////////////////////////////////////////

package httpServer

import (
	"fmt"
	"io/ioutil"
	"os"
)

type DirectoryConfiguration struct {
	DirectoryName        string
	DefaultFile          string
	Navigation           bool
	FollowSubdirectories bool
}

//	add a directory to be served
func (server *Server) ServeDirectory(context string, directoryConfig DirectoryConfiguration) error {

	var directoryHandler RequestHandler

	directoryHandler.Context = context
	directoryHandler.DirectoryConfig = directoryConfig

	server.Handler = append(server.Handler, directoryHandler)

	return nil
}

//	return the file content if the resource matches a file in the directory
func handleFilesFromDirectory(resource string, directoryConfig DirectoryConfiguration) ([]byte, error) {

	if len(resource) == 0 {

		fmt.Printf("[debug] attempting to find resource: %s\n", directoryConfig.DefaultFile)

		resourceFile, err := os.Open(directoryConfig.DirectoryName + "/" + directoryConfig.DefaultFile)
		if err == nil {
			resourceFileContent, _ := ioutil.ReadAll(resourceFile)

			return resourceFileContent, nil
		}
	} else {

		fmt.Printf("[debug] attempting to find resource: %s\n", resource)

		resourceFile, err := os.Open(directoryConfig.DirectoryName + "/" + resource)
		if err == nil {
			resourceFileContent, _ := ioutil.ReadAll(resourceFile)

			return resourceFileContent, nil
		}
	}

	return nil, nil
}
