package libvirt

import (
	"fmt"
	"govirt/pkg/logger"
	"govirt/pkg/response"
	"govirt/pkg/storagePool"
	"govirt/pkg/volume"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ctrl *LibvirtController) ListVolumesByPool(c *gin.Context) {
	poolName := c.Query("name")
	resultNum := c.Query("resultNum")
	pool, err := storagePool.GetStoragePool(poolName)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	// 将字符串转换为整数
	resultNumInt, err := strconv.Atoi(resultNum)
	if err != nil {
		logger.WarnString("libvirt", "ListVolumesByPool", fmt.Sprintf("resultNum转换失败: %s", err.Error()))
		resultNumInt = 100
	}

	volumes, err := volume.ListVolumes(pool, int32(resultNumInt))
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	response.Data(c, volumes)
}
