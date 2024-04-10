package quake

import (
	"encoding/json"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
)

type PortService struct {
	Port    int         `json:"port"`
	Service ServiceInfo `json:"service"`
}

type ServiceInfo struct {
	Name string `json:"name"`

	// HTTP-specific fields, only present when service is 'http' or related
	HTTP struct {
		Host string `json:"host"`
	} `json:"http,omitempty"`

	// Additional fields for other specific services can be added here if needed
}

type QuakeResponse struct {
	Data []PortService `json:"data"`
	Meta Meta          `json:"meta"`
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Total int `json:"total"` // This is the field you asked for
}

func ParseQuakeResult(reqBody ...string) (quakeRespList []QuakeResponse) {
	if len(reqBody) != 0 {
		for _, response := range reqBody {
			var quakeResponse QuakeResponse
			Error.HandleError(json.Unmarshal([]byte(response), &quakeResponse))
			quakeRespList = append(quakeRespList, quakeResponse)
		}
	}

	return quakeRespList
}
