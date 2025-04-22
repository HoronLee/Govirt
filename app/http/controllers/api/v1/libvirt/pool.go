package libvirt

import (
	"govirt/pkg/libvirt"
	"govirt/pkg/response"

	"github.com/gin-gonic/gin"
)

func (ctrl *LibvirtController) ListAllStoragePools(c *gin.Context) {
	pools, err := libvirt.ListAllStoragePools()
	if err != nil {
		response.Error(c, err, "列出所有存储池失败")
		return
	}
	c.JSON(200, libvirt.FormatPools(pools))
}
