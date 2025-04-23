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
	if err := InitStoragePool(); err != nil {
		logger.FatalString("libvirt", "初始化存储池", err.Error())
	}
	// domain.CreateATestDomain()
	logger.InfoString("libvirt", "初始化Libvirt控制器", "成功")
}

// InitStoragePool 初始化存储池
func InitStoragePool() error {
	pools := []struct {
		nameKey string
		pathKey string
	}{
		{"libvirt.pool.image.name", "libvirt.pool.image.path"},
		{"libvirt.pool.volume.name", "libvirt.pool.volume.path"},
	}

	for _, pool := range pools {
		if err := storagePool.InitSystemStoragePool(config.GetString(pool.nameKey), config.GetString(pool.pathKey)); err != nil {
			return err
		}
	}
	return nil
}

// TODO: 初始化网络