package daydaymap

import (
	"errors"
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/cmd"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/pkg"
	"github.com/shadowabi/AutoDomain_rebuild/pkg/daydaymap"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(DaydayMapCmd)
}

var DaydayMapCmd = &cobra.Command{
	Use:   "daydaymap",
	Short: "search domain from daydaymap",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if config.C.DaydaymapKey == "" {
			return errors.New("未配置 daydaymap 相关的凭证")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		define.Once.Do(pkg.GlobalRun)
		fmt.Printf("[+] daydaymap is working...\n")

		client := pkg.GenerateHTTPClient(define.TimeOut)
		reqStringList := pkg.MergeReqListToReqStringList("daydaymap", define.ReqIpList, define.ReqDomainList)
		reqBody := daydaymap.DayDayMapRequest(client, 1, 1, reqStringList...)
		reqResult := daydaymap.ParseDaydaymapResult(reqBody...)

		for i, _ := range reqResult {
			if int(reqResult[i].Data.Total) == 0 {
				continue
			}
			if int(reqResult[i].Data.Total) > 10000 {
				for j := 1; i <= reqResult[i].Data.Total/10000; j++ {
					reqBody = daydaymap.DayDayMapRequest(client, j, 10000, reqStringList[i])
					reqResult = append(reqResult, daydaymap.ParseDaydaymapResult(reqBody...)...)
				}
			} else {
				reqBody = daydaymap.DayDayMapRequest(client, 1, reqResult[i].Data.Total, reqStringList[i])
				reqResult = append(reqResult, daydaymap.ParseDaydaymapResult(reqBody...)...)
			}
		}

		chanNum := cap(reqResult)
		if chanNum != 0 {
			resultChannel := make(chan []string, chanNum)
			resultChannel <- daydaymap.PurgeDomainResult(reqResult...)
			close(resultChannel)

			pkg.FetchResultFromChanel(resultChannel)
		}
		fmt.Printf("[+] daydaymap search complete\n")
	},
}
