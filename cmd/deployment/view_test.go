// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package deployment

import (
	"flag"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/mock/server"
	"github.com/urfave/cli/v2"
)

var testDeploymentAppView = cli.NewApp()

// setup the command for tests
func init() {
	testDeploymentAppView.Commands = []*cli.Command{
		{
			Name: "view",
			Subcommands: []*cli.Command{
				&ViewCmd,
			},
		},
	}
	testDeploymentAppView.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "addr",
		},
		&cli.StringFlag{
			Name: "token",
		},
	}
}

func TestDeployment_View_Success(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	_ = cli.NewContext(testDeploymentAppView, set, nil)

	// setup server
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	// setup types
	tests := []struct {
		data []string
		want error
	}{
		{ // default output
			data: []string{"",
				"--addr", s.URL,
				"--token", "foobar",
				"view",
				"deployment",
				"--org", "github",
				"--repo", "octocat",
				"--deployment", "1",
			},
			want: nil,
		},
		{ // json output
			data: []string{"",
				"--addr", s.URL,
				"--token", "foobar",
				"view",
				"deployment",
				"--org", "github",
				"--repo", "octocat",
				"--deployment", "1",
				"--output", "json",
			},
			want: nil,
		},
	}

	// run test
	for _, test := range tests {
		got := testDeploymentAppView.Run(test.data)

		if got != test.want {
			t.Errorf("Run is %v, want %v", got, test.want)
		}
	}
}

func TestDeployment_View_Failure(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	_ = cli.NewContext(testDeploymentAppView, set, nil)

	// setup server
	gin.SetMode(gin.TestMode)

	s := httptest.NewServer(server.FakeHandler())

	// setup types
	tests := []struct {
		data []string
		want error
	}{
		{ // Error with invalid addr
			data: []string{"",
				"--token", "foobar",
				"view",
				"deployment",
				"--org", "github",
				"--repo", "octocat",
				"--deployment", "1",
			},
			want: fmt.Errorf("Invalid command: Flag '--addr' is not set or is empty"),
		},
		{ // Error with invalid token
			data: []string{"",
				"--addr", s.URL,
				"view",
				"deployment",
				"--org", "github",
				"--repo", "octocat",
				"--deployment", "1",
			},
			want: fmt.Errorf("Invalid command: Flag '--token' is not set or is empty"),
		},
		{ // Error with invalid org
			data: []string{"",
				"--addr", s.URL,
				"view",
				"deployment",
				"--repo", "octocat",
				"--deployment", "1",
			},
			want: fmt.Errorf("Invalid command: Flag '--org' is not set or is empty"),
		},
		{ // Error with invalid repo
			data: []string{"",
				"--addr", s.URL,
				"view",
				"deployment",
				"--org", "github",
				"--deployment", "1",
			},
			want: fmt.Errorf("Invalid command: Flag '--repo' is not set or is empty"),
		},
		{ // Error with invalid deployment
			data: []string{"",
				"--addr", s.URL,
				"view",
				"deployment",
				"--org", "github",
				"--repo", "octocat",
			},
			want: fmt.Errorf("Invalid command: Flag '--deployment' is not set or is empty"),
		},
	}

	// run test
	for _, test := range tests {
		got := testDeploymentAppView.Run(test.data)

		if got == test.want {
			t.Errorf("Run is %v, want %v", got, test.want)
		}
	}
}
