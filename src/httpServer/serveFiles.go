////////////////////////////////////////////////////////////////////////////////
//	serveFiles.go  -  Nov-13-2021  -  aldebap
//
//	HTTP Server Protocol
////////////////////////////////////////////////////////////////////////////////

package httpServer

import (
	"fmt"
	"io/fs"
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
func handleFilesFromDirectory(context string, resource string, directoryConfig DirectoryConfiguration) ([]byte, error) {

	var err error

	if len(resource) == 0 {

		if len(directoryConfig.DefaultFile) > 0 {
			fmt.Printf("[debug] attempting to find resource: %s\n", directoryConfig.DefaultFile)

			resourceFile, err := os.Open(directoryConfig.DirectoryName + "/" + directoryConfig.DefaultFile)
			if err == nil {
				resourceFileContent, _ := ioutil.ReadAll(resourceFile)

				return resourceFileContent, nil
			}
		} else if directoryConfig.Navigation {
			fmt.Printf("[debug] attempting to navigate the files in the directory: %s\n", directoryConfig.DirectoryName)

			//	generate an HTML response with a list of files and directories in the current directory
			directoryContent, err := fs.ReadDir(os.DirFS(directoryConfig.DirectoryName), ".")
			if err == nil {

				var directoryList string

				for _, directoryEntry := range directoryContent[0:] {
					directoryList += `<li>`
					if directoryEntry.IsDir() {

						directoryList += `<svg aria-label="Directory" aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-file-directory hx_color-icon-directory">
	    									<path fill-rule="evenodd" d="M1.75 1A1.75 1.75 0 000 2.75v10.5C0 14.216.784 15 1.75 15h12.5A1.75 1.75 0 0016 13.25v-8.5A1.75 1.75 0 0014.25 3h-6.5a.25.25 0 01-.2-.1l-.9-1.2c-.33-.44-.85-.7-1.4-.7h-3.5z"></path>
										</svg>`
					} else {

						directoryList += `<svg aria-label="File" aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-file color-fg-muted">
											<path fill-rule="evenodd" d="M3.75 1.5a.25.25 0 00-.25.25v11.5c0 .138.112.25.25.25h8.5a.25.25 0 00.25-.25V6H9.75A1.75 1.75 0 018 4.25V1.5H3.75zm5.75.56v2.19c0 .138.112.25.25.25h2.19L9.5 2.06zM2 1.75C2 .784 2.784 0 3.75 0h5.086c.464 0 .909.184 1.237.513l3.414 3.414c.329.328.513.773.513 1.237v8.086A1.75 1.75 0 0112.25 15h-8.5A1.75 1.75 0 012 13.25V1.75z"></path>
										</svg>`
					}
					directoryList += `<a href="` + context + "/" + directoryEntry.Name() + `">` + directoryEntry.Name() + `</a></li>`
				}

				resultHTML := `<!DOCTYPE html>
				<html lang="en">
					<head>
						<title>` + directoryConfig.DirectoryName + `</title>
					</head>
					<body>
						<ul>` + directoryList + `</ul>
					</body>`

				return []byte(resultHTML), nil
			}
		}
	} else {

		fmt.Printf("[debug] attempting to find resource: %s\n", resource)

		resourceFile, err := os.Open(directoryConfig.DirectoryName + "/" + resource)
		if err == nil {
			resourceFileContent, _ := ioutil.ReadAll(resourceFile)

			return resourceFileContent, nil
		}
	}

	return nil, err
}
