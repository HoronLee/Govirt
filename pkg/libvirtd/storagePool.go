// filepath: /home/horonlee/projects/govirt/pkg/libvirtd/storagePool.go
package libvirtd

import (
	"fmt"
	"govirt/pkg/helpers"
	"govirt/pkg/logger"
	"govirt/pkg/xmlDefine"
	"os"

	"github.com/digitalocean/go-libvirt"
)

// ListAllStoragePools 列出所有存储池
// 该函数返回所有存储池的列表，包括活动和非活动的存储池
func (vc *VirtConn) ListAllStoragePools() ([]libvirt.StoragePool, error) {
	pools, _, err := vc.Libvirt.ConnectListAllStoragePools(1, 0)
	if err != nil {
		return nil, err
	}
	return pools, nil
}

// CreateStoragePool 创建存储池
func (vc *VirtConn) CreateStoragePool(params *xmlDefine.PoolTemplateParams) (libvirt.StoragePool, error) {
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
	StoragePool, err := vc.Libvirt.StoragePoolDefineXML(xmlStr, 0)
	if err != nil {
		return libvirt.StoragePool{}, fmt.Errorf("定义存储池失败: %w", err)
	}

	// 设置自启动
	if err := vc.SetStoragePoolAutostart(StoragePool, params.Autostart); err != nil {
		return libvirt.StoragePool{}, fmt.Errorf("设置自启动失败: %w", err)
	}
	return StoragePool, nil
}

// DeleteStoragePool 移除存储池(不删除数据)
func (vc *VirtConn) DeleteStoragePool(pool libvirt.StoragePool) error {
	if err := vc.Libvirt.StoragePoolUndefine(pool); err != nil {
		return fmt.Errorf("移除存储池失败: %w", err)
	}
	return nil
}

// DropStoragePool 删除存储池(删除数据)
func (vc *VirtConn) DropStoragePool(pool libvirt.StoragePool) error {
	if err := vc.Libvirt.StoragePoolDelete(pool, 0); err != nil {
		return fmt.Errorf("删除存储池失败: %w", err)
	}
	return nil
}

// StartStoragePool 启动存储池
func (vc *VirtConn) StartStoragePool(pool libvirt.StoragePool) error {
	if err := vc.Libvirt.StoragePoolCreate(pool, 0); err != nil {
		return fmt.Errorf("启动存储池失败: %w", err)
	}
	return nil
}

// StopStoragePool 停止存储池
func (vc *VirtConn) StopStoragePool(pool libvirt.StoragePool) error {
	if err := vc.Libvirt.StoragePoolDestroy(pool); err != nil {
		return fmt.Errorf("停止存储池失败: %w", err)
	}
	return nil
}

// RefreshStoragePool 刷新存储池
func (vc *VirtConn) RefreshStoragePool(pool libvirt.StoragePool) error {
	if err := vc.Libvirt.StoragePoolRefresh(pool, 0); err != nil {
		return fmt.Errorf("刷新存储池失败: %w", err)
	}
	return nil
}

// InitSystemStoragePool 初始化多个存储池
func (vc *VirtConn) InitSystemStoragePool(params ...xmlDefine.PoolTemplateParams) error {
	for _, param := range params {
		// 获取存储池
		pool, err := vc.GetStoragePool(param.Name)
		if err != nil {
			logger.ErrorString("libvirt", "初始化存储池", err.Error())
			return err
		}

		// 如果存储池不存在，则创建
		if (pool == libvirt.StoragePool{}) {
			storagePool, err := vc.CreateStoragePool(&param)
			if err != nil {
				return fmt.Errorf("创建存储池 %s 失败: %w", param.Name, err)
			}
			if err := vc.StartStoragePool(storagePool); err != nil {
				return fmt.Errorf("启动存储池 %s 失败: %w", param.Name, err)
			}
			logger.InfoString("libvirt", "初始化存储池", fmt.Sprintf("成功创建并启动存储池 %s", param.Name))
		} else {
			logger.WarnString("libvirt", "初始化存储池", fmt.Sprintf("存储池 %s 已存在", param.Name))
		}
	}
	return nil
}

// GetStoragePool 根据 UUID、UUID字符串或名称获取存储池
func (vc *VirtConn) GetStoragePool(identifier any) (libvirt.StoragePool, error) {
	pools, err := vc.ListAllStoragePools()
	if err != nil {
		logger.ErrorString("libvirt", "获取存储池失败", err.Error())
		return libvirt.StoragePool{}, err
	}

	for _, pool := range pools {
		switch id := identifier.(type) {
		case libvirt.UUID:
			if pool.UUID == id {
				return pool, nil
			}
		case string:
			// 检查是否为UUID字符串格式
			if helpers.IsUUIDString(id) {
				uuid, err := helpers.UUIDStringToBytes(id)
				if err != nil {
					continue
				}
				if pool.UUID == uuid {
					return pool, nil
				}
			} else if pool.Name == id {
				// 如果不是UUID格式，则按名称匹配
				return pool, nil
			}
		default:
			return libvirt.StoragePool{}, fmt.Errorf("无效的标识符类型: %T", identifier)
		}
	}

	return libvirt.StoragePool{}, nil
}

// SetStoragePoolAutostart 设置存储池自动启动
func (vc *VirtConn) SetStoragePoolAutostart(pool libvirt.StoragePool, autostart bool) error {
	var autostartFlag int32 = 0
	if autostart {
		autostartFlag = 1
	}

	if err := vc.Libvirt.StoragePoolSetAutostart(pool, autostartFlag); err != nil {
		return fmt.Errorf("设置存储池自动启动失败: %w", err)
	}
	return nil
}

// GetStoragePoolAutostart 获取存储池自动启动状态
func (vc *VirtConn) GetStoragePoolAutostart(pool libvirt.StoragePool) (bool, error) {
	autostart, err := vc.Libvirt.StoragePoolGetAutostart(pool)
	if err != nil {
		return false, fmt.Errorf("获取存储池自动启动状态失败: %w", err)
	}

	return autostart == 1, nil
}
