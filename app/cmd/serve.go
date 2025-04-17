package cmd

import (
	"govirt/bootstrap"
	"govirt/pkg/config"
	"govirt/pkg/console"
	"govirt/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// CmdServe represents the available web sub-command.
var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Start web server",
	Run:   runApiServe,
	Args:  cobra.NoArgs,
}

// runApiServe 启动 Web 服务器
func runApiServe(cmd *cobra.Command, args []string) {
	if config.GetBool("rpc.enable") {
		runGrpc(cmd, args)
	} else {
		runWeb(cmd, args)
	}
}

// runWeb 启动 Web 服务器
func runWeb(cmd *cobra.Command, args []string) {

	// 设置 gin 的运行模式，支持 debug, release, test
	// release 会屏蔽调试信息，官方建议生产环境中使用
	// 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// 故此设置为 release，有特殊情况手动改为 debug 即可
	gin.SetMode(gin.ReleaseMode)

	// gin 实例
	router := gin.New()

	// 初始化路由绑定
	bootstrap.SetupRoute(router)

	// 运行服务器
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		logger.ErrorString("CMD", "serve", err.Error())
		console.Exit("Unable to start server, error:" + err.Error())
	}
}

// runGrpc 启动 RPC 服务器
func runGrpc(_ *cobra.Command, _ []string) {
	// Example: Initialize and start an RPC server (replace with actual implementation)
	err := bootstrap.SetupGRPCServer()
	if err != nil {
		logger.ErrorString("CMD", "rpc", err.Error())
		console.Exit("Unable to start RPC server, error: " + err.Error())
	}
}
