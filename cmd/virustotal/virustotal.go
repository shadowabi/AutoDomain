package virustotal

import (
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/cmd"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/pkg"
	"github.com/shadowabi/AutoDomain_rebuild/pkg/virustotal"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(VirusTotalCmd)
}

var VirusTotalCmd = &cobra.Command{
	Use:   "virustotal",
	Short: "search domain from virustotal",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		define.Once.Do(pkg.GlobalRun)
		fmt.Printf("[+] virustotal is working...\n")

		client := pkg.GenerateHTTPClient(define.TimeOut)

		reqDomainBody := virustotal.VirusTotalDomainRequest(client, define.ReqDomainList...)
		reqDomainResultList := virustotal.ParseVirusTotalDomainResult(reqDomainBody...)

		chanNum := len(reqDomainResultList)
		if chanNum != 0 {
			resultChannel := make(chan []string, chanNum)
			resultChannel <- virustotal.PurgeDomainResult(reqDomainResultList...)
			close(resultChannel)

			pkg.FetchResultFromChanel(resultChannel)
		}

		fmt.Printf("[+] virustotal search complete\n")
	},
}
