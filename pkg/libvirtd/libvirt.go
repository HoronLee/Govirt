package libvirtd

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/digitalocean/go-libvirt"
)

// VirtConn 自定义连接对象，扩展了libvirt.Libvirt的功能
type VirtConn struct {
	*libvirt.Libvirt
}

var (
	// connection 全局单例连接
	Conn      *VirtConn
	connMutex sync.Mutex
)

// InitConnection 初始化全局单例连接
func InitConnection(uri string) error {
	connMutex.Lock()
	defer connMutex.Unlock()

	if Conn != nil {
		return nil
	}

	parsedURI, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("解析URI失败: %v", err)
	}

	lv, err := libvirt.ConnectToURI(parsedURI)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}

	Conn = &VirtConn{lv}
	return nil
}
