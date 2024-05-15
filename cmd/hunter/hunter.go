package hunter

import (
	"errors"
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/cmd"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/pkg"
	"github.com/shadowabi/AutoDomain_rebuild/pkg/hunter"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(HunterCmd)
}

var HunterCmd = &cobra.Command{
	Use:   "hunter",
	Short: "search domain from hunter",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if config.C.HunterKey == "" {
			return errors.New("未配置 hunter 相关的凭证")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		define.Once.Do(pkg.GlobalRun)
		fmt.Printf("[+] hunter is working...\n")

		client := pkg.GenerateHTTPClient(define.TimeOut)

		reqStringList := pkg.MergeReqListToReqStringList("hunter", define.ReqIpList, define.ReqDomainList)
		reqBody := hunter.HunterRequest(client, 1, reqStringList...)
		reqResult := hunter.ParseHunterResult(reqBody...)

		for i, _ := range reqResult {
			if int(reqResult[i].Data.Total) > 1000 {
				pageList := net2.GeneratePageList(reqResult[i].Data.Total)
				for _, v := range pageList {
					reqBody2 := hunter.HunterRequest(client, v, reqStringList[i])
					reqResult2 := hunter.ParseHunterResult(reqBody2...)
					reqResult = append(reqResult, reqResult2...)
				}
			}
		}

		chanNum := cap(reqResult)
		if chanNum != 0 {
			resultChannel := make(chan []string, chanNum)
			resultChannel <- hunter.PurgeDomainResult(reqResult...)
			close(resultChannel)

			pkg.FetchResultFromChanel(resultChannel)
		}

		fmt.Printf("[+] hunter search complete\n")
	},
}
