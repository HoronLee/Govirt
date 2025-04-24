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
	needResult := c.Query("needResult")
	pool, err := storagePool.GetStoragePool(poolName)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	// 将字符串转换为整数
	resultNumInt, err := strconv.Atoi(needResult)
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

func (ctrl *LibvirtController) ListAllVolumesDetailsByPool(c *gin.Context) {
	poolName := c.Query("name")
	needResult := c.Query("needResult")
	pool, err := storagePool.GetStoragePool(poolName)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	// 将字符串转换为整数
	resultNumInt, err := strconv.Atoi(needResult)
	if err != nil {
		logger.WarnString("libvirt", "ListVolumesByPool", fmt.Sprintf("resultNum转换失败: %s", err.Error()))
		resultNumInt = 100
	}

	volumes, rRet, err := volume.ListVolumesDetail(pool, int32(resultNumInt), 0)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	data := struct {
		Volumes any    `json:"volumes"`
		RRet    uint32 `json:"rRet"`
	}{
		Volumes: volumes,
		RRet:    rRet,
	}
	response.Data(c, data)
}
