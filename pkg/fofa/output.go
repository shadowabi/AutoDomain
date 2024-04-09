package fofa

import (
	"strings"
)

func PurgeDomainResult(fofaResponse ...FofaResponse) (fofaDomainResult []string) {
	if len(fofaResponse) != 0 {
		for _, response := range fofaResponse {
			for _, v := range response.Results {
				if v[0] == "http" || v[1] == "https" {
					result := strings.Join([]string{v[0], "://", v[1]}, "")
					fofaDomainResult = append(fofaDomainResult, result)
				}
			}
		}
	}
	return fofaDomainResult
}
