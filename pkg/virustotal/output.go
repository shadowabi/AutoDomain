package virustotal

import "strings"

func PurgeDomainResult(virusTotalDomainResponse ...VirusTotalDomainResponse) (virusTotalDomainResult []string) {
	if len(virusTotalDomainResponse) != 0 {
		for _, response := range virusTotalDomainResponse {
			for _, v := range response.Data {
				result := strings.Join([]string{"http://", v.ID}, "")
				virusTotalDomainResult = append(virusTotalDomainResult, result)
			}
		}
	}
	return virusTotalDomainResult
}
