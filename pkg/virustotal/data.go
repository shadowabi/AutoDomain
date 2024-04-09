package virustotal

import (
	"encoding/json"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
)

type VirusTotalDomain struct {
	ID string `json:"id"`
}

type VirusTotalDomainResponse struct {
	Data []VirusTotalDomain `json:"data"`
	Meta Meta               `json:"meta"`
}

type Meta struct {
	Count int `json:"count"`
}

func ParseVirusTotalDomainResult(reqBody ...string) (virusTotalDomainRespList []VirusTotalDomainResponse) {
	for _, response := range reqBody {
		var virusTotalDomainResponse VirusTotalDomainResponse
		Error.HandleError(json.Unmarshal([]byte(response), &virusTotalDomainResponse))
		virusTotalDomainRespList = append(virusTotalDomainRespList, virusTotalDomainResponse)
	}
	return virusTotalDomainRespList
}
