package libvirt

import (
	"govirt/pkg/helpers"
	"govirt/pkg/libvirt"
	"govirt/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListAllDomains 列出所有域
func (ctrl *LibvirtController) ListAllDomains(c *gin.Context) {
	domains, err := libvirt.ListAllDomains()
	if err != nil {
		response.Error(c, err, "列出所有域失败")
		return
	}
	response.Data(c, libvirt.FormatDomains(domains))
}

// GetDomainState 获取指定域的状态
func (ctrl *LibvirtController) GetDomainStateByUUID(c *gin.Context) {
	uuid, err := helpers.UUIDStringToBytes(c.Query("uuid"))
	state, err := libvirt.GetDomainStateByUUID(uuid)
	if err != nil {
		response.Error(c, err, "获取域失败")
		return
	}
	response.Data(c, libvirt.DomainStateToString(state))
}

// UpdateDomainState 更新指定域的状态
func (ctrl *LibvirtController) UpdateDomainStateByUUID(c *gin.Context) {
	uuid, err := helpers.UUIDStringToBytes(c.Query("uuid"))
	if err != nil {
		response.Error(c, err)
		return
	}
	op := libvirt.StringToDomainOperation(c.Query("operation"))
	_, err = libvirt.UpdateDomainStateByUUID(uuid, op)
	if err != nil {
		response.Error(c, err, "更新域状态失败")
		return
	}
	response.Success(c)
}

// GetDomainDefineXML 获取指定域的定义XML
func (ctrl *LibvirtController) GetDomainDefineXML(c *gin.Context) {
	uuid, err := helpers.UUIDStringToBytes(c.Query("uuid"))
	if err != nil {
		response.Error(c, err)
		return
	}
	domain, _ := libvirt.GetDomainByUUID(uuid)
	xml, err := libvirt.GetDomainXMLDesc(domain)
	if err != nil {
		response.Error(c, err, "获取域定义XML失败")
		return
	}
	response.Data(c, xml)
}
