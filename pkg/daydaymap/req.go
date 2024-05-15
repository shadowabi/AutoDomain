package daydaymap

import (
	"bytes"
	"encoding/json"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"net/http"
	"time"
)

func DayDayMapRequest(client *http.Client, reqString string, totalList ...int) (respBody []string) {
	if len(totalList) != 0 {
		for _, total := range totalList {
			data := struct {
				Page    int    `json:"page"`
				Size    int    `json:"page_size"`
				Keyword string `json:"keyword"`
			}{
				Page:    1,
				Size:    total,
				Keyword: reqString,
			}
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
