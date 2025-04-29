// filepath: /home/horonlee/projects/govirt/pkg/libvirtd/network.go
package libvirtd

import (
	"fmt"
	"govirt/pkg/helpers"
	"govirt/pkg/logger"
	"govirt/pkg/xmlDefine"

	"github.com/digitalocean/go-libvirt"
)

// ListAllNetworks 列出所有网络
func (vc *VirtConn) ListAllNetworks() ([]libvirt.Network, error) {
	networks, _, err := vc.Libvirt.ConnectListAllNetworks(1, 0)
	if err != nil {
		return nil, err
	}
	return networks, nil
}

// CreateNetwork 创建网络
func (vc *VirtConn) CreateNetwork(params *xmlDefine.NetworkTemplateParams) (libvirt.Network, error) {
	// 为所有未设置的字段应用默认值
	xmlDefine.SetDefaults(params)

	// 渲染XML模板
	xmlStr, err := xmlDefine.RenderTemplate(xmlDefine.NetworkTemplate, params)
	if err != nil {
		return libvirt.Network{}, fmt.Errorf("渲染XML失败: %w", err)
	}

	// 定义网络
	network, err := vc.Libvirt.NetworkDefineXMLFlags(xmlStr, 0)
	if err != nil {
		return libvirt.Network{}, fmt.Errorf("定义存储池失败: %w", err)
	}

	// 设置自启动
	if err := vc.SetNetworkAutostart(network, params.Autostart); err != nil {
		return libvirt.Network{}, fmt.Errorf("设置自启动失败: %w", err)
	}
	return network, nil
}

// ActiveNetwork 启动网络
func (vc *VirtConn) ActiveNetwork(network libvirt.Network) error {
	if err := vc.Libvirt.NetworkCreate(network); err != nil {
		return fmt.Errorf("启动网络失败: %w", err)
	}
	return nil
}

// DeleteNetwork 删除网络
func (vc *VirtConn) DeleteNetwork(network libvirt.Network) error {
	if err := vc.Libvirt.NetworkDestroy(network); err != nil {
		return fmt.Errorf("停止网络失败: %w", err)
	}
	if err := vc.Libvirt.NetworkUndefine(network); err != nil {
		return fmt.Errorf("删除网络失败: %w", err)
	}
	return nil
}

// InitSystemNetwork 初始化系统网络，支持批量初始化
func (vc *VirtConn) InitSystemNetwork(params ...xmlDefine.NetworkTemplateParams) error {
	for _, param := range params {
		// 获取网络
		network, err := vc.GetNetwork(param.Name)
		if err != nil {
			logger.ErrorString("libvirt", "初始化网络", err.Error())
			return err
		}

		// 如果网络不存在，则创建
		if (network == libvirt.Network{}) {
			newNetwork, err := vc.CreateNetwork(&param)
			if err != nil {
				return fmt.Errorf("创建网络 %s 失败: %w", param.Name, err)
			}
			if err := vc.ActiveNetwork(newNetwork); err != nil {
				return fmt.Errorf("启动网络 %s 失败: %w", param.Name, err)
			}
			logger.InfoString("libvirt", "初始化网络", fmt.Sprintf("成功创建并启动网络 %s", param.Name))
		} else {
			logger.WarnString("libvirt", "初始化网络", fmt.Sprintf("网络 %s 已存在", param.Name))
		}
	}
	return nil
}

// GetNetwork 根据 UUID、UUID字符串或名称获取网络
func (vc *VirtConn) GetNetwork(identifier any) (libvirt.Network, error) {
	networks, err := vc.ListAllNetworks()
	if err != nil {
		logger.ErrorString("libvirt", "获取所有网络失败", err.Error())
		return libvirt.Network{}, err
	}

	for _, network := range networks {
		switch id := identifier.(type) {
		case libvirt.UUID:
			if network.UUID == id {
				return network, nil
			}
		case string:
			// 检查是否为UUID字符串格式
			if helpers.IsUUIDString(id) {
				uuid, err := helpers.UUIDStringToBytes(id)
				if err != nil {
					continue
				}
				if network.UUID == uuid {
					return network, nil
				}
			} else if network.Name == id {
				// 如果不是UUID格式，则按名称匹配
				return network, nil
			}
		default:
			return libvirt.Network{}, fmt.Errorf("无效的标识符类型: %T", identifier)
		}
	}

	return libvirt.Network{}, nil
}

// SetNetworkAutostart 设置网络自动启动
func (vc *VirtConn) SetNetworkAutostart(network libvirt.Network, autostart bool) error {
	var autostartFlag int32 = 0
	if autostart {
		autostartFlag = 1
	}

	if err := vc.Libvirt.NetworkSetAutostart(network, autostartFlag); err != nil {
		return fmt.Errorf("设置网络自动启动失败: %w", err)
	}
	return nil
}

// GetNetworkAutostart 获取网络自动启动状态
func (vc *VirtConn) GetNetworkAutostart(network libvirt.Network) (bool, error) {
	autostart, err := vc.Libvirt.NetworkGetAutostart(network)
	if err != nil {
		return false, fmt.Errorf("获取网络自动启动状态失败: %w", err)
	}

	return autostart == 1, nil
}
