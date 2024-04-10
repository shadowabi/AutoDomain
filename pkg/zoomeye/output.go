package zoomeye

import (
	"fmt"
	"strings"
)

func PurgeIpResult(zoomeyeIpResponse ...ZoomeyeIpResponse) (zoomeyeDomainResult []string) {
	if len(zoomeyeIpResponse) != 0 {
		for _, response := range zoomeyeIpResponse {
			if len(response.Matches) != 0 {
				for _, v := range response.Matches {
					if v.PortInfo.Service == "http" || v.PortInfo.Service == "https" {
						result := strings.Join([]string{v.PortInfo.Service, "://", v.Ip, ":", fmt.Sprintf("%v", v.PortInfo.Port)}, "")
						zoomeyeDomainResult = append(zoomeyeDomainResult, result)
					}
				}
			}
		}
	}
	return zoomeyeDomainResult
}

func PurgeDomainResult(zoomeyeDomainResponse ...ZoomeyeDomainResponse) (zoomeyeDomainResult []string) {
	if len(zoomeyeDomainResponse) != 0 {
		for _, response := range zoomeyeDomainResponse {
			for _, v := range response.List {
				result := strings.Join([]string{"http://", v.Name}, "")
				zoomeyeDomainResult = append(zoomeyeDomainResult, result)
			}
		}
	}
	return zoomeyeDomainResult
}
