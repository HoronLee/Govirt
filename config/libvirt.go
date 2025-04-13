// Package config libvirt配置文件
package config

import "gohub/pkg/config"

// init 函数在包初始化时自动执行，用于设置应用的配置信息
func init() {
	// 调用 config 包的 Add 方法，添加一个名为 "app" 的配置项
	config.Add("libvirt", func() map[string]any {
		// 返回一个包含应用配置信息的 map
		return map[string]any{
			"conURI": config.Env("CON_URI", "qemu:///system"),
		}
	})
}
