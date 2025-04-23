package network

import (
	"fmt"
	"govirt/pkg/libvirtd"
	"govirt/pkg/logger"

	"github.com/digitalocean/go-libvirt"
)

// GetNetworkByUUID 根据 UUID 获取网络
func GetNetworkByUUID(uuid libvirt.UUID) (libvirt.Network, error) {
	networks, err := ListAllNetworks()
	if err != nil {
		logger.ErrorString("libvirt", "获取域失败", err.Error())
		return libvirt.Network{}, err
	}
	for _, network := range networks {
		if network.UUID == uuid {
			return network, nil
		}
	}
	return libvirt.Network{}, nil
}

// SetNetworkAutostart 设置网络自动启动
func SetNetworkAutostart(network libvirt.Network, autostart bool) error {
	var autostartFlag int32 = 0
	if autostart {
		autostartFlag = 1
	}

	if err := libvirtd.Connection.NetworkSetAutostart(network, autostartFlag); err != nil {
		return fmt.Errorf("设置网络自动启动失败: %w", err)
	}
	return nil
}

// GetNetworkAutostart 获取网络自动启动状态
func GetNetworkAutostart(network libvirt.Network) (bool, error) {
	autostart, err := libvirtd.Connection.NetworkGetAutostart(network)
	if err != nil {
		return false, fmt.Errorf("获取网络自动启动状态失败: %w", err)
	}

	return autostart == 1, nil
}
