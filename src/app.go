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

//	initialize the application
func (a *App) Initialize(configurationFileName string) error {

	configurationDataFile, err := os.Open(configurationFileName)
	if err != nil {
		return err
	}

	//	load the configuration file
	fileContent, err := ioutil.ReadAll(configurationDataFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileContent, &a.Configuration)
	if err != nil {
		return err
	}

	defer configurationDataFile.Close()

	a.HttpServer.Initialize(a.Configuration.PortNumber)

	return nil
}

//	run the application
func (a *App) Run() error {

	//	add the directories list to the server configuration
	for i := 0; i < len(a.Configuration.ServerDirectory); i++ {

		var directoryConfig httpServer.DirectoryConfiguration

		directoryConfig.DirectoryName = a.Configuration.ServerDirectory[i].DirectoryName
		directoryConfig.DefaultFile = a.Configuration.ServerDirectory[i].DefaultFile
		directoryConfig.Navigation = a.Configuration.ServerDirectory[i].Navigation
		directoryConfig.FollowSubdirectories = a.Configuration.ServerDirectory[i].FollowSubdirectories

		a.HttpServer.ServeDirectory(a.Configuration.ServerDirectory[i].Context, directoryConfig)
	}

	return a.HttpServer.ListenAndServe()
}
