package storagePool

import (
	"fmt"
	"govirt/pkg/libvirtd"
	"govirt/pkg/xmlDefine"
	"os"

	"github.com/digitalocean/go-libvirt"
)

// ListAllStoragePools 列出所有存储池
// 该函数返回所有存储池的列表，包括活动和非活动的存储池
func ListAllStoragePools() ([]libvirt.StoragePool, error) {
	pools, _, err := libvirtd.Connection.ConnectListAllStoragePools(1, 0)
	if err != nil {
		return nil, err
	}
	return pools, nil
}

// CreateStoragePool 创建存储池
func CreateStoragePool(params *xmlDefine.PoolTemplateParams) (libvirt.StoragePool, error) {
	xmlDefine.SetDefaults(params)

	// 检查存储池目标路径是否存在，如不存在则创建
	if params.Path != "" {
		if _, err := os.Stat(params.Path); os.IsNotExist(err) {
			if err := os.MkdirAll(params.Path, 0755); err != nil {
				return libvirt.StoragePool{}, fmt.Errorf("创建存储池目标路径失败: %w", err)
			}
		} else if err != nil {
			return libvirt.StoragePool{}, fmt.Errorf("检查存储池目标路径状态失败: %w", err)
		}
	}

	// 渲染XML模板
	xmlStr, err := xmlDefine.RenderTemplate(xmlDefine.PoolTemplate, params)
	if err != nil {
		return libvirt.StoragePool{}, fmt.Errorf("渲染XML失败: %w", err)
	}

	// 定义存储池
	StoragePool, err := libvirtd.Connection.StoragePoolDefineXML(xmlStr, 0)
	if err != nil {
		return libvirt.StoragePool{}, fmt.Errorf("定义存储池失败: %w", err)
	}

	// 设置自启动
	if err := SetStoragePoolAutostart(StoragePool, params.Autostart); err != nil {
		return libvirt.StoragePool{}, fmt.Errorf("设置自启动失败: %w", err)
	}
	return StoragePool, nil
}

// DeleteStoragePool 移除存储池(不删除数据)
func DeleteStoragePool(pool libvirt.StoragePool) error {
	if err := libvirtd.Connection.StoragePoolUndefine(pool); err != nil {
		return fmt.Errorf("移除存储池失败: %w", err)
	}
	return nil
}

// DropStoragePool 删除存储池(删除数据)
func DropStoragePool(pool libvirt.StoragePool) error {
	if err := libvirtd.Connection.StoragePoolDelete(pool, 0); err != nil {
		return fmt.Errorf("删除存储池失败: %w", err)
	}
	return nil
}

// StartStoragePool 启动存储池
func StartStoragePool(pool libvirt.StoragePool) error {
	if err := libvirtd.Connection.StoragePoolCreate(pool, 0); err != nil {
		return fmt.Errorf("启动存储池失败: %w", err)
	}
	return nil
}

// StopStoragePool 停止存储池
func StopStoragePool(pool libvirt.StoragePool) error {
	if err := libvirtd.Connection.StoragePoolDestroy(pool); err != nil {
		return fmt.Errorf("停止存储池失败: %w", err)
	}
	return nil
}

// RefreshStoragePool 刷新存储池
func RefreshStoragePool(pool libvirt.StoragePool) error {
	if err := libvirtd.Connection.StoragePoolRefresh(pool, 0); err != nil {
		return fmt.Errorf("刷新存储池失败: %w", err)
	}
	return nil
}
