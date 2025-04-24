package network

import (
	"fmt"
	"govirt/pkg/libvirtd"
	"govirt/pkg/logger"
	"govirt/pkg/xmlDefine"

	"github.com/digitalocean/go-libvirt"
)

// InitSystemNetwork 初始化系统网络，支持批量初始化
func InitSystemNetwork(params ...xmlDefine.NetworkTemplateParams) error {
	for _, param := range params {
		// 获取网络
		network, err := GetNetwork(param.Name)
		if err != nil {
			logger.ErrorString("libvirt", "初始化网络", err.Error())
			return err
		}

		// 如果网络不存在，则创建
		if (network == libvirt.Network{}) {
			newNetwork, err := CreateNetwork(&param)
			if err != nil {
				return fmt.Errorf("创建网络 %s 失败: %w", param.Name, err)
			}
			if err := StartNetwork(newNetwork); err != nil {
				return fmt.Errorf("启动网络 %s 失败: %w", param.Name, err)
			}
			logger.InfoString("libvirt", "初始化网络", fmt.Sprintf("成功创建并启动网络 %s", param.Name))
		} else {
			logger.WarnString("libvirt", "初始化网络", fmt.Sprintf("网络 %s 已存在", param.Name))
		}
	}
	return nil
}

// GetNetwork 根据 UUID 或名称获取网络
func GetNetwork(identifier any) (libvirt.Network, error) {
	networks, err := ListAllNetworks()
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
			if network.Name == id {
				return network, nil
			}
		default:
			return libvirt.Network{}, fmt.Errorf("无效的标识符类型: %T", identifier)
		}
	}

	return libvirt.Network{}, nil
}

// SetNetworkAutostart 设置网络自动启动
func SetNetworkAutostart(network libvirt.Network, autostart bool) error {
	var autostartFlag int32 = 0
	if autostart {
		autostartFlag = 1
	}

	if err := libvirtd.Connection.NetworkSetAutostart(network, autostartFlag); err != nil {
		return fmt.Errorf("设置网络自动启动失败: %w", err)
	}
	return nil
}

// GetNetworkAutostart 获取网络自动启动状态
func GetNetworkAutostart(network libvirt.Network) (bool, error) {
	autostart, err := libvirtd.Connection.NetworkGetAutostart(network)
	if err != nil {
		return false, fmt.Errorf("获取网络自动启动状态失败: %w", err)
	}

	return autostart == 1, nil
}
