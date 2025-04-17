package rpc

import (
	libvirtPb "govirt/app/proto/libvirt" // 导入生成的 protobuf 代码
	"govirt/app/services"

	"google.golang.org/grpc"
)

// RegisterService 注册所有 gRPC 服务
func RegisterService(grpcServer *grpc.Server) {
	libvirtPb.RegisterLibvirtServiceServer(grpcServer, &services.LibvirtService{})
	// 如果有其他服务，可以在这里继续注册
	// e.g., anotherProto.RegisterAnotherServiceServer(grpcServer, &anotherService{})
}
