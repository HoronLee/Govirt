package libvirt

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/digitalocean/go-libvirt"
)

// VirtConnection 自定义连接对象，扩展了libvirt.Libvirt的功能
type VirtConnection struct {
	*libvirt.Libvirt
}

var (
	// connection 全局单例连接
	connection *VirtConnection
	connMutex  sync.Mutex
)

// InitConnection 初始化全局单例连接
func InitConnection(uri string) error {
	connMutex.Lock()
	defer connMutex.Unlock()

	if connection != nil {
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

	connection = &VirtConnection{lv}
	return nil
}
