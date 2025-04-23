package storagePool

import (
	"govirt/pkg/logger"
	"govirt/pkg/xmlDefine"

	"github.com/digitalocean/go-libvirt"
)

// InitSystemStoragePool 初始化默认存储池
func InitSystemStoragePool(name string, path string) error {
	// 获取默认存储池
	pool, err := GetStoragePoolByName(name)
	if err != nil {
		logger.ErrorString("libvirt", "初始化默认存储池", err.Error())
		return err
	}
	if pool.Name == "" {
		params := xmlDefine.PoolTemplateParams{
			Name: name,
			Path: path,
		}
		pool, err := CreateStoragePool(&params)
		if err != nil {
			return err
		}
		if err := StartStoragePool(pool); err != nil {
			return err
		}
	}
	return nil
}

// GetStoragePoolByUUID 根据 UUID 获取存储池
func GetStoragePoolByUUID(uuid libvirt.UUID) (libvirt.StoragePool, error) {
	pools, err := ListAllStoragePools()
	if err != nil {
		logger.ErrorString("libvirt", "获取存储池", err.Error())
		return libvirt.StoragePool{}, err
	}
	for _, domain := range pools {
		if domain.UUID == uuid {
			return domain, nil
		}
	}
	return libvirt.StoragePool{}, nil
}

// GetStoragePoolByName 根据名称获取存储池
func GetStoragePoolByName(name string) (libvirt.StoragePool, error) {
	pools, err := ListAllStoragePools()
	if err != nil {
		logger.ErrorString("libvirt", "获取存储池失败", err.Error())
		return libvirt.StoragePool{}, err
	}
	for _, domain := range pools {
		if domain.Name == name {
			return domain, nil
		}
	}
	return libvirt.StoragePool{}, nil
}
