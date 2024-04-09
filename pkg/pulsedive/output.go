package pulsedive

import "strings"

func PurgeDomainResult(pulsediveDomainResponse ...PulsediveDomainResponse) (pulsediveDomainResult []string) {
	if len(pulsediveDomainResponse) != 0 {
		for _, response := range pulsediveDomainResponse {
			for _, v := range response.Results {
				result := strings.Join([]string{"http://", v.Indicator}, "")
				pulsediveDomainResult = append(pulsediveDomainResult, result)
			}
		}
	}
	return pulsediveDomainResult
}
