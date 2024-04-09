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
		if config.C.FofaKey == "" || config.C.FofaMail == "" {
			return errors.New("未配置 fofa 相关的凭证")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		define.Once.Do(pkg.GlobalRun)
		fmt.Printf("[+] fofa is working...\n")

		client := pkg.GenerateHTTPClient(define.TimeOut)
		reqString := pkg.MergeReqListToReqString("fofa", define.ReqIpList, define.ReqDomainList)
		reqBody := fofa.FofaRequest(client, reqString, 1)
		reqResult := fofa.ParseFofaResult(reqBody...)

		if int(reqResult[0].Size) > 1000 {
			pageList := net2.GeneratePageList(reqResult[0].Size)
			reqBody2 := fofa.FofaRequest(client, reqString, pageList...)
			reqResult2 := fofa.ParseFofaResult(reqBody2...)
			reqResult = append(reqResult, reqResult2...)
		}

		chanNum := len(reqResult)
		if chanNum != 0 {
			resultChannel := make(chan []string, chanNum)
			resultChannel <- fofa.PurgeDomainResult(reqResult...)
			close(resultChannel)

			pkg.FetchResultFromChanel(resultChannel)
		}

		fmt.Printf("[+] fofa search complete\n")
	},
}
