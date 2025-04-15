package bootstrap

import (
	"govirt/pkg/config"
	"govirt/pkg/libvirt"
	"govirt/pkg/logger"
)

func InitLibvirt() {
	err := libvirt.InitConnection(config.GetString("libvirt.uri"))
	if err != nil {
		logger.FatalString("libvirt", "初始化libvirt连接失败", err.Error())
	}
	// 程序退出时自动关闭连接
	// defer libvirt.CloseConnection()
}
