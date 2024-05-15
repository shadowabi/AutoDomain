package quake

import (
	"errors"
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/cmd"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/pkg"
	"github.com/shadowabi/AutoDomain_rebuild/pkg/quake"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(QuakeCmd)
}

var QuakeCmd = &cobra.Command{
	Use:   "quake",
	Short: "search domain from quake",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if config.C.QuakeKey == "" {
			return errors.New("未配置 quake 相关的凭证")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		define.Once.Do(pkg.GlobalRun)
		fmt.Printf("[+] quake is working...\n")

		client := pkg.GenerateHTTPClient(define.TimeOut)

		reqStringList := pkg.MergeReqListToReqStringList("quake", define.ReqIpList, define.ReqDomainList)
		reqBody := quake.QuakeRequest(client, 1, reqStringList...)
		reqResult := quake.ParseQuakeResult(reqBody...)

		for i, _ := range reqResult {
			if int(reqResult[i].Meta.Pagination.Total) > 100 {
				pageList := net2.GeneratePageList(reqResult[i].Meta.Pagination.Total)
				for _, v := range pageList {
					reqBody2 := quake.QuakeRequest(client, v, reqStringList[i])
					reqResult2 := quake.ParseQuakeResult(reqBody2...)
					reqResult = append(reqResult, reqResult2...)
				}
			}
		}

		chanNum := cap(reqResult)
		if chanNum != 0 {
			resultChannel := make(chan []string, chanNum)
			resultChannel <- quake.PurgeDomainResult(reqResult...)
			close(resultChannel)

			pkg.FetchResultFromChanel(resultChannel)
		}

		fmt.Printf("[+] quake search complete\n")
	},
}
