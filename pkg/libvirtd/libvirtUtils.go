package libvirtd

import (
	"fmt"
	"govirt/pkg/config"
	"govirt/pkg/logger"
	"os"
	"strconv"
	"strings"
)

// GetConnection 获取全局单例连接
func GetConnection() *VirtConn {
	if Conn == nil {
		panic("libvirt连接未初始化，请先调用InitConnection")
	}
	return Conn
}

// CloseConnection 关闭连接
func CloseConnection() error {
	connMutex.Lock()
	defer connMutex.Unlock()

	if Conn == nil {
		return nil
	}

	err := Conn.Disconnect()
	Conn = nil
	return err
}

// ToStandardVersion 将数值转换为标准版本格式
func ToStandardVersion(v uint64) string {
	major := v / 1000000
	minor := (v % 1000000) / 1000
	release := v % 1000
	return fmt.Sprintf("%d.%d.%d", major, minor, release)
}

// ServerInfo 宿主机信息结构体
type ServerInfo struct {
	HOST_NAME   string
	HOST_UUID   string
	CPU_CORE    string
	MEMORY      string
	LIB_VERSION string
	LIB_URI     string
}

// GetServerInfo 获取宿主机信息
func GetServerInfo() (*ServerInfo, error) {
	info := &ServerInfo{}

	info.HOST_NAME = config.Get("libvirt.hostName")
	info.HOST_UUID = config.Get("libvirt.hostUUID")

	// 获取版本
	version, err := Conn.ConnectGetLibVersion()
	if err != nil {
		logger.FatalString("libvirt", "获取libvirt版本失败", err.Error())
		return nil, err
	}
	info.LIB_VERSION = ToStandardVersion(version)

	info.LIB_URI = config.GetString("libvirt.conURI")

	// 获取CPU核心数
	cpuInfo, err := os.ReadFile("/proc/cpuinfo")
	if err == nil {
		cpuCount := strings.Count(string(cpuInfo), "processor")
		info.CPU_CORE = fmt.Sprintf("%d cores", cpuCount)
	}

	// 获取内存信息
	memInfo, err := os.ReadFile("/proc/meminfo")
	if err == nil {
		lines := strings.Split(string(memInfo), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					memKB, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						info.MEMORY = fmt.Sprintf("%.2f GB", float64(memKB)/1024/1024)
					}
					break
				}
			}
		}
	}

	return info, nil
}
