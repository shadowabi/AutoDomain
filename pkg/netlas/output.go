package netlas

import "strings"

func PurgeDomainResult(netlasDomainResponse ...NetlasDomainResponse) (netlasDomainResult []string) {
	if len(netlasDomainResponse) != 0 {
		for _, response := range netlasDomainResponse {
			for _, v := range response.Items {
				result := strings.Join([]string{"http://", v.Data.Domain}, "")
				netlasDomainResult = append(netlasDomainResult, result)
			}
		}
	}
	return netlasDomainResult
}

func PurgeIpResult(netlasIpesponse ...NetlasIpResponse) (netlasIpResult []string) {
	if len(netlasIpesponse) != 0 {
		for _, response := range netlasIpesponse {
			for _, v := range response.Domains {
				result := strings.Join([]string{"http://", v}, "")
				netlasIpResult = append(netlasIpResult, result)
			}
		}
	}
	return netlasIpResult
}
