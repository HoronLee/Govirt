package libvirt

import (
	"govirt/pkg/helpers"
	"govirt/pkg/libvirtd"
	"govirt/pkg/response"
	"govirt/pkg/xmlDefine"

	"github.com/gin-gonic/gin"
)

func (ctrl *LibvirtController) ListAllStoragePools(c *gin.Context) {
	pools, err := libvirtd.Conn.ListAllStoragePools()
	if err != nil {
		response.Error(c, err, "列出所有存储池失败")
		return
	}
	c.JSON(200, helpers.FormatStructSlice(pools))
}

func (ctrl *LibvirtController) CreateStartStoragePool(c *gin.Context) {
	// 解析请求参数
	var params xmlDefine.PoolTemplateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, err, "解析请求参数失败")
		return
	}

	// 创建存储池
	pool, err := libvirtd.Conn.CreateStoragePool(&params)
	if err != nil {
		response.Error(c, err, "创建存储池失败")
		return
	}

	// 启动存储池
	if err := libvirtd.Conn.StartStoragePool(pool); err != nil {
		response.Error(c, err, "启动存储池失败")
		return
	}
	response.Created(c, helpers.FormatUUIDInStruct(pool))
}

func (ctrl *LibvirtController) DeleteStoragePool(c *gin.Context) {
	uuid := c.Query("pool_identifier")
	pool, err := libvirtd.Conn.GetStoragePool(uuid)
	if err != nil {
		response.Error(c, err, "获取存储池失败")
		return
	}
	if err := libvirtd.Conn.StopStoragePool(pool); err != nil {
		response.Error(c, err, "停止存储池失败")
		return
	}
	// 删除存储池
	if err := libvirtd.Conn.DeleteStoragePool(pool); err != nil {
		response.Error(c, err, "删除存储池失败")
		return
	}
	response.Success(c)
}
