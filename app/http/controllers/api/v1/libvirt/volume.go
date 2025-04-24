package libvirt

import (
	"govirt/pkg/response"
	"govirt/pkg/storagePool"
	"govirt/pkg/volume"
	"govirt/pkg/xmlDefine"

	"github.com/gin-gonic/gin"
)

func (ctrl *LibvirtController) ListVolumesSummaryByPool(c *gin.Context) {
	identifier := c.Query("pool_identifier")

	pool, err := storagePool.GetStoragePool(identifier)
	if err != nil {
		response.Abort500(c, "Failed to get storage pool: "+err.Error())
		return
	}

	resultNum, err := volume.GetVolumeNum(pool)
	if err != nil {
		response.Abort500(c, "Failed to get volume count: "+err.Error())
		return
	}
	volumes, err := volume.ListVolumesSummary(pool, resultNum)
	if err != nil {
		response.Abort500(c, "Failed to list volume summaries: "+err.Error())
		return
	}

	response.Data(c, volumes)
}

func (ctrl *LibvirtController) ListAllVolumesDetailsByPool(c *gin.Context) {
	identifier := c.Query("pool_identifier")
	pool, err := storagePool.GetStoragePool(identifier)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	volumes, rRet, err := volume.ListVolumesDetails(pool, 1, 0)
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

// CreateVolume 创建存储卷
func (ctrl *LibvirtController) CreateVolume(c *gin.Context) {
	identifier := c.Query("pool_identifier")
	pool, err := storagePool.GetStoragePool(identifier)
	if err != nil {
		response.Abort500(c, err.Error())
	}
	// 解析请求参数
	var params xmlDefine.VolumeTemplateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, err, "解析请求参数失败")
		return
	}

	vol, err := volume.CreateVolume(pool, &params, 0)
	if err != nil {
		response.Error(c, err, "创建存储卷失败")
		return
	}

	response.Created(c, vol)
}
