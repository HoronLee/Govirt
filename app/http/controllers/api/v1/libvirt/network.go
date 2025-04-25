package libvirt

import (
	"govirt/pkg/helpers"
	"govirt/pkg/network"
	"govirt/pkg/response"
	"govirt/pkg/xmlDefine"

	"github.com/gin-gonic/gin"
)

// ListAllNetworks 列出所有网络
func (ctrl *LibvirtController) ListAllNetworks(c *gin.Context) {
	networks, err := network.ListAllNetworks()
	if err != nil {
		response.Error(c, err, "列出所有网络失败")
		return
	}
	response.Data(c, helpers.FormatStructSlice(networks))
}

// CreateNetwork 创建网络
func (ctrl *LibvirtController) CreateNetwork(c *gin.Context) {
	// 获取请求参数
	var params xmlDefine.NetworkTemplateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, err, "请求参数错误")
		return
	}

	// 创建网络
	nw, err := network.CreateNetwork(&params)
	if err != nil {
		response.Error(c, err, "创建网络失败")
		return
	}

	// 启动网络
	if err := network.ActiveNetwork(nw); err != nil {
		response.Error(c, err, "启动网络失败")
		return
	}

	response.Success(c)
}

// DeketeNetwork 删除网络
func (ctrl *LibvirtController) DeleteNetwork(c *gin.Context) {
	ni := c.Query("network_identifier")
	if ni == "" {
		response.Error(c, nil, "网络名称不能为空")
		return
	}
	nw, err := network.GetNetwork(ni)
	if err != nil {
		response.Error(c, nil, "网络不存在")
		return
	}
	if err := network.DeleteNetwork(nw); err != nil {
		response.Error(c, err, "删除网络失败")
		return
	}

	response.Success(c)
}

// ActiveNetwork 启动网络
func (ctrl *LibvirtController) ActiveNetwork(c *gin.Context) {
	ni := c.Query("network_identifier")
	if ni == "" {
		response.Error(c, nil, "网络名称不能为空")
		return
	}
	nw, err := network.GetNetwork(ni)
	if err != nil {
		response.Error(c, nil, "网络不存在")
		return
	}
	if err := network.ActiveNetwork(nw); err != nil {
		response.Error(c, err, "启动网络失败")
		return
	}

	response.Success(c)
}
