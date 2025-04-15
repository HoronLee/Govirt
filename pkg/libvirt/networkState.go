package libvirt

import (
	"govirt/pkg/logger"

	"github.com/digitalocean/go-libvirt"
)

// ListAllNetworks 列出所有网络
func ListAllNetworks() ([]libvirt.Network, error) {
	networks, _, err := GetConnection().ConnectListAllNetworks(1, libvirt.ConnectListNetworksActive|libvirt.ConnectListNetworksInactive)
	if err != nil {
		logger.ErrorString("libvirt", "列出所有网络失败", err.Error())
		return nil, err
	}
	return networks, nil
}
