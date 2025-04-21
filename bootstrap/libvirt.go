package bootstrap

import (
	"govirt/pkg/config"
	"govirt/pkg/libvirt"
	"govirt/pkg/logger"
)

func InitLibvirt() {
	err := libvirt.InitConnection(config.GetString("libvirt.conURI"))
	if err != nil {
		logger.FatalString("libvirt", "初始化libvirt连接", err.Error())
	}
	// libvirt.CreateATestDomain()
}
