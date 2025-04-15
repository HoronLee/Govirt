package libvirt

import (
	"govirt/pkg/helpers"
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

// FormatNetworks 格式化网络信息
func FormatNetworks(domains []libvirt.Network) []map[string]any {
	var formattedNetworks []map[string]any
	for _, d := range domains {
		// state, _ := GetDomainStateByUUID(d.UUID) // 使用下划线忽略错误
		formattedNetworks = append(formattedNetworks, map[string]any{
			"Name": d.Name,
			"UUID": helpers.UUIDBytesToString(d.UUID),
			// "State": DomainStateToString(state),
		})
	}
	return formattedNetworks
}
