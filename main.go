package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
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
	}

	apiClient, err := options.CreateClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	networkDomain, err := resolveNetworkDomain(apiClient, options)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("WARNING - about to delete network domain '%s' (Id = '%s') in datacenter '%s'. Are you sure you want to proceed?",
		networkDomain.Name,
		networkDomain.ID,
		networkDomain.DatacenterID,
	)
}
