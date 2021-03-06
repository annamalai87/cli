// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/constants"

	"github.com/gosuri/uitable"

	"github.com/urfave/cli/v2"
	yaml "gopkg.in/yaml.v2"
)

// GetCmd defines the command for getting a list of repositories.
var GetCmd = cli.Command{
	Name:        "repo",
	Aliases:     []string{"repos"},
	Description: "Use this command to get a list of repositories.",
	Usage:       "Display a list of repositories",
	Action:      get,
	Flags: []cli.Flag{

		// optional flags that can be supplied to a command
		&cli.IntFlag{
			Name:    "page",
			Aliases: []string{"p"},
			Usage:   "Print a specific page of repos",
			Value:   1,
		},
		&cli.IntFlag{
			Name:    "per-page",
			Aliases: []string{"pp"},
			Usage:   "Expand the number of items contained within page",
			Value:   10,
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Print the output in wide, yaml or json format",
		},
	},
	CustomHelpTemplate: fmt.Sprintf(`%s
EXAMPLES:
 1. Get repositories.
    $ {{.HelpName}}
 2. Get repositories with wide view output.
    $ {{.HelpName}} --output wide
 3. Get repositories with yaml output.
    $ {{.HelpName}} --output yaml
 4. Get repositories with json output.
    $ {{.HelpName}} --output json
`, cli.CommandHelpTemplate),
}

// helper function to execute logs cli command
func get(c *cli.Context) error {
	// create a vela client
	client, err := vela.NewClient(c.String("addr"), nil)
	if err != nil {
		return err
	}

	// set token from global config
	client.Authentication.SetTokenAuth(c.String("token"))

	// set the page options based on user input
	opts := &vela.ListOptions{
		Page:    c.Int("page"),
		PerPage: c.Int("per-page"),
	}

	repositories, _, err := client.Repo.GetAll(opts)
	if err != nil {
		return err
	}

	switch c.String("output") {
	case "json":
		output, err := json.MarshalIndent(repositories, "", "    ")
		if err != nil {
			return err
		}

		fmt.Println(string(output))

	case "yaml":
		output, err := yaml.Marshal(repositories)
		if err != nil {
			return err
		}

		fmt.Println(string(output))

	case "wide":
		table := uitable.New()
		table.MaxColWidth = 200
		table.Wrap = true
		table.AddRow("ORG/REPO", "STATUS", "EVENTS", "VISIBILITY", "BRANCH", "REMOTE")

		for _, r := range *repositories {
			events := ""

			if r.GetAllowPush() {
				events += fmt.Sprintf("%s,", constants.EventPush)
			}

			if r.GetAllowPull() {
				events += fmt.Sprintf("%s,", constants.EventPull)
			}

			if r.GetAllowTag() {
				events += fmt.Sprintf("%s,", constants.EventTag)
			}

			if r.GetAllowDeploy() {
				events += fmt.Sprintf("%s,", constants.EventDeploy)
			}

			if r.GetAllowComment() {
				events += fmt.Sprintf("%s,", constants.EventComment)
			}

			table.AddRow(r.GetFullName(), r.GetActive(), strings.TrimSuffix(events, ","), r.GetVisibility(), r.GetBranch(), r.GetLink())
		}

		fmt.Println(table)

	default:
		table := uitable.New()
		table.MaxColWidth = 200
		table.Wrap = true // wrap columns

		table.AddRow("ORG/REPO", "STATUS", "EVENTS", "VISIBILITY", "BRANCH")

		for _, r := range *repositories {
			events := ""

			if r.GetAllowPush() {
				events += fmt.Sprintf("%s,", constants.EventPush)
			}

			if r.GetAllowPull() {
				events += fmt.Sprintf("%s,", constants.EventPull)
			}

			if r.GetAllowTag() {
				events += fmt.Sprintf("%s,", constants.EventTag)
			}

			if r.GetAllowDeploy() {
				events += fmt.Sprintf("%s,", constants.EventDeploy)
			}

			if r.GetAllowComment() {
				events += fmt.Sprintf("%s,", constants.EventComment)
			}

			table.AddRow(r.GetFullName(), r.GetActive(), strings.TrimSuffix(events, ","), r.GetVisibility(), r.GetBranch())
		}

		fmt.Println(table)
	}

	return nil
}
