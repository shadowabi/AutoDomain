package zoomeye

import (
	"errors"
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/cmd"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/pkg"
	"github.com/shadowabi/AutoDomain_rebuild/pkg/zoomeye"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(ZoomeyeCmd)
}

var ZoomeyeCmd = &cobra.Command{
	Use:   "zoomeye",
	Short: "search domain from zoomeye",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if config.C.ZoomeyeKey == "" {
			return errors.New("未配置 zoomeye 相关的凭证")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		define.Once.Do(pkg.GlobalRun)
		fmt.Printf("[+] zoomeye is working...\n")

		client := pkg.GenerateHTTPClient(define.TimeOut)

		reqIpBody := zoomeye.ZoomeyeIpRequest(client, define.ReqIpList...)
		reqIpResultList := zoomeye.ParseZoomeyeIpResult(reqIpBody...)

		reqDomainBody := zoomeye.ZoomeyeDomainRequest(client, define.ReqDomainList...)
		reqDomainResultList := zoomeye.ParseZoomeyeDomainResult(reqDomainBody...)

		chanNum := len(reqDomainResultList) + len(reqIpResultList)
		if chanNum != 0 {
			resultChannel := make(chan []string, chanNum)
			if len(reqDomainResultList) != 0 {
				resultChannel <- zoomeye.PurgeDomainResult(reqDomainResultList...)
			}

			if len(reqIpResultList) != 0 {
				resultChannel <- zoomeye.PurgeIpResult(reqIpResultList...)
			}

			close(resultChannel)

			pkg.FetchResultFromChanel(resultChannel)
		}

		fmt.Printf("[+] zoomeye search complete\n")
	},
}
