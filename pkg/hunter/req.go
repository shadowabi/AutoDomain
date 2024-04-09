package hunter

import (
	"encoding/base64"
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"net/http"
	"time"
)

func HunterRequest(client *http.Client, reqString string, page ...int) (respBody []string) {
	if len(page) != 0 {
		reqString = base64.URLEncoding.EncodeToString([]byte(reqString))
		for _, num := range page {
			url := fmt.Sprintf("https://hunter.qianxin.com/openApi/search?api-key=%s&search=%s&page=%v&page_size=100&is_web=3",
				config.C.HunterKey, reqString, num)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", define.UserAgent)

			resp, err := client.Do(req)
			time.Sleep(500 * time.Millisecond)
			Error.HandleError(err)
			respBody = append(respBody, net2.HandleResponse(resp))
			resp.Body.Close()
		}
	}
	return respBody
}
