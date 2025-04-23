// Package config libvirt配置文件
package config

import (
	"govirt/pkg/config"
)

func init() {
	config.Add("libvirt", func() map[string]any {
		return map[string]any{
			"hostName": config.Env("HOST_NAME", "R430"),
			"conURI":   config.Env("CON_URI", "qemu:///system"),
			"pool": map[string]any{
				"image": map[string]any{
					"name": config.Env("IMAGE_POOL_NAME", "images"),
					"path": config.Env("IMAGE_POOL_PATH", "/var/lib/libvirt/images"),
				},
				"volume": map[string]any{
					"name": config.Env("VOLUME_POOL_NAME", "volumes"),
					"path": config.Env("VOLUME_POOL_PATH", "/var/lib/libvirt/volumes"),
				},
			},
			"network": map[string]any{
				"internal": map[string]any{
					"name": config.Env("INTERNAL_NETWORK_NAME", "internal"),
				},
				"external": map[string]any{
					"name": config.Env("EXTERNAL_NETWORK_NAME", "external"),
				},
			},
		}
	})
}
