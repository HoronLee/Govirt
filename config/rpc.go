package config

import (
	"govirt/pkg/config"
)

func init() {
	config.Add("rpc", func() map[string]any {
		return map[string]any{
			"enable":  config.Env("RPC_ENABLE", "true"),
			"address": config.Env("RPC_ADDRESS", "0.0.0.0"),
			"port":    config.Env("RPC_PORT", "8000"),
			"name":    config.Env("RPC_NAME", "govirt"),
			"timeout": config.Env("RPC_TIMEOUT", "1s"),
		}
	})
}
