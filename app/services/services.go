package services

import (
	"context"
	"errors"
	libvirtPb "govirt/app/proto/libvirt"
	"govirt/pkg/libvirt"
)

// LibvirtService 定义 RPC 服务
type LibvirtService struct {
	libvirtPb.UnimplementedLibvirtServiceServer // 嵌入未实现的默认方法
}

// GetLibVersion 实现 libvirt.LibvirtServiceServer 接口
func (s *LibvirtService) GetLibVersion(ctx context.Context, req *libvirtPb.GetLibVersionRequest) (*libvirtPb.GetLibVersionResponse, error) {
	version, err := libvirt.GetLibVersion()
	if err != nil {
		return nil, errors.New("获取 libvirt 版本失败: " + err.Error())
	}
	return &libvirtPb.GetLibVersionResponse{Version: version}, nil
}
