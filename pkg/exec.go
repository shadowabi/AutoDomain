package pkg

import (
	"fmt"
	"net/http"
	"os"
	"bufio"
	"sync"
	"strings"
	"regexp"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"crypto/tls"
	"time"
	"strconv"
	"sort"
)

var (
	Drs []string // 存放主域名结果
	Irs []string // 存放IP结果
	Rs2 []string // 存放资产测绘查询结果
	keyword string // 存放关键词
	flag int // 区别IP和域名
	Config Configure //存放配置文件，需传递
	Modes = []string{"fofa", "quake", "hunter", "zoomeye", "vt", "netlas", "pulsedive"} //模块
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4621.0 Safari/537.36"
	client = &http.Client{ //http.client
        Timeout: 15 * time.Second,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }
	mu sync.Mutex
)

type Configure struct {
    FofaMail   	string `json:"FofaMail"`
    FofaKey 	string `json:"FofaKey"`
    QuakeKey   	string `json:"QuakeKey"`
    HunterKey 	string `json:"HunterKey"`
    ZoomeyeKey  string `json:"ZoomeyeKey"`
}

func ReadConfig() {
    data, _ := ioutil.ReadFile("./config/config.json")

    // 解码 JSON 数据
    err := json.Unmarshal(data, &Config)
    if err != nil {
        fmt.Println("config.json配置出错!")
        os.Exit(1)
    }
}

func ReadFile(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("[-] 无法打开此文件")
        os.Exit(1)
    }
    scan := bufio.NewScanner(file)
    var wg sync.WaitGroup
    for scan.Scan() {
        line := strings.TrimSpace(scan.Text())
        wg.Add(1)
        go Match(line, &wg)
    }
    wg.Wait()
    file.Close()
}

func Match(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	ipRegex := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(?:/\d{1,2}|)`)
	ip := ipRegex.FindString(url)
	if ip != "" {
		mu.Lock()
		Irs = append(Irs, ip)
		mu.Unlock()
	}

	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		url = strings.Replace(url, "http://", "", 1)
		url = strings.Replace(url, "https://", "", 1)
		if strings.Contains(url, "/") || strings.Contains(url, "\\") {
			url = url[:strings.IndexAny(url, "/\\")]
		}
	}

	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,11}$`)
	domain := domainRegex.FindString(url)
	if domain != "" {
		subDomainRegex := regexp.MustCompile(`([a-z0-9][a-z0-9\-]*?\.(?:\w{2,4})(?:\.(?:cn|hk))?)$`)
		subDomain := subDomainRegex.FindString(domain)
		if subDomain != "" && !Contains(Drs, subDomain) {
			mu.Lock()
			Drs = append(Drs, subDomain)
			mu.Unlock()
		}
	}
}

func Generate(mode string, wg *sync.WaitGroup) {
	defer wg.Done()
	grammar := ""

	if Contains(Modes, mode) {	
		if mode == "fofa" || mode == "hunter" {
			grammar = "="
		} else if mode == "quake" {
			grammar = ":"
		}

		if grammar != "" {
			if len(Drs) > 0 {
				for _, i := range Drs {
					mu.Lock()
					keyword = keyword + "domain" + grammar + strings.TrimSpace(i) + " || "
					mu.Unlock()
				}
			}

			if len(Irs) > 0 {
				for _, i := range Irs {
					mu.Lock()
					keyword = keyword + "ip" + grammar + strings.TrimSpace(i) + " || "
					mu.Unlock()
				}
			}

			mu.Lock()
			keyword = strings.TrimSuffix(keyword, " || ")
			mu.Unlock()
		}

		if len(Drs) == 0 && len(Irs) == 0 {
			fmt.Println("[-] 无效目标，退出程序！")
			os.Exit(0)
		}

		Scan(mode, keyword)
	} else {
		fmt.Println("[-] 参数错误！")
		os.Exit(0)
	}
}

func Scan(mode string, keyword string) {
	if mode == "fofa" {
		Fofa(keyword)
	}

	if mode == "quake" {
		Quake()
	}

	if mode == "hunter" {
		Hunter(keyword)
	}

	if mode == "zoomeye" {
		Zoomeye()
	}

	if mode == "vt" {
		VirusTotal()
	}

	if mode == "netlas" {
		Netlas()
	}

	if mode == "pulsedive" {
		Pulsedive()
	}

	keyword = "" // 清空变量
}

func ParseBody(resp *http.Response) ([]byte) {

	body, _ := ioutil.ReadAll(resp.Body)
	var data interface{} 
	json.Unmarshal(body, &data)
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	return prettyJSON
}

func DoBatch(requests []*http.Request) ([]*http.Response, error) {
	var wg sync.WaitGroup
	responses := make([]*http.Response, len(requests))

	// 设置WaitGroup的计数器
	wg.Add(len(requests))

	// 并发执行请求
	for i, req := range requests {
		go func(i int, req *http.Request) {
			defer wg.Done()

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			responses[i] = resp
		}(i, req)
	}

	// 等待所有请求完成
	wg.Wait()

	return responses, nil
}

func Contains(slice interface{}, item interface{}) bool {
	switch slice := slice.(type) {
	case []string:
		for _, s := range slice {
			if s == item.(string) {
				return true
			}
		}
	}
	return false
}


func OutPut() {
	if len(Rs2) == 0 {
		fmt.Println("[-] 没有搜集到子域名")
		os.Exit(1)
	}

	file, err := os.OpenFile("./result.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    sort.Strings(Rs2)
    if Rs2 != nil {
        for _, i := range Rs2 {
        	fmt.Println(i)
            fmt.Fprintln(writer, i)
        }
    }

    fmt.Println("已保存到result.txt文件")

}

func Fofa(keyword string) {
	fmt.Println("[+]fofa is working...")
	keyword = base64.URLEncoding.EncodeToString([]byte(keyword))
	url := fmt.Sprintf("https://fofa.info/api/v1/search/all?email=%s&key=%s&qbase64=%s&full=false&fields=protocol,host&size=1000", Config.FofaMail, Config.FofaKey, keyword)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	body := ParseBody(resp)

	var datas map[string]interface{}
	json.Unmarshal(body, &datas)
	if results, ok := datas["results"].([]interface{}); ok {
		for _, result := range results {
			if data, ok := result.([]interface{}); ok {
				if len(data) >= 2 {
					if protocol, ok := data[0].(string); ok {
						if address, ok := data[1].(string); ok {
							var _url string

							if protocol == "http" || protocol == "https" {
								if address[:4] == "http" || address[:5] == "https" {
									_url = address
								} else {
									_url = fmt.Sprintf("%s://%s", protocol, address)
								}
							} else if protocol == "" {
								_url = fmt.Sprintf("http://%s", address)
							}

							if _url != "" && !Contains(Rs2, _url) {
								mu.Lock()
								Rs2 = append(Rs2, _url)
								mu.Unlock()
							}
						}
					}
				}
			}
		}
	}
}

func Quake() {
	fmt.Println("[+]quake is working...")

	payload := strings.NewReader("query=" + keyword + "&start=0&size=100&include=service.name&include=port&include=service.http.host")
	req, err := http.NewRequest("POST", "https://quake.360.net/api/v3/search/quake_service", payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-QuakeToken", Config.QuakeKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body := ParseBody(resp)

	var datas map[string]interface{}

	err = json.Unmarshal(body, &datas)
	if err != nil {
		fmt.Println(err)
		return
	}

	codeStr := fmt.Sprintf("%v", datas["code"])
	code, _ := strconv.Atoi(codeStr)
	if results, ok := datas["data"].([]interface{}); ok && code == 0 {
		for _, result := range results {
			if data, ok := result.(map[string]interface{}); ok {
				port := ""
				if portFloat, ok := data["port"].(float64); ok {
					portInt := int(portFloat)
					if portInt != 80 && portInt != 443 {
						port = ":" + strconv.Itoa(portInt)
					}
				}
				if serviceName, ok := data["service"].(map[string]interface{})["name"].(string); ok {
					if serviceName == "http/ssl" && data["service"].(map[string]interface{})["http"] != nil {
						if host, ok := data["service"].(map[string]interface{})["http"].(map[string]interface{})["host"].(string); ok {
							_url := "https://" + host + port
							if host != "" && port != "" && !Contains(Rs2, _url) {
								mu.Lock()
								Rs2 = append(Rs2, _url)
								mu.Unlock()
							}
						}
					} else if serviceName == "http" && data["service"].(map[string]interface{})["http"] != nil {
						if host, ok := data["service"].(map[string]interface{})["http"].(map[string]interface{})["host"].(string); ok {
							_url := "http://" + host + port
							if host != "" && port != "" && !Contains(Rs2, _url) {
								mu.Lock()
								Rs2 = append(Rs2, _url)
								mu.Unlock()
							}
						}
					}
				}
			}
		}
	}
}

func Hunter(keyword string) {
	fmt.Println("[+]hunter is working...")
	keyword = base64.URLEncoding.EncodeToString([]byte(keyword))
	url := fmt.Sprintf("https://hunter.qianxin.com/openApi/search?api-key=%s&search=%s&page=1&page_size=100&is_web=3", Config.HunterKey, keyword)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body := ParseBody(resp)

	var datas map[string]interface{}
	err = json.Unmarshal(body, &datas)
	if err != nil {
		fmt.Println(err)
		return
	}

	if data, ok := datas["data"].(map[string]interface{}); ok {
		if arr, ok := data["arr"].([]interface{}); ok {
			for _, item := range arr {
				if data, ok := item.(map[string]interface{}); ok {
					if _url, ok := data["url"].(string); ok {
						_url = strings.TrimSpace(_url)
						if _url != "" && !Contains(Rs2, _url) {
							mu.Lock()
							Rs2 = append(Rs2, _url)
							mu.Unlock()
						}
					}
				}
			}
		}
	}
}

func Zoomeye() {
	fmt.Println("[+]zoomeye is working...")
	_url := ""
	grs := []*http.Request{} // 存放异步处理结果
	gIp := []string{}        // 存放需要异步请求的ip
	gDm := []string{}        // 存放需要异步请求的域名
	header1 := http.Header{
		"User-Agent": []string{userAgent},
		"API-KEY":    []string{Config.ZoomeyeKey},
	}

	if len(Irs) > 0 {
		for _, v := range Irs {
			gIp = append(gIp, fmt.Sprintf("https://api.zoomeye.org/host/search?query=ip:%s&facets=port", v))
		}
	}

	if len(Drs) > 0 {
		for _, v := range Drs {
			gDm = append(gDm, fmt.Sprintf("https://api.zoomeye.org/domain/search?q=%s&type=1", v))
		}
	}

	if len(gIp) > 0 {
		for _, v := range gIp {
			req, err := http.NewRequest("GET", v, nil)
			if err != nil {
				continue
			}
			req.Header = header1
			grs = append(grs, req)
		}
	}

	if len(gDm) > 0 {
		for _, v := range gDm {
			req, err := http.NewRequest("GET", v, nil)
			if err != nil {
				fmt.Println(err)
				continue
			}
			req.Header = header1
			grs = append(grs, req)
		}
	}

	responses, err := DoBatch(grs)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, resp := range responses {
		if resp != nil && resp.StatusCode == http.StatusOK {
			body := ParseBody(resp)

			var datas map[string]interface{}
			err = json.Unmarshal(body, &datas)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if total, ok := datas["total"].(float64); ok && total != 0 {
				if matches, ok := datas["matches"].([]interface{}); ok {
					for _, match := range matches {
						if matchData, ok := match.(map[string]interface{}); ok {
							if portInfo, ok := matchData["portinfo"].(map[string]interface{}); ok {
								port := strconv.Itoa(int(portInfo["port"].(float64)))
								if service, ok := portInfo["service"].(string); ok && (service == "http" || service == "https") {
									protocol := service
									_url = protocol + "://" + matchData["ip"].(string) + ":" + port
								} else {
									_url = "http://" + matchData["ip"].(string) + ":" + port
								}
								if _url != "" && !Contains(Rs2, _url) {
									mu.Lock()
									Rs2 = append(Rs2, _url)
									mu.Unlock()
								}
							}
						}
					}
				}

				if lists, ok := datas["list"].([]interface{}); ok {
					for _, list := range lists {
						if listData, ok := list.(map[string]interface{}); ok {
							_url = "http://" + listData["name"].(string)
							if _url != "" && !Contains(Rs2, _url) {
								mu.Lock()
								Rs2 = append(Rs2, _url)
								mu.Unlock()
							}
						}
					}
				}
			}
		}
	}
}

func VirusTotal() {
	fmt.Println("[+]virustotal is working...")
	grs := []*http.Request{}
	gDm := []string{}
	header1 := http.Header{
		"User-Agent": 					[]string{userAgent},
		"X-Vt-Anti-Abuse-Header":    	[]string{"1"},
		"X-Tool":						[]string{"vt-ui-main"},
		"Accept-Ianguage":				[]string{"en-US,en;q=0.9,es;q=0.8"},
	}

	if len(Drs) > 0 {
		for _, v := range Drs {
			gDm = append(gDm, fmt.Sprintf("https://www.virustotal.com/ui/domains/%s/subdomains?relationships=resolutions&cursor=eyJsaW1pdCI6IDIwMCwgIm9mZnNldCI6IDB9&limit=200", v))
		}
	}

	if len(gDm) > 0 {
		for _, v := range gDm {
			req, err := http.NewRequest("GET", v, nil)
			if err != nil {
				continue
			}
			req.Header = header1
			grs = append(grs, req)
		}
	}

	responses, err := DoBatch(grs)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, resp := range responses {
		if resp != nil && resp.StatusCode == http.StatusOK {
			body := ParseBody(resp)
			var datas map[string]interface{}

			err = json.Unmarshal(body, &datas)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if results, ok := datas["data"].([]interface{}); ok {
				for _, result := range results{
					if resultData, ok := result.(map[string]interface{}); ok {
						_url := "http://" + resultData["id"].(string)
						if _url != "" && !Contains(Rs2, _url){
							mu.Lock()
							Rs2 = append(Rs2, _url)
							mu.Unlock()
						}
					}
				}
			}
		}
	}
}

func Netlas() {
	fmt.Println("[+]netlas is working...")

	_url := ""
	grs := []*http.Request{}
	gIp := []string{}
	gDm := []string{}

	if len(Irs) > 0 {
		for _, v := range Irs {
			gIp = append(gIp, fmt.Sprintf("https://app.netlas.io/api/host/%s/?source_type=include&fields=related_domains", v))
		}
	}

	if len(Drs) > 0 {
		for _, v := range Drs {
			gDm = append(gDm, fmt.Sprintf("https://app.netlas.io/api/domains/?q=*.%s&source_type=include&fields=domain", v))
			gDm = append(gDm, fmt.Sprintf("https://app.netlas.io/api/host/%s/?source_type=include&fields=related_domains", v))
		}
	}

	if len(gIp) > 0 {
		for _, v := range gIp {
			req, err := http.NewRequest("GET", v, nil)
			if err != nil {
				continue
			}
			req.Header.Set("User-Agent", userAgent)
			grs = append(grs, req)
		}
	}

	if len(gDm) > 0 {
		for _, v := range gDm {
			req, err := http.NewRequest("GET", v, nil)
			if err != nil {
				fmt.Println(err)
				continue
			}
			req.Header.Set("User-Agent", userAgent)
			grs = append(grs, req)
		}
	}

	responses, err := DoBatch(grs)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, resp := range responses {
		if resp != nil && resp.StatusCode == http.StatusOK {
			body := ParseBody(resp)

			var datas map[string]interface{}
			err = json.Unmarshal(body, &datas)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if relatedDomains, ok := datas["related_domains"].([]interface{}); ok {
				for _, relatedDomain := range relatedDomains {
					if domain, ok := relatedDomain.(string); ok {
						_url = "http://" + domain
						if _url != "" && !Contains(Rs2, _url) {
							mu.Lock()
							Rs2 = append(Rs2, _url)
							mu.Unlock()
						}
					}
				}
			}

			if items, ok := datas["items"].([]interface{}); ok {
				for _, item := range items {
					if itemData, ok := item.(map[string]interface{}); ok {
						if domain, ok := itemData["data"].(map[string]interface{})["domain"].(string); ok {
							_url = "http://" + domain
							if _url != "" && !Contains(Rs2, _url) {
								mu.Lock()
								Rs2 = append(Rs2, _url)
								mu.Unlock()
							}
						}
					}
				}
			}
		}
	}
}

func Pulsedive() {
	fmt.Println("[+]pulsedive is working...")
	_url := ""
	grs := []*http.Request{} // 存放异步处理结果
	gDm := []string{}        // 存放需要异步请求的域名

	for _, v := range Drs {
		gDm = append(gDm, fmt.Sprintf("https://pulsedive.com/api/explore.php?q=ioc%%3d*.%s%%20active%%3dtrue", v))
	}

	if len(gDm) > 0 {
		for _, v := range gDm {
			req, err := http.NewRequest("GET", v, nil)
			if err != nil {
				fmt.Println(err)
				continue
			}
			req.Header.Set("User-Agent", userAgent)

			grs = append(grs, req)
		}
	}

	responses, err := DoBatch(grs)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, resp := range responses {
		if resp != nil && resp.StatusCode == http.StatusOK {
			body := ParseBody(resp)
			var datas map[string]interface{}
			err = json.Unmarshal(body, &datas)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if results, ok := datas["results"].([]interface{}); ok {
				for _, result := range results {
					if resultData, ok := result.(map[string]interface{}); ok {
						if indicator, ok := resultData["indicator"].(string); ok {
							_url = "http://" + indicator
							if _url != "" && !Contains(Rs2, _url) {
								mu.Lock()
								Rs2 = append(Rs2, _url)
								mu.Unlock()
							}
						}
					}
				}
			}
		}
	}
}