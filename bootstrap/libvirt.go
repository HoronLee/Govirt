package bootstrap

import (
	"govirt/pkg/config"
	"govirt/pkg/libvirtd"
	"govirt/pkg/logger"
)

func InitLibvirt() {
	err := libvirtd.InitConnection(config.GetString("libvirt.conURI"))
	if err != nil {
		logger.FatalString("libvirt", "初始化libvirt连接", err.Error())
	}
	// domain.CreateATestDomain()
}
