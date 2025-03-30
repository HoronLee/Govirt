// Package requests 处理请求数据和表单验证
package requests

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// ValidatorFunc 验证函数类型
type ValidatorFunc func(any, *gin.Context) map[string][]string

// Validate 控制器里调用示例：
//
//	if ok := requests.Validate(c, &requests.UserSaveRequest{}, requests.UserSave); !ok {
//	    return
//	}
func Validate(c *gin.Context, obj any, handler ValidatorFunc) bool {

	// 1. 解析请求，支持 JSON 数据、表单请求和 URL Query
	if err := c.ShouldBind(obj); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。",
			"error":   err.Error(),
		})
		fmt.Println(err.Error())
		return false
	}

	// 2. 表单验证
	errs := handler(obj, c)

	// 3. 判断验证是否通过
	if len(errs) > 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "请求验证不通过，具体请查看 errors",
			"errors":  errs,
		})
		return false
	}

	return true
}
func validate(data any, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {

	// 配置选项
	opts := govalidator.Options{
		Data:          data,     // 要验证的数据
		Rules:         rules,    // 验证规则
		TagIdentifier: "valid",  // 模型中的 Struct 标签标识符，用于标识需要验证的字段
		Messages:      messages, // 自定义的错误消息
	}

	// 开始验证
	return govalidator.New(opts).ValidateStruct() // 创建一个新的验证器并验证结构体，返回验证结果
}
