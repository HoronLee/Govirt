package libvirt

import (
	"gohub/pkg/helpers"
	"gohub/pkg/logger"

	"github.com/digitalocean/go-libvirt"
)

var (
	// Libvirt 全局 libvirt 连接对象
	Libvirt *libvirt.Libvirt
)

// 获取 libvirt 版本
func GetLibVersion() (string, error) {
	v, err := Libvirt.ConnectGetLibVersion()
	if err != nil {
		logger.FatalString("libvirt", "获取libvirt版本失败", err.Error())
		return "", err
	}
	version := helpers.ToStandardVersion(v)
	return version, nil
}
