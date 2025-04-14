package libvirt

import (
	"fmt"
	"govirt/pkg/logger"

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
	version := ToStandardVersion(v)
	return version, nil
}

// ToStandardVersion 将数值转换为标准版本格式
func ToStandardVersion(v uint64) string {
	major := v / 1000000
	minor := (v % 1000000) / 1000
	release := v % 1000
	return fmt.Sprintf("%d.%d.%d", major, minor, release)
}
