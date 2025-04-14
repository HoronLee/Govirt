package bootstrap

import (
	"govirt/pkg/config"
	libvirtPKG "govirt/pkg/libvirt"
	"govirt/pkg/logger"
	"net/url"

	"github.com/digitalocean/go-libvirt"
)

// 此函数用于连接到 libvirt 并列出所有域的信息
func InitLibvirt() {
	// 解析 libvirt 的 URI
	uri, _ := url.Parse(config.GetString("libvirt.uri"))
	// 连接到 libvirt
	var err error
	libvirtPKG.Libvirt, err = libvirt.ConnectToURI(uri)
	if err != nil {
		logger.FatalString("libvirt", "连接到libvirtd失败", err.Error())
	}

	// 使用 defer 确保在函数返回时断开连接
	// defer func() {
	// 	if err = Libvirt.Disconnect(); err != nil {
	// 		logger.FatalJSON("libvirt", "无法断开与libvirt的连接", err)
	// 	}
	// }()
}
