package main

import (
	"fmt"
	"log"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
)

func resolveNetworkDomain(apiClient *compute.Client, options programOptions) (networkDomain *compute.NetworkDomain, err error) {
	log.Printf("Resolve network domain '%s' in datacenter '%s'...",
		options.NetworkDomain,
		options.Datacenter,
	)

	networkDomain, err = apiClient.GetNetworkDomainByName(options.NetworkDomain, options.Datacenter)
	if err != nil {
		return
	}

	if networkDomain == nil {
		err = fmt.Errorf("Unable to find network domain '%s' in datacenter '%s'", options.NetworkDomain, options.Datacenter)
	}

	return
}
