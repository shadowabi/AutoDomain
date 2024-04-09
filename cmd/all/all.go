package all

import (
	"github.com/shadowabi/AutoDomain_rebuild/cmd"
	"github.com/shadowabi/AutoDomain_rebuild/cmd/fofa"
	"github.com/shadowabi/AutoDomain_rebuild/cmd/hunter"
	"github.com/shadowabi/AutoDomain_rebuild/cmd/netlas"
	"github.com/shadowabi/AutoDomain_rebuild/cmd/pulsedive"
	"github.com/shadowabi/AutoDomain_rebuild/cmd/quake"
	"github.com/shadowabi/AutoDomain_rebuild/cmd/virustotal"
	"github.com/shadowabi/AutoDomain_rebuild/cmd/zoomeye"
	"github.com/shadowabi/AutoDomain_rebuild/pkg"
	"github.com/spf13/cobra"
	"sync"
)

func init() {
	cmd.RootCmd.AddCommand(AllCmd)
	AllCmd.AddCommand(fofa.FofaCmd)
	AllCmd.AddCommand(hunter.HunterCmd)
	AllCmd.AddCommand(netlas.NetlasCmd)
	AllCmd.AddCommand(pulsedive.PulsediveCmd)
	AllCmd.AddCommand(quake.QuakeCmd)
	AllCmd.AddCommand(virustotal.VirusTotalCmd)
	AllCmd.AddCommand(zoomeye.ZoomeyeCmd)
}

var AllCmd = &cobra.Command{
	Use:   "all",
	Short: "search domain from all engine",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		pkg.GlobalRun()
		wg.Add(len(cmd.Commands()))
		for _, child := range cmd.Commands() {
			if child.Use != "" {
				err := child.PersistentPreRunE(cmd, args)
				if err != nil {
					wg.Done()
					continue
				}

				go func(child *cobra.Command) {
					child.Run(cmd, args)
					wg.Done()
				}(child)
			}
		}
		wg.Wait()
	},
}
