// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"flag"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/mock/server"
	"github.com/urfave/cli/v2"
)

var testRepoAppGet = cli.NewApp()

// setup the command for tests
func init() {
	testRepoAppGet.Commands = []*cli.Command{
		{
			Name: "get",
			Subcommands: []*cli.Command{
				&GetCmd,
			},
		},
	}
	testRepoAppGet.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "addr",
		},
		&cli.StringFlag{
			Name: "token",
		},
	}
}

func TestRepo_Get_Success(t *testing.T) {

	set := flag.NewFlagSet("test", 0)
	_ = cli.NewContext(testRepoAppGet, set, nil)

	// setup server
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	// setup types
	tests := []struct {
		data []string
		want error
	}{
		// default output
		{data: []string{
			"", "--addr", s.URL, "--token", "foobar",
			"get", "repo"}, want: nil},

		// json output
		{data: []string{
			"", "--addr", s.URL, "--token", "foobar",
			"get", "repo", "--o", "json"}, want: nil},

		// yaml output
		{data: []string{
			"", "--addr", s.URL, "--token", "foobar",
			"get", "repo", "--o", "yaml"}, want: nil},

		// wide output
		{data: []string{
			"", "--addr", s.URL, "--token", "foobar",
			"get", "repo", "--o", "wide"}, want: nil},

		// page default output
		{data: []string{
			"", "--addr", s.URL, "--token", "foobar",
			"get", "repo", "--p", "2"}, want: nil},

		// per page output
		{data: []string{
			"", "--addr", s.URL, "--token", "foobar",
			"get", "repo", "--pp", "20"}, want: nil},
	}

	// run test
	for _, test := range tests {
		got := testRepoAppGet.Run(test.data)

		if got != test.want {
			t.Errorf("Run is %v, want %v", got, test.want)
		}
	}
}

func TestRepo_Get_Failure(t *testing.T) {

	set := flag.NewFlagSet("test", 0)
	_ = cli.NewContext(testRepoAppGet, set, nil)

	// setup server
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	// setup types
	tests := []struct {
		data []string
		want error
	}{
		// ´Error with invalid addr
		{data: []string{
			"", "--token", "foobar",
			"get", "repo"},
			want: fmt.Errorf("Invalid command: Flag '--addr' is not set or is empty")},

		// ´Error with invalid token
		{data: []string{
			"", "--addr", s.URL,
			"get", "repo"},
			want: fmt.Errorf("Invalid command: Flag '--token' is not set or is empty")},
	}

	// run test
	for _, test := range tests {
		got := testRepoAppGet.Run(test.data)

		if got == test.want {
			t.Errorf("Run is %v, want %v", got, test.want)
		}
	}
}
