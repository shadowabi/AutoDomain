package daydaymap

import (
	"fmt"
)

func PurgeDomainResult(daydaymapResponse ...DaydaymapResponse) (daydaymapDomainResult []string) {
	if len(daydaymapResponse) != 0 {
		for _, response := range daydaymapResponse {
			for _, v := range response.Data.List {
				if v.Service == "http" || v.Service == "https" {
					result := ""
					if v.Domain != "" {
						result = fmt.Sprintf("%s://%s:%v", v.Service, v.Domain, v.Port)
					} else {
						result = fmt.Sprintf("%s://%v:%v", v.Service, v.Ip, v.Port)
					}
					daydaymapDomainResult = append(daydaymapDomainResult, result)
				}
			}
		}
	}
	return daydaymapDomainResult
}
