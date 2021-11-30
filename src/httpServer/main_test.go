////////////////////////////////////////////////////////////////////////////////
//	main_test.go  -  Nov/25/2021  -  aldebap
//
//	httpServer package main unit tests
////////////////////////////////////////////////////////////////////////////////

package httpServer

import (
	"errors"
	"net"
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
		regexCompile = regexp.Compile

		if want.Error() != got.Error() {

			t.Errorf("invalid error returned: got %q, want %q", got, want)
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
		regexCompile = regexp.Compile

		if want.Error() != got.Error() {

			t.Errorf("invalid error returned: got %q, want %q", got, want)
		}
	})
}

//	Server.ListenAndServe() scenario 01 - test to check correct error handling when Initialize() is not called before
func TestListenAndServe_scenario01(t *testing.T) {

	t.Run("check correct error handling when Initialize() is not called before", func(t *testing.T) {

		const errorMsg = "Server object must be initialized before Listening and Serving"

		want := errors.New(errorMsg)
		httpServer := Server{}

		got := httpServer.ListenAndServe()
		if want.Error() != got.Error() {

			t.Errorf("invalid error returned: got %q, want %q", got, want)
		}
	})
}

//	Server.ListenAndServe() scenario 02 - test to check correct error handling when net.ResolveTCPAddr() call fails
func TestListenAndServe_scenario02(t *testing.T) {

	t.Run("test to check correct error handling when net.ResolveTCPAddr() call fails", func(t *testing.T) {

		const errorMsg = "Cannot resolve TCP address"

		//	point resolveTCPAddr to a mock function
		resolveTCPAddr = func(network string, address string) (*net.TCPAddr, error) {

			return nil, errors.New(errorMsg)
		}

		want := errors.New(errorMsg)
		httpServer := Server{}

		httpServer.Initialize(5001)
		got := httpServer.ListenAndServe()
		resolveTCPAddr = net.ResolveTCPAddr

		if want.Error() != got.Error() {

			t.Errorf("invalid error returned: got %q, want %q", got, want)
		}
	})
}

//	Server.ListenAndServe() scenario 03 - test to check correct error handling when net.ListenTCP() call fails
func TestListenAndServe_scenario03(t *testing.T) {

	t.Run("test to check correct error handling when net.ListenTCP() call fails", func(t *testing.T) {

		const errorMsg = "Cannot listen the port of TCP address"

		//	point listenTCP to a mock function
		listenTCP = func(network string, laddr *net.TCPAddr) (*net.TCPListener, error) {

			return nil, errors.New(errorMsg)
		}

		want := errors.New(errorMsg)
		httpServer := Server{}

		httpServer.Initialize(5001)
		got := httpServer.ListenAndServe()
		listenTCP = net.ListenTCP

		if want.Error() != got.Error() {

			t.Errorf("invalid error returned: got %q, want %q", got, want)
		}
	})
}
