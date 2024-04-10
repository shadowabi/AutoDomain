package zoomeye

import (
	"encoding/json"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
)

type ZoomeyeIpResponse struct {
	Matches []MatchInfo `json:"matches"`
}

type MatchInfo struct {
	Ip       string   `json:"ip"`
	PortInfo PortInfo `json:"portinfo"`
	Honeypot int      `json:"honeypot"`
}

type PortInfo struct {
	Port    int    `json:"port"`
	Service string `json:"service"`
}

func ParseZoomeyeIpResult(reqBody ...string) (zoomeyeIpRespList []ZoomeyeIpResponse) {
	if len(reqBody) != 0 {
		for _, response := range reqBody {
			var zoomeyeIpResponse ZoomeyeIpResponse
			Error.HandleError(json.Unmarshal([]byte(response), &zoomeyeIpResponse))
			if len(zoomeyeIpResponse.Matches) != 0 {
				for _, match := range zoomeyeIpResponse.Matches {
					if match.Honeypot == 0 {
						zoomeyeIpRespList = append(zoomeyeIpRespList, zoomeyeIpResponse)
					}
				}
			}
		}
	}
	return zoomeyeIpRespList
}

type ZoomeyeDomainResponse struct {
	List []SiteInfo `json:"list"`
}

type SiteInfo struct {
	Name string `json:"name"`
}

func ParseZoomeyeDomainResult(reqBody ...string) (zoomeyeDomainRespList []ZoomeyeDomainResponse) {
	if len(reqBody) != 0 {
		for _, response := range reqBody {
			var zoomeyeDomainResponse ZoomeyeDomainResponse
			Error.HandleError(json.Unmarshal([]byte(response), &zoomeyeDomainResponse))
			zoomeyeDomainRespList = append(zoomeyeDomainRespList, zoomeyeDomainResponse)
		}
	}
	return zoomeyeDomainRespList
}
