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
