package hunter

import (
	"encoding/json"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
)

type ArrElement struct {
	Url string `json:"url"`
}

type Data struct {
	Total int          `json:"time"`
	Arr   []ArrElement `json:"arr"`
}

type HunterResponse struct {
	Data Data `json:"data"`
}

func ParseHunterResult(reqBody ...string) (hunterRespList []HunterResponse) {
	for _, response := range reqBody {
		var hunterResponse HunterResponse
		Error.HandleError(json.Unmarshal([]byte(response), &hunterResponse))
		hunterRespList = append(hunterRespList, hunterResponse)
	}
	return hunterRespList
}
