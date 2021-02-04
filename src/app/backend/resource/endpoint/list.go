package endpoint

import (
	v1 "k8s.io/api/core/v1"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
)

type EndpointList struct {
	ListMeta  api.ListMeta `json:"listMeta"`
	Endpoints []Endpoint   `json"endpoints"`
}

func toEndpointList(endpoints []v1.Endpoints) *EndpointList {
	endpointList := EndpointList{
		Endpoints: make([]Endpoint, 0),
		ListMeta:  api.ListMeta{TotalItems: len(endpoints)},
	}

	for _, endpoint := range endpoints {
		for _, subSets := range endpoint.Subsets {
			for _, address := range subSets.Addresses {
				endpointList.Endpoints = append(endpointList.Endpoints, *toEndpoint(address, subSets.Ports, true))
			}
			for _, notReadyAddress := range subSets.NotReadyAddresses {
				endpointList.Endpoints = append(endpointList.Endpoints, *toEndpoint(notReadyAddress, subSets.Ports, false))
			}
		}
	}

	return &endpointList
}
