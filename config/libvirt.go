// Package config libvirt配置文件
package config

import "govirt/pkg/config"

func init() {
	config.Add("libvirt", func() map[string]any {
		return map[string]any{
			"hostName":       config.Env("HOST_NAME", "R430"),
			"conURI":         config.Env("CON_URI", "qemu:///system"),
			"imagePollPath":  config.Env("IMAGE_POOL_PATH", "/var/lib/libvirt/images"),
			"volumePoolPath": config.Env("VOLUME_POOL_PATH", "/var/lib/libvirt/volumes"),
		}
	})
}
