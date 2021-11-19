////////////////////////////////////////////////////////////////////////////////
//	app_test.go  -  Nov/18/2021  -  aldebap
//
//	App unit tests
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"io/ioutil"
	"os"
	"testing"
)

//	generate configuration without directories to perform some tests
func generateConfigWithoutDirectories() (*os.File, string, error) {

	//	create a temporary configuration file without directories
	tmpConfigFile, err := ioutil.TempFile(os.TempDir(), "config")
	if err != nil {
		return nil, "Unable to create temporary configuration file: %v", err
	}
	defer os.Remove(tmpConfigFile.Name())

	//	write the configuration to file
	configurationData := []byte(`{
			"portNumber": 1234
		}`)
	_, err = tmpConfigFile.Write(configurationData)
	if nil != err {
		return nil, "Unable to write to temporary configuration file: %v", err
	}

	err = tmpConfigFile.Close()
	if nil != err {
		return nil, "Unable to close temporary configuration file: %v", err
	}

	return tmpConfigFile, "", nil
}

//	scenario 01 - test to load configuration with invalid file name
func TestInitialize_scenario01(t *testing.T) {

	t.Run("load configuration with invalid file name", func(t *testing.T) {

		//	load the configuration file
		want := error(nil)
		testApp := App{}

		got := testApp.Initialize("invalid.config")
		if want != got {

			t.Errorf("invalid configuration loading result: got %v, want %v", got, want)
		}
	})
}

//	scenario 02 - test to load configuration without directories to perform some tests
func TestInitialize_scenario02(t *testing.T) {

	t.Run("load configuration without directories", func(t *testing.T) {

		configFile, msg, err := generateConfigWithoutDirectories()
		if nil != err {
			t.Fatalf(msg, err)
		}

		//	load the configuration file
		want := int32(1234)
		testApp := App{}

		err = testApp.Initialize(configFile.Name())
		if nil != err {

			t.Fatalf("Error loading configuration file: %v", err)
		}
		got := testApp.Configuration.PortNumber
		if want != got {

			t.Errorf("invalid port number: got %d, want %d", got, want)
		}
	})
}
