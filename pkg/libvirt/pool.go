package libvirt

import (
	"github.com/digitalocean/go-libvirt"
)

func ListAllStoragePools() ([]libvirt.StoragePool, error) {
	pools, _, err := connection.ConnectListAllStoragePools(1, libvirt.ConnectListStoragePoolsActive|libvirt.ConnectListStoragePoolsInactive)
	if err != nil {
		return nil, err
	}
	return pools, nil
}
