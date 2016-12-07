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
