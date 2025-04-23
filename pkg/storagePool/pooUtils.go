package storagePool

import (
	"fmt"
	"govirt/pkg/libvirtd"
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

// SetStoragePoolAutostart 设置存储池自动启动
func SetStoragePoolAutostart(pool libvirt.StoragePool, autostart bool) error {
	var autostartFlag int32 = 0
	if autostart {
		autostartFlag = 1
	}

	if err := libvirtd.Connection.StoragePoolSetAutostart(pool, autostartFlag); err != nil {
		return fmt.Errorf("设置存储池自动启动失败: %w", err)
	}
	return nil
}

// GetStoragePoolAutostart 获取存储池自动启动状态
func GetStoragePoolAutostart(pool libvirt.StoragePool) (bool, error) {
	autostart, err := libvirtd.Connection.StoragePoolGetAutostart(pool)
	if err != nil {
		return false, fmt.Errorf("获取存储池自动启动状态失败: %w", err)
	}

	return autostart == 1, nil
}
