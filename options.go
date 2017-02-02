/*
   Copyright 2016 Dimension Data

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/jessevdk/go-flags"
)

type programOptions struct {
	Region        string `short:"r" long:"region" description:"The CloudControl region to use (e.g. AU, NA, etc)."`
	Datacenter    string `short:"d" long:"datacenter" description:"The name CloudControl data centre containing the resource(s) to export (e.g. AU10, NA9, etc)."`
	NetworkDomain string `short:"n" long:"networkdomain" description:"The name of tje network domain to nuke."`
	Force         bool   `short:"f" long:"force" description:"Destroy the network domain without prompting."`
	Verbose       bool   `short:"v" long:"verbose" description:"Display detailed information about the program's activities."`
	Version       bool   `long:"version" description:"Display program version info."`
	ShowHelp      bool   `short:"?" long:"help" description:"Show program help."`
}

// Validate the programOptions.
func (options programOptions) Validate() error {
	if options.Region == "" {
		return fmt.Errorf("Must specify the target region.")
	}

	if options.Datacenter == "" {
		return fmt.Errorf("Must specify the target datacenter.")
	}

	if options.NetworkDomain == "" {
		return fmt.Errorf("Must specify the target network domain.")
	}

	return nil
}

// Create a CloudControl client.
func (options programOptions) CreateClient() (client *compute.Client, err error) {
	if options.Region == "" {
		err = fmt.Errorf("Must specify the target CloudControl region.")

		return
	}

	username := os.Getenv("MCP_USER")
	if username == "" {
		err = fmt.Errorf("The MCP_USER environment variable has not been set. Set it to your CloudControl username.")

		return
	}

	password := os.Getenv("MCP_PASSWORD")
	if username == "" {
		err = fmt.Errorf("The MCP_PASSWORD environment variable has not been set. Set it to your CloudControl password.")

		return
	}

	client = compute.NewClient(options.Region, username, password)

	return
}

func parseOptions() programOptions {
	options := programOptions{}

	parser := flags.NewParser(&options, flags.Default)
	_, err := parser.ParseArgs(os.Args)
	if err == nil {
		err = options.Validate()
	}
	if err != nil {
		switch err.(type) {
		case *flags.Error:
			// Ignore, since the parser will print it out anyway.
		default:
			fmt.Print(err.Error())
		}

		showHelp()

		os.Exit(1)
	}

	return options
}

func showHelp() {
	options := programOptions{}

	parser := flags.NewParser(&options, flags.Default)
	fmt.Println()
	parser.WriteHelp(os.Stdout)
}
