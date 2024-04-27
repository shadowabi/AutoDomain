package cmd

import (
	"errors"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "AutoDomain",
	Short: "AutoDomain is a web mapping domain name tool",
	Long: "     _         _        ____                        _        \n" +
		"    / \\  _   _| |_ ___ |  _ \\  ___  _ __ ___   __ _(_)_ __   \n" +
		"   / _ \\| | | | __/ _ \\| | | |/ _ \\| '_ ` _ \\ / _` | | '_ \\  \n" +
		"  / ___ \\ |_| | || (_) | |_| | (_) | | | | | | (_| | | | | | \n" +
		" /_/   \\_\\__,_|\\__\\___/|____/ \\___/|_| |_| |_|\\__,_|_|_| |_| \n" +
		` 
        github.com/shadowabi/AutoDomain

AutoDomain是一个集成网络空间测绘系统的工具。
AutoDomain is a tool for integrating cyberspace mapping systems.
`,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if define.Url != "" && define.File != "" {
			Error.HandleFatal(errors.New("参数不可以同时存在"))
			return
		}
		if define.Url == "" && define.File == "" {
			Error.HandleFatal(errors.New("必选参数为空，请输入 -u 参数或 -f 参数"))
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var logLevel string

func init() {
	RootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "info", "设置日志等级 (Set log level) [trace|debug|info|warn|error|fatal|panic]")
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.SetHelpFunc(customHelpFunc)
	RootCmd.PersistentFlags().StringVarP(&define.File, "file", "f", "", "从文件中读取目标地址 (Input FILENAME)")
	RootCmd.PersistentFlags().StringVarP(&define.Url, "url", "u", "", "输入目标地址 (Input [ip|domain|url])")
	RootCmd.PersistentFlags().IntVarP(&define.TimeOut, "timeout", "t", 15, "输入每个 http 请求的超时时间 (Enter the timeout period for every http request)")
	RootCmd.PersistentFlags().StringVarP(&define.OutPut, "output", "o", "./result.txt", "输入结果文件输出的位置 (Enter the location of the scan result output)")
}

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:  RootCmd,
		Headings: cc.HiGreen + cc.Underline,
		Commands: cc.Cyan + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Cyan + cc.Bold,
	})
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func customHelpFunc(cmd *cobra.Command, args []string) {
	cmd.Usage()
}
