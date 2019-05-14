// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"yunion.io/x/structarg"
)

// define default arguments
type Options struct {
	structarg.BaseOptions

	Region  string `help:"Region name or ID"`
	Port    int    `help:"The port that the service runs on"`
	Address string `help:"The IP address to serve on (set to 0.0.0.0 for all interfaces)" default:"0.0.0.0"`

	AuthURL       string   `help:"Keystone auth URL" alias:"auth-uri"`
	AdminUser     string   `help:"Admin username"`
	AdminDomain   string   `help:"Admin user domain"`
	AdminPassword string   `help:"Admin password"`
	AdminProject  string   `help:"Admin project" default:"system" alias:"admin-tenant-name"`
	CorsHosts     []string `help:"List of hostname that allow CORS"`

	SqlConnection string `help:"SQL connection string"`

	DNSServer    string   `help:"Address of DNS server"`
	DNSDomain    string   `help:"Domain suffix for virtual servers"`
	DNSResolvers []string `help:"Upstream DNS resolvers"`

	Debug        bool     `help:"Show debug information"`
	Timeout      int      `default:"600" help:"Maximal number of seconds to wait for a response"`
	AuthURLStr   string   `default:"$AUTH_URL" help:"Authentication URL, default to env[AUTH_URL]"`
	EndpointType string   `default:"publicURL" help:"Default to env[ENPOINT_TYPE] or publicURL" choices:"publicURL|internalURL"`
	Endpoints    []string `help:"endpoints" json:"end-point" default:"e1,e2"`
	SUBCOMMAND   string   `help:"climc subcommand" subcommand:"true"`
}

// argument
type HelpOptions struct {
	SUBCOMMAND string `help:"Sub-command name"`
}

type TestOptions struct {
	NAME string `help:"Test name"`
	Arg1 string `help:"Argument1"`
	Arg2 string `help:"Argument1"`
}

func showErrorAndExit(e error) {
	fmt.Printf("Error: %s\n", e)
	os.Exit(1)
}

func main() {
	parser, e := structarg.NewArgumentParser(&Options{},
		"structargtest",
		`Command-line interface test prog`,
		`See "structargtest help COMMAND" for help on a subcommand.`)
	subcmd := parser.GetSubcommand()
	if subcmd == nil {
		showErrorAndExit(fmt.Errorf("No subcommand argument"))
	}
	// add subcomamnd
	subcmd.AddSubParser(&HelpOptions{}, "help", "Show help information of a subcommand", func(suboptions *HelpOptions) error {
		helpstr, e := subcmd.SubHelpString(suboptions.SUBCOMMAND)
		if e != nil {
			return e
		} else {
			fmt.Print(helpstr)
			return nil
		}
	})
	subcmd.AddSubParser(&TestOptions{}, "test", "Run a test", func(suboptions *TestOptions) error {
		fmt.Printf("Run test %s with argument \"%s\" and \"%s\"\n", suboptions.NAME, suboptions.Arg1, suboptions.Arg2)
		return nil
	})
	subcmd.AddSubParser(&TestOptions{}, "test2", "Run a test 2", func(suboptions *TestOptions) error {
		fmt.Printf("Run test %s with argument \"%s\" and \"%s\"\n", suboptions.NAME, suboptions.Arg1, suboptions.Arg2)
		return nil
	})
	e = parser.ParseArgs2(os.Args[1:], false, false)
	options := parser.Options().(*Options)
	if len(options.Config) > 0 {
		ec := parser.ParseFile(options.Config)
		if ec != nil {
			fmt.Printf("Error reading config file: %s\n", ec)
			os.Exit(1)
		}
	}
	parser.SetDefault()

	if options.Help {
		fmt.Print(parser.HelpString())
	} else {
		fmt.Printf("################## Options #################\n")
		fmt.Printf("%#v\n", options)
		fmt.Printf("############################################\n")
		subcmd := parser.GetSubcommand()
		if subcmd == nil {
			if e != nil {
				showErrorAndExit(e)
			}
		} else {
			subparser := subcmd.GetSubParser()
			if e != nil {
				if subparser != nil {
					fmt.Print(subparser.Usage())
				}
				showErrorAndExit(e)
			} else {
				suboptions := subparser.Options()
				e = subcmd.Invoke(suboptions)
				if e != nil {
					showErrorAndExit(e)
				}
			}
		}
	}
}
