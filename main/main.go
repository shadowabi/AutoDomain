/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"errors"
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
	"github.com/shadowabi/AutoDomain_rebuild/utils/Compare"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	"github.com/shadowabi/AutoDomain_rebuild/utils/File"
	"github.com/shadowabi/AutoDomain_rebuild/utils/log"
	"strings"
)

func init() {
	log.Init("trace")
	configFile := pkg.GetPwd()
	configFile = strings.Join([]string{configFile, "/config.json"}, "")
	err := File.FileNonExistCreate(configFile)
	Error.HandleFatal(err)
	config.SpecificInit(configFile)
	if pkg.IsEmptyConfig(config.C) == true {
		Error.HandleFatal(errors.New("请配置config.json"))
		return
	}
}

func main() {
	cmd.RootCmd.Println(cmd.RootCmd.Long)
	cmd.Execute()

	define.ResultList = Compare.RemoveDuplicates(define.ResultList)
	pkg.WriteToFile(define.ResultList, define.OutPut)
	fmt.Printf("[+] The output is in %s\n", define.OutPut)
}
