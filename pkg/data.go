package pkg

import (
	"bufio"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Compare"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	"os"
	"regexp"
	"strings"
	"sync"
)

func ParseFileParameter(fileName string) (fileHostList []string) {
	file, err := os.Open(fileName)
	Error.HandlePanic(err)
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		fileHostList = append(fileHostList, line)
	}
	file.Close()
	return fileHostList
}

func TripProtocolString(param string) string {
	if strings.Contains(param, "http://") || strings.Contains(param, "https://") {
		param = strings.Replace(param, "http://", "", 1)
		param = strings.Replace(param, "https://", "", 1)
		if strings.Contains(param, "/") || strings.Contains(param, "\\") {
			param = param[:strings.IndexAny(param, "/\\")]
		}
	}
	return param
}

func ConvertToReqIpList(param ...string) (reqIpList []string) {
	if len(param) != 0 {
		for _, host := range param {
			host := TripProtocolString(host)
			ipRegex := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(?:/\d{1,2}|)`)
			ip := ipRegex.FindString(host)
			if ip != "" && !Compare.IsStringInStringArray(ip, reqIpList) {
				reqIpList = append(reqIpList, ip)
			}
		}
	}
	return reqIpList
}

func ConvertToReqDomainList(param ...string) (reqDomainList []string) {
	if len(param) != 0 {
		for _, host := range param {
			host := TripProtocolString(host)
			domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,11}$`)
			domain := domainRegex.FindString(host)
			if domain != "" {
				subDomainRegex := regexp.MustCompile(`([a-z0-9][a-z0-9\-]*?\.(?:\w{2,4})(?:\.(?:cn|hk))?)$`)
				subDomain := subDomainRegex.FindString(domain)
				if subDomain != "" && !Compare.IsStringInStringArray(subDomain, reqDomainList) {
					reqDomainList = append(reqDomainList, subDomain)
				}
			}
		}
	}
	return reqDomainList
}

func MergeReqListToReqString(mode string, reqIpList []string, reqDomainList []string) (reqString string) {
	grammar := define.ModeToGrammar[mode]
	if grammar != "" {
		for _, ip := range reqIpList {
			reqString += "ip" + grammar + ip + " || "
		}
		for _, ip := range reqDomainList {
			reqString += "domain" + grammar + ip + " || "
		}
	}
	reqString = strings.TrimSuffix(reqString, " || ")
	return reqString
}

func FetchResultFromChanel(resultChannel chan []string) {
	var mu sync.Mutex
	for v := range resultChannel {
		mu.Lock()
		define.ResultList = append(define.ResultList, v...)
		mu.Unlock()
	}
}