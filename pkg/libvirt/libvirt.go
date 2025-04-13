package libvirt

import (
	"gohub/pkg/config"
	"gohub/pkg/helpers"
	"gohub/pkg/logger"
	"net/url"

	"github.com/digitalocean/go-libvirt"
)

var (
	// Libvirt 全局 libvirt 连接对象
	Libvirt *libvirt.Libvirt
)

// 此函数用于连接到 libvirt 并列出所有域的信息
func InitLibvirt() {
	// 解析 libvirt 的 URI
	uri, _ := url.Parse(config.GetString("libvirt.uri"))
	// 连接到 libvirt
	var err error
	Libvirt, err = libvirt.ConnectToURI(uri)
	if err != nil {
		logger.FatalString("libvirt", "连接到libvirtd失败", err.Error())
	}

	// 使用 defer 确保在函数返回时断开连接
	// defer func() {
	// 	if err = Libvirt.Disconnect(); err != nil {
	// 		logger.FatalJSON("libvirt", "无法断开与libvirt的连接", err)
	// 	}
	// }()

	// logger.InfoString("libvirt", "libvirt版本", GetLibVersion())

	// domains := ListAllDomains()
	// // 遍历所有域并打印信息
	// for _, d := range domains {
	// 	fmt.Printf("%d\t%s\t%x\n", d.ID, d.Name, d.UUID)
	// }
}

// 获取 libvirt 版本
func GetLibVersion() (string, error) {
	v, err := Libvirt.ConnectGetLibVersion()
	if err != nil {
		logger.FatalString("libvirt", "获取libvirt版本失败", err.Error())
		return "", err
	}
	version := helpers.ToStandardVersion(v)
	return version, nil
}
