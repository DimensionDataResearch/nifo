package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
)

// Destroy the target network domain.
func nuke(apiClient *compute.Client, networkDomainID string) error {
	logger.Printf("Destroying network domain '%s'...", networkDomainID)

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

	err = nukeVLANs(apiClient, networkDomainID)
	if err != nil {
		return err
	}

	err = nukeNetworkDomain(apiClient, networkDomainID)
	if err != nil {
		return err
	}

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

		page.Next()
	}

	for _, natRule := range natRules {
		logger.Printf("Deleting NAT rule '%s' (%s -> %s)...",
			natRule.ID,
			natRule.ExternalIPAddress,
			natRule.InternalIPAddress,
		)

		err := apiClient.DeleteNATRule(natRule.ID)
		if err != nil {
			return err
		}

		logger.Printf("Deleted NAT rule '%s' (%s -> %s).",
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

		page.Next()
	}

	for _, publicIPBlock := range publicIPBlocks {
		logger.Printf("Deleting public IP block '%s'...",
			publicIPBlock.ID,
		)

		err := apiClient.RemovePublicIPBlock(publicIPBlock.ID)
		if err != nil {
			return err
		}

		logger.Printf("Deleted public IP block '%s'...",
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

		page.Next()
	}

	asyncLock := &sync.Mutex{}
	deletionComplete := &sync.WaitGroup{}
	deletionComplete.Add(len(servers))

	failed := false
	for _, server := range servers {
		go func(server compute.Server) {
			defer deletionComplete.Done()

			var err error
			if server.Started {
				err = hardStopServer(apiClient, server.ID)
				if err != nil {
					logger.Println(err)
					failed = true

					return
				}
			}

			asyncLock.Lock()
			logger.Printf("Destroying server '%s' ('%s')...",
				server.Name,
				server.ID,
			)

			err = apiClient.DeleteServer(server.ID)
			asyncLock.Unlock()
			if err != nil {
				logger.Println(err)
				failed = true

				return
			}

			err = apiClient.WaitForDelete(compute.ResourceTypeServer, server.ID, 5*time.Minute)
			if err != nil {
				logger.Println(err)
				failed = true

				return
			}

			logger.Printf("Destroyed server '%s' ('%s').",
				server.Name,
				server.ID,
			)

		}(server)
	}

	deletionComplete.Wait()
	if failed {
		return fmt.Errorf("Destroy failed for one or more servers in network domain '%s'.", networkDomainID)
	}

	return nil
}

func hardStopServer(apiClient *compute.Client, serverID string) error {
	logger.Printf("Stopping server '%s'...", serverID)

	err := apiClient.PowerOffServer(serverID)
	if err != nil {
		return err
	}

	_, err = apiClient.WaitForChange(compute.ResourceTypeServer, serverID, "Stop server", 5*time.Minute)
	if err != nil {
		return err
	}

	logger.Printf("Stopped server '%s'...", serverID)

	return nil
}

func nukeVLANs(apiClient *compute.Client, networkDomainID string) error {
	var vlans []compute.VLAN

	page := compute.DefaultPaging()
	page.PageSize = 20
	for {
		result, err := apiClient.ListVLANs(networkDomainID, page)
		if err != nil {
			return err
		}
		if result.IsEmpty() {
			break
		}

		vlans = append(vlans, result.VLANs...)

		page.Next()
	}

	for _, vlan := range vlans {
		logger.Printf("Deleting VLAN '%s'...",
			vlan.ID,
		)

		err := apiClient.DeleteVLAN(vlan.ID)
		if err != nil {
			return err
		}

		err = apiClient.WaitForDelete(compute.ResourceTypeVLAN, vlan.ID, 5*time.Minute)
		if err != nil {
			return err
		}

		logger.Printf("Deleted VLAN '%s'...",
			vlan.ID,
		)
	}

	return nil
}

func nukeNetworkDomain(apiClient *compute.Client, networkDomainID string) error {
	logger.Printf("Deleting network domain '%s'...", networkDomainID)

	err := apiClient.DeleteNetworkDomain(networkDomainID)
	if err != nil {
		return err
	}

	logger.Printf("Deleted network domain '%s'.", networkDomainID)

	return nil
}
