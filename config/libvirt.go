// Package config libvirt配置文件
package config

import (
	"govirt/pkg/config"
)

func init() {
	config.Add("libvirt", func() map[string]any {
		return map[string]any{
			"hostName": config.Env("HOST_NAME", "R430"),
			"hostUUID": config.Env("HOST_UUID", "12345678-1234-5678-1234-567812345678"),
			"conURI":   config.Env("CON_URI", "qemu:///system"),
		}
	})
}
