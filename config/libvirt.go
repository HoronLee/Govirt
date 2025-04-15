// Package config libvirt配置文件
package config

import "govirt/pkg/config"

func init() {
	config.Add("libvirt", func() map[string]any {
		return map[string]any{
			"conURI": config.Env("CON_URI", "qemu:///system"),
		}
	})
}
