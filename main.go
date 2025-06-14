package main

import (
	"fmt"
	"govirt/app/cmd"
	"govirt/app/cmd/make"
	"govirt/bootstrap"
	btsConfig "govirt/config"
	"govirt/pkg/config"
	"govirt/pkg/console"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// 加载 config 目录下的配置信息
	btsConfig.Initialize()
}

func main() {

	// 应用的主入口，默认调用 cmd.CmdServe 命令
	var rootCmd = &cobra.Command{
		Use:   "govirt",
		Short: "A libvirt client for managing virtual machines",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			// 配置初始化，依赖命令行 --env 参数
			config.InitConfig(cmd.Env)
			// 初始化 Logger
			bootstrap.SetupLogger()
			// 初始化数据库
			bootstrap.SetupDB()
			// 初始化 apikey
			bootstrap.InitApikey()
			// 初始化libvirt
			bootstrap.InitLibvirt()
		},
	}

	// 注册子命令
	rootCmd.AddCommand(
		cmd.CmdServe,
		make.CmdMake,
	)

	// 配置默认运行 Web 服务
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 注册全局参数，--env
	cmd.RegisterGlobalFlags(rootCmd)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
