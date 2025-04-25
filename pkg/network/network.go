package network

import (
	"fmt"
	"govirt/pkg/libvirtd"
	"govirt/pkg/xmlDefine"

	"github.com/digitalocean/go-libvirt"
)

// ListAllNetworks 列出所有网络
func ListAllNetworks() ([]libvirt.Network, error) {
	networks, _, err := libvirtd.Connection.ConnectListAllNetworks(1, libvirt.ConnectListNetworksActive|libvirt.ConnectListNetworksInactive)
	if err != nil {
		return nil, err
	}
	return networks, nil
}

// CreateNetwork 创建网络
func CreateNetwork(params *xmlDefine.NetworkTemplateParams) (libvirt.Network, error) {
	// 为所有未设置的字段应用默认值
	xmlDefine.SetDefaults(params)

	// 渲染XML模板
	xmlStr, err := xmlDefine.RenderTemplate(xmlDefine.NetworkTemplate, params)
	if err != nil {
		return libvirt.Network{}, fmt.Errorf("渲染XML失败: %w", err)
	}

	// 定义网络
	network, err := libvirtd.Connection.NetworkDefineXMLFlags(xmlStr, 0)
	if err != nil {
		return libvirt.Network{}, fmt.Errorf("定义存储池失败: %w", err)
	}

	// 设置自启动
	if err := SetNetworkAutostart(network, params.Autostart); err != nil {
		return libvirt.Network{}, fmt.Errorf("设置自启动失败: %w", err)
	}
	return network, nil
}

// ActiveNetwork 启动网络
func ActiveNetwork(network libvirt.Network) error {
	if err := libvirtd.Connection.NetworkCreate(network); err != nil {
		return fmt.Errorf("启动网络失败: %w", err)
	}
	return nil
}

// DeleteNetwork 删除网络
func DeleteNetwork(network libvirt.Network) error {
	if err := libvirtd.Connection.NetworkDestroy(network); err != nil {
		return fmt.Errorf("停止网络失败: %w", err)
	}
	if err := libvirtd.Connection.NetworkUndefine(network); err != nil {
		return fmt.Errorf("删除网络失败: %w", err)
	}
	return nil
}
