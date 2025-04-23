package storagePool

import (
	"fmt"
	"govirt/pkg/libvirtd"
	"govirt/pkg/logger"
	"govirt/pkg/xmlDefine"

	"github.com/digitalocean/go-libvirt"
)

// InitSystemStoragePool 初始化多个存储池
func InitSystemStoragePool(params ...xmlDefine.PoolTemplateParams) error {
	for _, param := range params {
		// 获取存储池
		pool, err := GetStoragePoolByName(param.Name)
		if err != nil {
			logger.ErrorString("libvirt", "初始化存储池", err.Error())
			return err
		}

		// 如果存储池不存在，则创建
		if pool.Name == "" {
			newPool, err := CreateStoragePool(&param)
			if err != nil {
				return fmt.Errorf("创建存储池 %s 失败: %w", param.Name, err)
			}
			if err := StartStoragePool(newPool); err != nil {
				return fmt.Errorf("启动存储池 %s 失败: %w", param.Name, err)
			}
			logger.InfoString("libvirt", "初始化存储池", fmt.Sprintf("成功创建并启动存储池 %s", param.Name))
		} else {
			logger.WarnString("libvirt", "初始化存储池", fmt.Sprintf("存储池 %s 已存在", param.Name))
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
