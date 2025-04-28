package libvirt

import (
	"govirt/pkg/image"
	"govirt/pkg/response"

	"github.com/gin-gonic/gin"
)

func (ctrl *LibvirtController) CreateImageFromLocalFile(c *gin.Context) {
	var req struct {
		Name           string
		SourceFilePath string
		PoolName       string
		OSType         string
		Arch           string
		ImageType      string
		Description    string
		MinDisk        uint64
		MinRam         uint64
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, "请求参数错误")
		return
	}

	image, err := image.CreateImageFromLocalFile(
		req.Name,
		req.SourceFilePath,
		req.PoolName,
		req.OSType,
		req.Arch,
		req.ImageType,
		req.Description,
		req.MinDisk,
		req.MinRam,
	)
	if err != nil {
		response.Error(c, err, "创建镜像失败")
		return
	}

	response.Data(c, image)
}

func (ctrl *LibvirtController) DeleteImage(c *gin.Context) {
	idOrUUID := c.Query("image_id")

	if idOrUUID == "" {
		response.BadRequest(c, nil, "缺少镜像ID或UUID")
		return
	}

	err := image.DeleteImage(idOrUUID)
	if err != nil {
		response.Error(c, err, "删除镜像失败")
		return
	}

	response.Success(c)
}
