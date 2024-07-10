package netlas

import (
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"net/http"
	"time"
)

func NetlasDomainRequest(client *http.Client, reqDomainList ...string) (responseDomainList []string) {
	if len(reqDomainList) != 0 {
		for _, host := range reqDomainList {
			url := fmt.Sprintf("https://app.netlas.io/api/domains/?q=*.%s&source_type=include&fields=domain", host)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", define.UserAgent)

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

func NetlasIpRequest(client *http.Client, reqIpList ...string) (responseIpList []string) {
	if len(reqIpList) != 0 {
		for _, host := range reqIpList {
			url := fmt.Sprintf("https://app.netlas.io/api/host/%s/?fields=domains", host)
			req, _ := http.NewRequest("GET", url, nil)

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
