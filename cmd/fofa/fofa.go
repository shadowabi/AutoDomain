package fofa

import (
	"errors"
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/cmd"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/pkg"
	"github.com/shadowabi/AutoDomain_rebuild/pkg/fofa"
	net2 "github.com/shadowabi/AutoDomain_rebuild/utils/response"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(FofaCmd)
}

var FofaCmd = &cobra.Command{
	Use:   "fofa",
	Short: "search domain from fofa",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if config.C.FofaKey == "" {
			return errors.New("未配置 fofa 相关的凭证")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		define.Once.Do(pkg.GlobalRun)
		fmt.Printf("[+] fofa is working...\n")

		client := pkg.GenerateHTTPClient(define.TimeOut)
		reqStringList := pkg.MergeReqListToReqStringList("fofa", define.ReqIpList, define.ReqDomainList)
		reqBody := fofa.FofaRequest(client, 1, reqStringList...)
		reqResult := fofa.ParseFofaResult(reqBody...)

		for i, _ := range reqResult {
			if int(reqResult[i].Size) > 1000 {
				pageList := net2.GeneratePageList(reqResult[i].Size)
				for _, v := range pageList {
					reqBody2 := fofa.FofaRequest(client, v, reqStringList[i])
					reqResult2 := fofa.ParseFofaResult(reqBody2...)
					reqResult = append(reqResult, reqResult2...)
				}
			}
		}

		chanNum := cap(reqResult)
		if chanNum != 0 {
			resultChannel := make(chan []string, chanNum)
			resultChannel <- fofa.PurgeDomainResult(reqResult...)
			close(resultChannel)

			pkg.FetchResultFromChanel(resultChannel)
		}

		fmt.Printf("[+] fofa search complete\n")
	},
}
