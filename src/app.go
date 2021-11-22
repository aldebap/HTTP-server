////////////////////////////////////////////////////////////////////////////////
//	app.go  -  Nov-9-2021  -  aldebap
//
//	HTTP Web Server Application
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"HTTP-server/httpServer"
)

//	configuration file structures
type directoryConfiguration struct {
	Context              string `json:"context"`
	DirectoryName        string `json:"directory"`
	DefaultFile          string `json:"defaultFile"`
	Navigation           bool   `json:"directoryNavigation"`
	FollowSubdirectories bool   `json:"followSubdirectories"`
}

type configurationData struct {
	PortNumber      int32                    `json:"portNumber"`
	ServerDirectory []directoryConfiguration `json:"serverDirectory"`
}

type App struct {
	Configuration configurationData
	HttpServer    httpServer.Server
}

//	variables to functions that need to be mocked
var readAll = ioutil.ReadAll

//	initialize the application
func (a *App) Initialize(configurationFileName string) error {

	configurationDataFile, err := os.Open(configurationFileName)
	if err != nil {
		return err
	}
	defer configurationDataFile.Close()

	//	load the configuration file
	fileContent, err := readAll(configurationDataFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileContent, &a.Configuration)
	if err != nil {
		return err
	}

	a.HttpServer.Initialize(a.Configuration.PortNumber)

	return nil
}

//	run the application
func (a *App) Run() error {

	//	add the directories list to the server configuration
	for _, directory := range a.Configuration.ServerDirectory[0:] {

		var serverDirectoryConfig httpServer.DirectoryConfiguration

		serverDirectoryConfig.DirectoryName = directory.DirectoryName
		serverDirectoryConfig.DefaultFile = directory.DefaultFile
		serverDirectoryConfig.Navigation = directory.Navigation
		serverDirectoryConfig.FollowSubdirectories = directory.FollowSubdirectories

		a.HttpServer.ServeDirectory(directory.Context, serverDirectoryConfig)
	}

	return a.HttpServer.ListenAndServe()
}
