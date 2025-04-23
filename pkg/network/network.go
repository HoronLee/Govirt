package network

import (
	"fmt"
	"govirt/pkg/libvirtd"
	"govirt/pkg/xmlDefine"

	"github.com/digitalocean/go-libvirt"
	"github.com/google/uuid"
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

	// 如果未提供UUID，则自动生成一个
	if params.UUID == "" {
		params.UUID = uuid.New().String()
	}
	// 渲染XML模板
	xmlStr, err := xmlDefine.RenderTemplate(xmlDefine.PoolTemplate, params)
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
