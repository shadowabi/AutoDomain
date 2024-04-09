package pulsedive

import (
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"net/http"
	"time"
)

func PulsediveDomainRequest(client *http.Client, reqDomainList ...string) (responseDomainList []string) {
	if len(reqDomainList) != 0 {
		for _, host := range reqDomainList {
			url := fmt.Sprintf("https://pulsedive.com/api/explore.php?q=ioc%%3d*.%s%%20%%20retired%%3d0&limit=50&key=%v", host, config.C.PulsediveKey)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", define.UserAgent)

			resp, err := client.Do(req)
			time.Sleep(500 * time.Millisecond)
			Error.HandleError(err)
			responseDomainList = append(responseDomainList, net2.HandleResponse(resp))
		}
	}
	return responseDomainList
}
