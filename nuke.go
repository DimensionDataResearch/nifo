package main

import (
	"log"

	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
)

// Destroy the target network domain.
func nuke(apiClient *compute.Client, networkDomainID string) error {
	err := nukeNATRules(apiClient, networkDomainID)
	if err != nil {
		return err
	}

	err = nukePublicIPBlocks(apiClient, networkDomainID)
	if err != nil {
		return err
	}

	err = nukeServers(apiClient, networkDomainID)
	if err != nil {
		return err
	}

	// TODO: Nuke VIP nodes, VIP pools, and virtual listeners.

	return nil
}

func nukeNATRules(apiClient *compute.Client, networkDomainID string) error {
	var natRules []compute.NATRule

	page := compute.DefaultPaging()
	page.PageSize = 20
	for {
		result, err := apiClient.ListNATRules(networkDomainID, page)
		if err != nil {
			return err
		}
		if result.IsEmpty() {
			break
		}

		natRules = append(natRules, result.Rules...)
	}

	for _, natRule := range natRules {
		log.Printf("Deleting NAT rule '%s' (%s -> %s)...",
			natRule.ID,
			natRule.ExternalIPAddress,
			natRule.InternalIPAddress,
		)

		err := apiClient.DeleteNATRule(natRule.ID)
		if err != nil {
			return err
		}

		log.Printf("Deleted NAT rule '%s' (%s -> %s).",
			natRule.ID,
			natRule.ExternalIPAddress,
			natRule.InternalIPAddress,
		)
	}

	return nil
}

func nukePublicIPBlocks(apiClient *compute.Client, networkDomainID string) error {
	var publicIPBlocks []compute.PublicIPBlock

	page := compute.DefaultPaging()
	page.PageSize = 20
	for {
		result, err := apiClient.ListPublicIPBlocks(networkDomainID, page)
		if err != nil {
			return err
		}
		if result.IsEmpty() {
			break
		}

		publicIPBlocks = append(publicIPBlocks, result.Blocks...)
	}

	for _, publicIPBlock := range publicIPBlocks {
		log.Printf("Deleting public IP block '%s'...",
			publicIPBlock.ID,
		)

		err := apiClient.RemovePublicIPBlock(publicIPBlock.ID)
		if err != nil {
			return err
		}

		log.Printf("Deleted public IP block '%s'...",
			publicIPBlock.ID,
		)
	}

	return nil
}

func nukeServers(apiClient *compute.Client, networkDomainID string) error {
	var servers []compute.Server

	page := compute.DefaultPaging()
	page.PageSize = 20
	for {
		result, err := apiClient.ListServersInNetworkDomain(networkDomainID, page)
		if err != nil {
			return err
		}
		if result.IsEmpty() {
			break
		}

		servers = append(servers, result.Items...)
	}

	// TODO: Parallelise.
	var err error
	for _, server := range servers {
		if server.Started {
			err = hardStopServer(apiClient, server.ID)
			if err != nil {
				return err
			}
		}

		log.Printf("Deleting server '%s' ('%s')...",
			server.Name,
			server.ID,
		)

		err = apiClient.DeleteServer(server.ID)
		if err != nil {
			return err
		}

		err = apiClient.WaitForDelete(compute.ResourceTypeServer, server.ID, 5*time.Minute)
		if err != nil {
			return err
		}

		log.Printf("Deleted server '%s' ('%s').",
			server.Name,
			server.ID,
		)
	}

	return nil
}

func hardStopServer(apiClient *compute.Client, serverID string) error {
	log.Printf("Stopping server '%s'...", serverID)

	err := apiClient.PowerOffServer(serverID)
	if err != nil {
		return err
	}

	_, err = apiClient.WaitForChange(compute.ResourceTypeServer, serverID, "Stop server", 5*time.Minute)
	if err != nil {
		return err
	}

	log.Printf("Stopped server '%s'...", serverID)

	return nil
}
