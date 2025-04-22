package libvirt

import (
	"govirt/pkg/helpers"

	"github.com/digitalocean/go-libvirt"
)

// FormatPools 格式化池信息
func FormatPools(pools []libvirt.StoragePool) []map[string]any {
	var formattedPools []map[string]any
	for _, d := range pools {
		formattedPools = append(formattedPools, map[string]any{
			"Name": d.Name,
			"UUID": helpers.UUIDBytesToString(d.UUID),
		})
	}
	return formattedPools
}
