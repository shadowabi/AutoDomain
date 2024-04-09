package pulsedive

import (
	"encoding/json"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
)

type Result struct {
	Indicator string `json:"indicator"`
}

type PulsediveDomainResponse struct {
	Results []Result
}

func ParsePulsediveDomainResult(reqBody ...string) (pulsediveDomainResult []PulsediveDomainResponse) {
	for _, response := range reqBody {
		var pulsediveDomainResponse PulsediveDomainResponse
		Error.HandleError(json.Unmarshal([]byte(response), &pulsediveDomainResponse))
		pulsediveDomainResult = append(pulsediveDomainResult, pulsediveDomainResponse)
	}
	return pulsediveDomainResult
}
