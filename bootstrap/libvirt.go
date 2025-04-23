package bootstrap

import (
	"govirt/pkg/config"
	"govirt/pkg/libvirtd"
	"govirt/pkg/logger"
	"govirt/pkg/storagePool"
)

func InitLibvirt() {
	err := libvirtd.InitConnection(config.GetString("libvirt.conURI"))
	if err != nil {
		logger.FatalString("libvirt", "初始化libvirt连接", err.Error())
	}
	// domain.CreateATestDomain()
	err = storagePool.InitSystemStoragePool(config.GetString("libvirt.pool.image.name"), config.GetString("libvirt.pool.image.path"))
	if err != nil {
		logger.FatalString("libvirt", "初始化镜像存储池", err.Error())
	}
	err = storagePool.InitSystemStoragePool(config.GetString("libvirt.pool.volume.name"), config.GetString("libvirt.pool.volume.path"))
	if err != nil {
		logger.FatalString("libvirt", "初始化数据存储池", err.Error())
	}
	logger.InfoString("libvirt", "初始化Libvirt控制器", "成功")
}
