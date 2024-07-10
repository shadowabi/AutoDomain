/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/cmd"
	_ "github.com/shadowabi/AutoDomain_rebuild/cmd/all"
	_ "github.com/shadowabi/AutoDomain_rebuild/cmd/daydaymap"
	_ "github.com/shadowabi/AutoDomain_rebuild/cmd/fofa"
	_ "github.com/shadowabi/AutoDomain_rebuild/cmd/hunter"
	_ "github.com/shadowabi/AutoDomain_rebuild/cmd/netlas"
	_ "github.com/shadowabi/AutoDomain_rebuild/cmd/pulsedive"
	_ "github.com/shadowabi/AutoDomain_rebuild/cmd/quake"
	_ "github.com/shadowabi/AutoDomain_rebuild/cmd/virustotal"
	_ "github.com/shadowabi/AutoDomain_rebuild/cmd/zoomeye"
	"github.com/shadowabi/AutoDomain_rebuild/config"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/pkg"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
)

func init() {
	config.InitConfigure("config.yaml")
}

func main() {
	if pkg.IsEmptyConfig(config.C) {
		Error.HandleFatal(fmt.Errorf("请配置 config.yaml"))
		return
	}

	fmt.Println(cmd.RootCmd.Long)
	cmd.Execute()

	fmt.Printf("[+] The output is in %s\n", define.OutPut)
}
