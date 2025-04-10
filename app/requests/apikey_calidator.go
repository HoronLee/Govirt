package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type CreateApikeyRequest struct {
	Name string `json:"name,omitempty" valid:"name"`
}

func CreateApikey(data any, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"name": []string{"required", "between:3,20"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:密钥名称不能为空",
			"between:密钥名称长度需在 3~20 之间",
		},
	}

	errs := validate(data, rules, messages)

	return errs
}
