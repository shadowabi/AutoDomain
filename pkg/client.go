package pkg

import (
	"crypto/tls"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	"net/http"
	"os"
	"time"
)

func GenerateHTTPClient(timeOut int) *http.Client {
	client := &http.Client{
		Timeout: time.Duration(timeOut) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return client
}

func GetPwd() (homePath string) {
	homePath, err := os.Getwd()
	Error.HandlePanic(err)
	return homePath
}

func IsEmptyConfig(c define.Configure) bool {
	return c.FofaMail == "" &&
		c.FofaKey == "" &&
		c.QuakeKey == "" &&
		c.HunterKey == "" &&
		c.ZoomeyeKey == "" &&
		c.PulsediveKey == ""
}

func GlobalRun() {
	if define.File != "" {
		define.HostList = ParseFileParameter(define.File)
	} else {
		define.HostList = append(define.HostList, define.Url)
	}
	define.ReqIpList = ConvertToReqIpList(define.HostList...)
	define.ReqDomainList = ConvertToReqDomainList(define.HostList...)

	define.TimeOut = (len(define.ReqIpList) + len(define.ReqDomainList)) * define.TimeOut
}