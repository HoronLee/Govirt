package bootstrap

import (
	"fmt"
	"net"

	"govirt/pkg/config"
	"govirt/pkg/logger"
	"govirt/pkg/rpc"

	"google.golang.org/grpc"
)

// SetupGRPCServer 初始化并启动 gRPC 服务器
func SetupGRPCServer() error {
	// 获取配置中的地址和端口
	address := config.GetString("rpc.address")
	port := config.GetString("rpc.port")
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		logger.ErrorString("gRPC", "listen", err.Error())
		return fmt.Errorf("无法启动监听器: %w", err)
	}
	logger.InfoString("gRPC", "listen", "监听地址: "+address+":"+port)

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册服务
	rpc.RegisterService(grpcServer)

	// 启动 gRPC 服务
	logger.InfoString("gRPC", "启动", "gRPC 服务已启动")
	if err := grpcServer.Serve(listener); err != nil {
		logger.ErrorString("gRPC", "serve", err.Error())
		return fmt.Errorf("gRPC 服务启动失败: %w", err)
	}

	return nil
}
