package fofa

import (
	"encoding/base64"
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"net/http"
	"time"
)

func FofaRequest(client *http.Client, page int, reqStringList ...string) (respBody []string) {
	if len(reqStringList) != 0 {
		for _, reqString := range reqStringList {
			reqString = base64.URLEncoding.EncodeToString([]byte(reqString))
			url := fmt.Sprintf("https://fofa.info/api/v1/search/all?key=%s&qbase64=%s&full=false&fields=protocol,host&size=1000&page=%v", config.C.FofaKey, reqString, page)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", define.UserAgent)

			resp, err := client.Do(req)
			time.Sleep(500 * time.Millisecond)
			if err != nil {
				continue
			}
			respBody = append(respBody, net2.HandleResponse(resp))
			resp.Body.Close()
		}
	}
	return respBody
}
