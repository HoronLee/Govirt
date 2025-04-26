package v1

import (
	"govirt/app/models/apikey"
	"govirt/app/requests"
	"govirt/pkg/database"
	"govirt/pkg/helpers"
	"govirt/pkg/response"

	"github.com/gin-gonic/gin"
)

type ApikeyController struct {
	BaseAPIController
}

func (ctrl *ApikeyController) ListApikey(c *gin.Context) {
	var apikeys []apikey.Apikey
	database.DB.Find(&apikeys)

	result := make([]struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}, len(apikeys))

	for i, v := range apikeys {
		result[i] = struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}{
			ID:   uint(v.ID),
			Name: v.Name,
		}
	}

	response.Data(c, result)
}

func (ctrl *ApikeyController) CreateApikey(c *gin.Context) {
	request := requests.CreateApikeyRequest{}
	if ok := requests.Validate(c, &request, requests.CreateApikey); !ok {
		return
	}
	if (&apikey.Apikey{Name: request.Name}).IsExist() {
		response.Abort403(c, "API Key 已存在")
		return
	}
	key := helpers.RandomString(64)
	apikey := apikey.Apikey{
		Name: request.Name,
		Key:  key,
	}
	apikey.Create()
	response.Data(c, gin.H{
		"name": apikey.Name,
		"key":  key,
	})
}

func (ctrl *ApikeyController) DeleteApikey(c *gin.Context) {
	apikeyModel, err := apikey.GetFromName(c.Param("name"))
	if err != nil {
		response.Abort404(c)
		return
	}
	if apikeyModel.ID == 0 {
		response.Abort404(c)
		return
	}

	rowsAffected, err := apikeyModel.Delete()
	if err != nil {
		response.Abort500(c, "删除失败，请稍后尝试~")
		return
	}
	if rowsAffected > 0 {
		response.Success(c)
	} else {
		response.Abort500(c, "删除失败，请稍后尝试~")
	}
}
