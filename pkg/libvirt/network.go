package libvirt

import (
	"github.com/digitalocean/go-libvirt"
)

// ListAllNetworks 列出所有网络
func ListAllNetworks() ([]libvirt.Network, error) {
	networks, _, err := connection.ConnectListAllNetworks(1, libvirt.ConnectListNetworksActive|libvirt.ConnectListNetworksInactive)
	if err != nil {
		return nil, err
	}
	return networks, nil
}
