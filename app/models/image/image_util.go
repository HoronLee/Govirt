package image

import (
	"fmt"
	"govirt/pkg/database"
	"govirt/pkg/helpers"
)

func Get(idstr string) (image Image, err error) {
	result := database.DB.Where("id", idstr).First(&image)
	return image, result.Error
}

func GetByUUID(uuid string) (image Image, err error) {
	result := database.DB.Where("uuid = ?", uuid).First(&image)
	return image, result.Error
}

func GetBy(field, value string) (image Image, err error) {
	// 使用 fmt.Sprintf 来动态构建列名查询
	// 注意：如果 field 参数来自用户输入，需要进行严格验证以防止 SQL 注入
	// 在 GetByID 调用此函数时，field 固定为 "name"，是安全的。
	result := database.DB.Where(fmt.Sprintf("`%s` = ?", field), value).First(&image) // 在列名两边加上反引号 `` 更安全
	return image, result.Error
}

// GetByID 根据标识符（名称或UUID）查找镜像记录
func GetByID(identifier string) (image Image, err error) {
	if helpers.IsUUIDString(identifier) {
		// 直接调用 GetByUUID 可能更清晰
		result := database.DB.Where("uuid = ?", identifier).First(&image)
		err = result.Error
		// image, err = GetByUUID(identifier) // 或者保持原来的调用
	} else {
		// 直接查询 name 列
		result := database.DB.Where("name = ?", identifier).First(&image)
		err = result.Error
		// image, err = GetBy("name", identifier) // 或者保持原来的调用，但需确保 GetBy 已修复
	}
	if err != nil {
		// 可以考虑检查 err 是否为 gorm.ErrRecordNotFound 并返回更具体的错误信息
		return image, fmt.Errorf("查找镜像失败: %w", err)
	}
	return image, nil
}

func GetByStatus(status string) (image Image, err error) {
	result := database.DB.Where("status = ?", status).First(&image)
	return image, result.Error
}

func All() (images []Image, err error) {
	result := database.DB.Find(&images)
	return images, result.Error
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(Image{}).Where("? = ?", field, value).Count(&count)
	return count > 0
}
