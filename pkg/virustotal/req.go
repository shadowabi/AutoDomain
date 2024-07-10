package virustotal

import (
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"net/http"
	"time"
)

func VirusTotalDomainRequest(client *http.Client, reqDomainList ...string) (responseDomainList []string) {
	if len(reqDomainList) != 0 {
		for _, host := range reqDomainList {
			url := fmt.Sprintf("https://www.virustotal.com/ui/domains/%s/subdomains?relationships=resolutions&cursor=eyJsaW1pdCI6IDIwMCwgIm9mZnNldCI6IDB9&limit=200", host)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", define.UserAgent)
			req.Header.Set("X-Vt-Anti-Abuse-Header", "1")
			req.Header.Set("X-Tool", "vt-ui-main")
			req.Header.Set("Accept-Ianguage", "en-US,en;q=0.9,es;q=0.8")

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
