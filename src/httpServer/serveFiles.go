////////////////////////////////////////////////////////////////////////////////
//	serveFiles.go  -  Nov-13-2021  -  aldebap
//
//	HTTP Server Protocol
////////////////////////////////////////////////////////////////////////////////

package httpServer

type DirectoryConfiguration struct {
	Context              string
	DirectoryName        string
	DefaultFile          string
	Navigation           bool
	FollowSubdirectories bool
}

//	add a directory to be served
func (server *Server) ServeDirectory(directoryConfig DirectoryConfiguration) error {

	server.Directory = append(server.Directory, directoryConfig)

	return nil
}
