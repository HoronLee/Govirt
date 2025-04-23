package network

import (
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
