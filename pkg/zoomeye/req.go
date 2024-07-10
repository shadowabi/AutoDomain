package zoomeye

import (
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"net/http"
	"time"
)

func ZoomeyeDomainRequest(client *http.Client, reqDomainList ...string) (responseDomainList []string) {
	if len(reqDomainList) != 0 {
		for _, host := range reqDomainList {
			url := fmt.Sprintf("https://api.zoomeye.org/domain/search?q=%s&type=1", host)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", define.UserAgent)
			req.Header.Set("API-KEY", config.C.ZoomeyeKey)

			resp, err := client.Do(req)
			time.Sleep(500 * time.Millisecond)
			if err != nil {
				continue
			}
			responseDomainList = append(responseDomainList, net2.HandleResponse(resp))
		}
	}
	return responseDomainList
}

func ZoomeyeIpRequest(client *http.Client, reqIpList ...string) (responseIpList []string) {
	if len(reqIpList) != 0 {
		for _, host := range reqIpList {
			url := fmt.Sprintf("https://api.zoomeye.org/host/search?query=ip:%s&facets=port", host)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Add("User-Agent", define.UserAgent)
			req.Header.Add("API-KEY", config.C.ZoomeyeKey)

			resp, err := client.Do(req)
			time.Sleep(500 * time.Millisecond)
			if err != nil {
				continue
			}
			responseIpList = append(responseIpList, net2.HandleResponse(resp))
		}
	}
	return responseIpList
}
