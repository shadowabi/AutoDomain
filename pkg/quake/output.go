package quake

import (
	"fmt"
	"strings"
)

func PurgeDomainResult(quakeResponse ...QuakeResponse) (quakeDomainResult []string) {
	if len(quakeResponse) != 0 {
		for _, response := range quakeResponse {
			for _, v := range response.Data {
				if v.Service.Name == "http" || v.Service.Name == "http/ssl" {
					if v.Service.Name == "http/ssl" {
						v.Service.Name = "https"
					}
					if v.Service.HTTP.Host == "" { // error data
						continue
					}
					result := strings.Join([]string{v.Service.Name, "://", v.Service.HTTP.Host, ":", fmt.Sprintf("%v", v.Port)}, "")
					quakeDomainResult = append(quakeDomainResult, result)
				}
			}
		}
	}
	return quakeDomainResult
}
