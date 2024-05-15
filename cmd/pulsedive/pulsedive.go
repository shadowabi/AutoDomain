package pulsedive

import (
	"errors"
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/cmd"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/pkg"
	"github.com/shadowabi/AutoDomain_rebuild/pkg/pulsedive"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(PulsediveCmd)
}

var PulsediveCmd = &cobra.Command{
	Use:   "pulsedive",
	Short: "search domain from pulsedive",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if config.C.PulsediveKey == "" {
			return errors.New("未配置 pulsedive 相关的凭证")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		pkg.GlobalRun()
		fmt.Printf("[+] pulsedive is working...\n")

		client := pkg.GenerateHTTPClient(define.TimeOut)

		reqDomainBody := pulsedive.PulsediveDomainRequest(client, define.ReqDomainList...)
		reqDomainResultList := pulsedive.ParsePulsediveDomainResult(reqDomainBody...)

		chanNum := cap(reqDomainResultList)
		if chanNum != 0 {
			resultChannel := make(chan []string, chanNum)
			resultChannel <- pulsedive.PurgeDomainResult(reqDomainResultList...)
			close(resultChannel)

			pkg.FetchResultFromChanel(resultChannel)
		}

		fmt.Printf("[+] pulsedive search complete\n")
	},
}
