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
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", 0)

func main() {
	log.SetPrefix("[VERBOSE] ")
	log.SetFlags(0) // No date / time prefix.

	options := parseOptions()

	if options.ShowHelp {
		showHelp()

		return
	} else if options.Version {
		fmt.Printf("NukeItFromOrbit %s\n", ProductVersion)

		return
	}

	if options.Verbose {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	apiClient, err := options.CreateClient()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	networkDomain, err := resolveNetworkDomain(apiClient, options)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if !options.Force {
		fmt.Printf("WARNING - about to delete network domain '%s' (Id = '%s') in datacenter '%s'. Are you sure you want to proceed?\n",
			networkDomain.Name,
			networkDomain.ID,
			networkDomain.DatacenterID,
		)
		fmt.Printf("Type yes to continue: ")
		stdin := bufio.NewReader(os.Stdin)
		confirmation, _, err := stdin.ReadLine()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if string(confirmation) != "yes" {
			os.Exit(2)
		}
	}

	err = nuke(apiClient, networkDomain.ID)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
