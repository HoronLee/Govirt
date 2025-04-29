package helpers

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/digitalocean/go-libvirt"
)

// UUIDBytesToString 将UUID字节数组转换为标准字符串格式
func UUIDBytesToString(uuid []byte) string {
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		uuid[0], uuid[1], uuid[2], uuid[3],
		uuid[4], uuid[5], uuid[6], uuid[7],
		uuid[8], uuid[9], uuid[10], uuid[11],
		uuid[12], uuid[13], uuid[14], uuid[15])
}

// UUIDStringToBytes 将标准UUID字符串转换为字节数组
func UUIDStringToBytes(uuidStr string) (libvirt.UUID, error) {
	var uuid [16]byte
	_, err := fmt.Sscanf(uuidStr, "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		&uuid[0], &uuid[1], &uuid[2], &uuid[3],
		&uuid[4], &uuid[5], &uuid[6], &uuid[7],
		&uuid[8], &uuid[9], &uuid[10], &uuid[11],
		&uuid[12], &uuid[13], &uuid[14], &uuid[15])
	if err != nil {
		return libvirt.UUID{}, fmt.Errorf("无效的UUID格式")
	}
	return uuid, nil
}

// FormatStructSlice 将结构体切片中的每个元素转换为格式化后的映射
func FormatStructSlice(slice any) []map[string]any {
	var result []map[string]any

	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		return result
	}

	for i := range val.Len() {
		elem := val.Index(i).Interface()
		formattedElem := FormatUUIDInStruct(elem)
		result = append(result, formattedElem)
	}

	return result
}

// FormatUUIDInStruct 将结构体中任何类型为libvirt.UUID的字段转换为字符串格式
// 返回一个新的map，保留原始字段名和原始值，但UUID字段被转换为字符串
func FormatUUIDInStruct(obj any) map[string]any {
	result := make(map[string]any)

	// 获取对象的反射值
	val := reflect.ValueOf(obj)

	// 如果是指针，获取指针指向的元素
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 只处理结构体
	if val.Kind() != reflect.Struct {
		return result
	}

	// 获取结构体的类型
	typ := val.Type()

	// 遍历结构体的所有字段
	for i := range val.NumField() {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		// 检查字段是否为UUID类型
		if field.Type() == reflect.TypeOf(libvirt.UUID{}) {
			// 将UUID转换为字符串
			uuid := field.Interface().(libvirt.UUID)
			result[fieldName] = UUIDBytesToString(uuid[:])
		} else {
			// 保持原始值
			result[fieldName] = field.Interface()
		}
	}

	return result
}

// GenerateUUIDString 生成一个新的UUID字符串
func GenerateUUIDString() string {
	uuid := [16]byte{}
	for i := range uuid {
		uuid[i] = byte(rand.Intn(256))
	}
	return UUIDBytesToString(uuid[:])
}

// IsUUIDString 检查字符串是否为有效的UUID格式
func IsUUIDString(uuidStr string) bool {
	if len(uuidStr) != 36 {
		return false
	}

	var uuid [16]byte
	_, err := fmt.Sscanf(uuidStr, "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		&uuid[0], &uuid[1], &uuid[2], &uuid[3],
		&uuid[4], &uuid[5], &uuid[6], &uuid[7],
		&uuid[8], &uuid[9], &uuid[10], &uuid[11],
		&uuid[12], &uuid[13], &uuid[14], &uuid[15])
	return err == nil
}
