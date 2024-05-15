package netlas

import (
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/cmd"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/pkg"
	"github.com/shadowabi/AutoDomain_rebuild/pkg/netlas"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(NetlasCmd)
}

var NetlasCmd = &cobra.Command{
	Use:   "netlas",
	Short: "search domain from netlas",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		define.Once.Do(pkg.GlobalRun)
		fmt.Printf("[+] netlas is working...\n")

		client := pkg.GenerateHTTPClient(define.TimeOut)

		reqDomainBody := netlas.NetlasDomainRequest(client, define.ReqDomainList...)
		reqDomainResultList := netlas.ParseNetlasDomainResult(reqDomainBody...)

		reqIpBody := netlas.NetlasIpRequest(client, define.ReqIpList...)
		reqIpResultList := netlas.ParseNetlasIpResult(reqIpBody...)

		chanNum := cap(reqDomainResultList) + cap(reqIpResultList)
		if chanNum != 0 {
			resultChannel := make(chan []string, chanNum)
			if len(reqDomainResultList) != 0 {
				resultChannel <- netlas.PurgeDomainResult(reqDomainResultList...)
			}

			if len(reqIpResultList) != 0 {
				resultChannel <- netlas.PurgeIpResult(reqIpResultList...)
			}

			close(resultChannel)

			pkg.FetchResultFromChanel(resultChannel)
		}

		fmt.Printf("[+] netlas search complete\n")
	},
}
