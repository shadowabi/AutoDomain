package daydaymap

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"net/http"
	"time"
)

type DaydaymapData struct {
	Page    int    `json:"page"`
	Size    int    `json:"page_size"`
	Keyword string `json:"keyword"`
}

func DayDayMapRequest(client *http.Client, page int, total int, reqStringList ...string) (respBody []string) {
	if len(reqStringList) != 0 {
		for _, reqString := range reqStringList {
			reqString = base64.URLEncoding.EncodeToString([]byte(reqString))
			data := DaydaymapData{Page: page, Size: total, Keyword: reqString}
			dataJson, _ := json.Marshal(data)
			dataReq := bytes.NewBuffer(dataJson)
			req, _ := http.NewRequest("POST", "https://www.daydaymap.com/api/v1/raymap/search/all", dataReq)
			req.Header.Set("User-Agent", define.UserAgent)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("api-key", config.C.DaydaymapKey)

			resp, err := client.Do(req)
			time.Sleep(500 * time.Millisecond)
			Error.HandleError(err)
			respBody = append(respBody, net2.HandleResponse(resp))
			resp.Body.Close()
		}
	}
	return respBody
}
