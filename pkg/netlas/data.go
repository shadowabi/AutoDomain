package netlas

import (
	"encoding/json"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
)

type NetlasDomainResponse struct {
	Items     []DomainData `json:"items"`
	Took      int          `json:"took"`
	Timestamp int64        `json:"timestamp"`
}

type DomainData struct {
	Data Domain `json:"data"`
}

type Domain struct {
	Domain string `json:"domain"`
}

func ParseNetlasDomainResult(reqBody ...string) (netlasDomainRespList []NetlasDomainResponse) {
	if len(reqBody) != 0 {
		for _, response := range reqBody {
			var netlasDomainResponse NetlasDomainResponse
			Error.HandleError(json.Unmarshal([]byte(response), &netlasDomainResponse))
			netlasDomainRespList = append(netlasDomainRespList, netlasDomainResponse)
		}
	}
	return netlasDomainRespList
}

type NetlasIpResponse struct {
	Domains []string `json:"domains"`
}

func ParseNetlasIpResult(reqBody ...string) (netlasIpRespList []NetlasIpResponse) {
	if len(reqBody) != 0 {
		for _, response := range reqBody {
			var netlasDomainResponse NetlasIpResponse
			Error.HandleError(json.Unmarshal([]byte(response), &netlasDomainResponse))
			netlasIpRespList = append(netlasIpRespList, netlasDomainResponse)
		}
	}
	return netlasIpRespList
}
