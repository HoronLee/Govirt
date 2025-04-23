package config

import "govirt/pkg/config"

func init() {
	config.Add("network", func() map[string]any {
		return map[string]any{
			"internal": map[string]any{
				"name":        config.Env("INTERNAL_NETWORK_NAME", "internal"),
				"domainName":  config.Env("INTERNAL_NETWORK_DOMAIN_NAME", "internal"),
				"forwardMode": config.Env("INTERNAL_NETWORK_FORWARD_MODE", "nat"),
				"ip":          config.Env("INTERNAL_NETWORK_IP", "192.168.200.1"),
				// "netmask":     config.Env("INTERNAL_NETWORK_NETMASK", "255.255.255.0"),
				"dhcpStart": config.Env("INTERNAL_NETWORK_DHCP_START", "192.168.200.2"),
				"dhcpEnd":   config.Env("INTERNAL_NETWORK_DHCP_END", "192.168.200.254"),
			},
			"external": map[string]any{
				"name":        config.Env("EXTERNAL_NETWORK_NAME", "external"),
				"domainName":  config.Env("EXTERNAL_NETWORK_DOMAIN_NAME", "external"),
				"forwardMode": config.Env("EXTERNAL_NETWORK_FORWARD_MODE", "route"),
				"ip":          config.Env("EXTERNAL_NETWORK_IP", "192.168.250.1"),
				// "netmask":     config.Env("EXTERNAL_NETWORK_NETMASK", "255.255.255.0"),
				"dhcpStart": config.Env("EXTERNAL_NETWORK_DHCP_START", "192.168.250.2"),
				"dhcpEnd":   config.Env("EXTERNAL_NETWORK_DHCP_END", "192.168.250.254"),
			},
		}
	})
}
