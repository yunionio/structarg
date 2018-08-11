package main

import (
	"fmt"
	"os"

	"yunion.io/x/structarg"
)

// define default arguments
type Options struct {
	Help         bool   `help:"Show help" short-token:"h"`
	Debug        bool   `help:"Show debug information"`
	Timeout      int    `default:"600" help:"Maximal number of seconds to wait for a response"`
	AuthURLStr   string `default:"$AUTH_URL" help:"Authentication URL, default to env[AUTH_URL]"`
	EndpointType string `default:"publicURL" help:"Default to env[ENPOINT_TYPE] or publicURL" choices:"publicURL|internalURL"`
	Config       string `help:"Configuration file path"`
	SUBCOMMAND   string `help:"climc subcommand" subcommand:"true"`
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
	e = parser.ParseArgs(os.Args[1:], false)
	options := parser.Options().(*Options)
	if len(options.Config) > 0 {
		ec := parser.ParseFile(options.Config)
		if ec != nil {
			fmt.Printf("Error reading config file: %s\n", ec)
			os.Exit(1)
		}
	}
	if options.Help {
		fmt.Print(parser.HelpString())
	} else {
		fmt.Printf("################## Options #################\n")
		fmt.Printf("AuthURLStr = %s\n", options.AuthURLStr)
		fmt.Printf("Timeout = %d\n", options.Timeout)
		fmt.Printf("EndpointType = %s\n", options.EndpointType)
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
