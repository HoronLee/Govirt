package libvirt

import (
	"govirt/pkg/response"
	"govirt/pkg/storagePool"
	"govirt/pkg/volume"
	"govirt/pkg/xmlDefine"

	"github.com/gin-gonic/gin"
)

func (ctrl *LibvirtController) ListVolumesSummary(c *gin.Context) {
	identifier := c.Query("pool_identifier")

	pool, err := storagePool.GetStoragePool(identifier)
	if err != nil {
		response.Abort500(c, "Failed to get storage pool: "+err.Error())
		return
	}

	volumes, err := volume.ListVolumesSummary(pool)
	if err != nil {
		response.Abort500(c, "Failed to list volume summaries: "+err.Error())
		return
	}

	response.Data(c, volumes)
}

func (ctrl *LibvirtController) ListVolumesDetails(c *gin.Context) {
	identifier := c.Query("pool_identifier")
	pool, err := storagePool.GetStoragePool(identifier)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	volumes, _, err := volume.ListVolumesDetails(pool, 0)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	response.Data(c, volumes)
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

// DeleteVolume 删除存储卷
func (ctrl *LibvirtController) DeleteVolume(c *gin.Context) {
	identifier := c.Query("pool_identifier")
	pool, err := storagePool.GetStoragePool(identifier)
	if err != nil {
		response.Error(c, err)
	}
	volumeName := c.Query("volume_name")
	if volumeName == "" {
		response.BadRequest(c, nil, "卷名称不能为空")
		return
	}
	err = volume.DeleteVolume(pool, volumeName, 0)
	if err != nil {
		response.Error(c, err, "删除存储卷失败")
		return
	}

	response.Success(c)
}

// CloneVolume 克隆存储卷
func (ctrl *LibvirtController) CloneVolume(c *gin.Context) {
	svn := c.Query("source_volume_name")
	if svn == "" {
		response.BadRequest(c, nil, "源卷名称不能为空")
		return
	}

	spi := c.Query("source_pool_identifier")
	dpi := c.DefaultQuery("destination_pool_identifier", spi)

	spool, err := storagePool.GetStoragePool(spi)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	dpool, err := storagePool.GetStoragePool(dpi)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	svol, err := volume.GetVolume(spool, svn)
	if err != nil {
		response.BadRequest(c, err, "源卷不存在")
		return
	}
	// 解析请求参数
	var params xmlDefine.VolumeTemplateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, err, "解析请求参数失败")
		return
	}

	// vol, err := volume.CloneVolume(spool, svn, dpool, &params, 0)
	vol, err := volume.CloneVolume(dpool, &params, svol, 0)
	if err != nil {
		response.Error(c, err, "克隆存储卷失败")
		return
	}

	response.Created(c, vol)
}
