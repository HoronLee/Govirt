package storagePool

import (
	"govirt/pkg/logger"

	"github.com/digitalocean/go-libvirt"
)

// GetStoragePoolByUUID 根据 UUID 获取存储池
func GetStoragePoolByUUID(uuid libvirt.UUID) (libvirt.StoragePool, error) {
	pools, err := ListAllStoragePools()
	if err != nil {
		logger.ErrorString("libvirt", "获取存储池失败", err.Error())
		return libvirt.StoragePool{}, err
	}
	for _, domain := range pools {
		if domain.UUID == uuid {
			return domain, nil
		}
	}
	return libvirt.StoragePool{}, nil
}
