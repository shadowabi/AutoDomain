package fofa

import (
	"encoding/json"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
)

type FofaResponse struct {
	Size    int        `json:"size"`
	Results [][]string `json:"results"`
}

func ParseFofaResult(reqBody ...string) (fofaRespList []FofaResponse) {
	for _, response := range reqBody {
		var fofaResponse FofaResponse
		Error.HandleError(json.Unmarshal([]byte(response), &fofaResponse))
		fofaRespList = append(fofaRespList, fofaResponse)
	}
	return fofaRespList
}
