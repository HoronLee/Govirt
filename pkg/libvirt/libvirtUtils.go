package libvirt

import (
	"fmt"
	"govirt/pkg/logger"
)

// GetConnection 获取全局单例连接
func GetConnection() *VirtConnection {
	if connection == nil {
		panic("libvirt连接未初始化，请先调用InitConnection")
	}
	return connection
}

// CloseConnection 关闭连接
func CloseConnection() error {
	connMutex.Lock()
	defer connMutex.Unlock()

	if connection == nil {
		return nil
	}

	err := connection.Disconnect()
	connection = nil
	return err
}

// 获取 libvirt 版本
func GetLibVersion() (string, error) {
	v, err := GetConnection().ConnectGetLibVersion()
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
