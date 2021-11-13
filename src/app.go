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
type configurationData struct {
	PortNumber int32 `json:"portNumber"`
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
	//	TODO: need to add error handling here
	fileContent, _ := ioutil.ReadAll(configurationDataFile)

	json.Unmarshal(fileContent, &a.Configuration)

	defer configurationDataFile.Close()

	a.HttpServer.Initialize(a.Configuration.PortNumber)

	return nil
}

//	run the application
func (a *App) Run() error {

	return a.HttpServer.ListenAndServe()
}
