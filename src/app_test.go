////////////////////////////////////////////////////////////////////////////////
//	app_test.go  -  Nov/18/2021  -  aldebap
//
//	App unit tests
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

//	generate invalid Json configuration file
func generateInvalidConfig() (*os.File, string, error) {

	//	create an invalid temporary configuration file
	tmpConfigFile, err := ioutil.TempFile(os.TempDir(), "config")
	if err != nil {
		return nil, "Unable to create temporary configuration file: %v", err
	}

	//	write the configuration to file
	configurationData := []byte(`{
			"portNumber": 1234
		`)
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

//	generate configuration without directories to perform some tests
func generateConfigWithoutDirectories() (*os.File, string, error) {

	//	create a temporary configuration file without directories
	tmpConfigFile, err := ioutil.TempFile(os.TempDir(), "config")
	if err != nil {
		return nil, "Unable to create temporary configuration file: %v", err
	}

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
		want := errors.New("open invalid.config: no such file or directory")
		testApp := App{}

		got := testApp.Initialize("invalid.config")
		if want.Error() != got.Error() {

			t.Errorf("invalid configuration loading result: got %q, want %q", got, want)
		}
	})
}

//	scenario 02 - test error handling on attempt to load configuration file
func TestInitialize_scenario02(t *testing.T) {

	t.Run("error handling on attempt to load configuration file", func(t *testing.T) {

		//	point readAll to a mock function
		readAll = func(io.Reader) ([]byte, error) {

			return nil, errors.New("cannot read the input file")
		}

		tmpConfigFile, msg, err := generateInvalidConfig()
		if nil != err {
			t.Fatalf(msg, err)
		}
		defer os.Remove(tmpConfigFile.Name())

		//	try to load the configuration file
		want := errors.New("cannot read the input file")
		testApp := App{}

		got := testApp.Initialize(tmpConfigFile.Name())
		if want.Error() != got.Error() {

			t.Errorf("invalid configuration loading result: got %q, want %q", got, want)
		}

		//	point readAll back to io.ReadAll
		readAll = ioutil.ReadAll
	})
}

//	scenario 03 - test error handling on attempt to parse JSon content
func TestInitialize_scenario03(t *testing.T) {

	t.Run("error handling on attempt to parse JSon content", func(t *testing.T) {

		tmpConfigFile, msg, err := generateInvalidConfig()
		if nil != err {
			t.Fatalf(msg, err)
		}
		defer os.Remove(tmpConfigFile.Name())

		//	try to load the configuration file
		want := errors.New("unexpected end of JSON input")
		testApp := App{}

		got := testApp.Initialize(tmpConfigFile.Name())
		if want.Error() != got.Error() {

			t.Errorf("invalid configuration loading result: got %q, want %q", got, want)
		}
	})
}

//	scenario 04 - test to load configuration without directories to perform some tests
func TestInitialize_scenario04(t *testing.T) {

	t.Run("load configuration without directories", func(t *testing.T) {

		tmpConfigFile, msg, err := generateConfigWithoutDirectories()
		if nil != err {
			t.Fatalf(msg, err)
		}
		defer os.Remove(tmpConfigFile.Name())

		//	load the configuration file
		want := int32(1234)
		testApp := App{}

		err = testApp.Initialize(tmpConfigFile.Name())
		if nil != err {

			t.Fatalf("Error loading configuration file: %v", err)
		}
		got := testApp.Configuration.PortNumber
		if want != got {

			t.Errorf("invalid port number: got %d, want %d", got, want)
		}
	})
}
