////////////////////////////////////////////////////////////////////////////////
//	main_test.go  -  Nov/25/2021  -  aldebap
//
//	httpServer package main unit tests
////////////////////////////////////////////////////////////////////////////////

package httpServer

import (
	"errors"
	"regexp"
	"testing"
)

//	Server.Initialize() scenario 01 - test to check correct error handling when startLineRegEx compilation fails
func TestInitialize_scenario01(t *testing.T) {

	t.Run("check correct error handling when startLineRegEx compilation fails", func(t *testing.T) {

		const errorMsg = "error compiling the regular expression"

		//	point regexCompile to a mock function
		regexCompile = func(expr string) (*regexp.Regexp, error) {

			if expr == `^(\S+)\s+(\S+)\s+(\S.*)$` {

				return nil, errors.New(errorMsg)
			}

			return regexp.Compile(expr)
		}

		//	load the configuration file
		want := errors.New(errorMsg)
		httpServer := Server{}

		got := httpServer.Initialize(5001)
		if want.Error() != got.Error() {

			t.Errorf("invalid regexp compilation: got %q, want %q", got, want)
		}
	})
}

//	Server.Initialize() scenario 02 - test to check correct error handling when requestHeaderRegEx compilation fails
func TestInitialize_scenario02(t *testing.T) {

	t.Run("check correct error handling when requestHeaderRegEx compilation fails", func(t *testing.T) {

		const errorMsg = "error compiling the regular expression"

		//	point regexCompile to a mock function
		regexCompile = func(expr string) (*regexp.Regexp, error) {

			if expr == `^(\S+):\s+(\S.*)$` {

				return nil, errors.New(errorMsg)
			}

			return regexp.Compile(expr)
		}

		//	load the configuration file
		want := errors.New(errorMsg)
		httpServer := Server{}

		got := httpServer.Initialize(5001)
		if want.Error() != got.Error() {

			t.Errorf("invalid regexp compilation: got %q, want %q", got, want)
		}
	})
}
