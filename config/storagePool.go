package config

import "govirt/pkg/config"

func init() {
	config.Add("pool", func() map[string]any {
		return map[string]any{
			"image": map[string]any{
				"name": config.Env("IMAGE_POOL_NAME", "images"),
				"path": config.Env("IMAGE_POOL_PATH", "/var/lib/libvirt/images"),
			},
			"volume": map[string]any{
				"name": config.Env("VOLUME_POOL_NAME", "volumes"),
				"path": config.Env("VOLUME_POOL_PATH", "/var/lib/libvirt/volumes"),
			},
		}
	})
}
