package cmd

import (
	"os"
	"github.com/spf13/cobra"
	exec "github.com/shadowabi/AutoDomain/pkg"
	"strings"
	"sync"
)

var (
	file	  string
	url		  string
	mode 	  string
)


var RootCmd = &cobra.Command{
	Use:   "Serverless_PortScan",
	Short: "Serverless_PortScan is used to scan ports using cloud functions.",
	Long: 
"     _         _        ____                        _        \n" +
"    / \\  _   _| |_ ___ |  _ \\  ___  _ __ ___   __ _(_)_ __   \n" +
"   / _ \\| | | | __/ _ \\| | | |/ _ \\| '_ ` _ \\ / _` | | '_ \\  \n" +
"  / ___ \\ |_| | || (_) | |_| | (_) | | | | | | (_| | | | | | \n" +
" /_/   \\_\\__,_|\\__\\___/|____/ \\___/|_| |_| |_|\\__,_|_|_| |_| \n" +
                                                            
                                                            
                                                                              
` 
        github.com/shadowabi/AutoDomain

AutoDomain是一个集成网络空间测绘系统的工具。
AutoDomain is a tool for integrating cyberspace mapping systems.
`,
}



func init() {
	exec.ReadConfig()
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.Flags().StringVarP(&file, "file", "f", "", "从文件中读取目标地址 (Input FILENAME)")
	RootCmd.Flags().StringVarP(&url, "url", "u", "", "输入目标地址 (Input IP/DOMAIN/URL)")
	RootCmd.Flags().StringVarP(&mode, "mode", "m", "all" , "可选择特定的测绘模块，例如fofa、quake、hunter、vt、netlas、pulsedive，默认all为全选 (Specific mapping modules can be selected, such as fofa, quake, hunter, vt, netlas, pulsedive, and all is selected by default)")
}


func Execute(){
	err := RootCmd.Execute()

	var wg sync.WaitGroup
    if url != "" {
    	wg.Add(1)
        go exec.Match(strings.TrimSpace(url), &wg)
        wg.Wait()
    } else if file != "" {
        exec.ReadFile(file)
    } else {
    	RootCmd.Usage()
    }

    if mode == "all" {
		for _, m := range exec.Modes {
			wg.Add(1)
			go exec.Generate(m, &wg)
		}
	} else {
		wg.Add(1)
		go exec.Generate(mode,&wg)
	}
	wg.Wait()
	
	exec.OutPut()

	if err != nil {
		os.Exit(1)
	}
}