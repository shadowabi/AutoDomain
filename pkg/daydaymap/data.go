package daydaymap

import (
	"encoding/base64"
	"encoding/json"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	"strings"
)

func MergeReqListToReqString(reqIpList []string, reqDomainList []string) (reqString string) {
	for _, host := range reqIpList {
		reqString += "ip:" + "\"" + host + "\"" + " || "
	}
	for _, host := range reqDomainList {
		reqString += "domain:" + "\"" + host + "\"" + " || "
	}
	reqString = strings.TrimSuffix(reqString, " || ")
	reqString = base64.URLEncoding.EncodeToString([]byte(reqString))
	return reqString
}

type ArrElement struct {
	Domain  string `json:"domain"`
	Port    int    `json:"port"`
	Ip      string `json:"ip"`
	Service string `json:"service"`
}

type Data struct {
	List  []ArrElement `json:"list"`
	Total int          `json:"total"`
}

type DaydaymapResponse struct {
	Data Data `json:"data"`
}

func ParseDaydaymapResult(reqBody ...string) (daydaymapRespList []DaydaymapResponse) {
	if len(reqBody) != 0 {
		for _, response := range reqBody {
			var daydaymapResponse DaydaymapResponse
			Error.HandleError(json.Unmarshal([]byte(response), &daydaymapResponse))
			daydaymapRespList = append(daydaymapRespList, daydaymapResponse)
		}
	}
	return daydaymapRespList
}
