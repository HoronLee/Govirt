package xmlDefine

import (
	"bytes"
	"fmt"
	"govirt/pkg/config"
	"reflect"
	"strconv"
	"text/template"
)

// 定义模板函数映射
var templateFuncs = template.FuncMap{
	"notEmpty": func(s string) bool {
		return s != ""
	},
	// greater类型的方法用于处理负数值自动使用默认值，用法如下
	// {{if greaterThanZero .MaxMem}}
	// <memory unit="KiB">{{.MaxMem}}</memory>
	// {{else}}
	// <memory unit="KiB">1048576</memory> <!-- 默认1GB -->
	// {{end}}
	"greaterThanZero": func(n int) bool {
		return n > 0
	},
	"greaterThanZero64": func(n int64) bool {
		return n > 0
	},
	// eq和ne用于判断两个数值是否相等，下面是一个示例用法
	// {{if eq .VncPort -1}}
	// <!-- 使用自动端口 -->
	// {{else}}
	// <!-- 使用指定端口 -->
	// {{end}}
	"eq": func(a, b any) bool {
		return a == b
	},
	"ne": func(a, b any) bool {
		return a != b
	},
}

// RenderTemplate 通用模板渲染函数，接收任意模板字符串和数据结构
// templateStr: 包含模板变量的模板字符串
// data: 包含模板变量值的结构体实例
func RenderTemplate(templateStr string, data any) (string, error) {
	tmpl, err := template.New("template").Funcs(templateFuncs).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("执行模板渲染失败: %w", err)
	}

	return buf.String(), nil
}

// 如果需要缓存已解析的模板以提高性能
var templateCache = make(map[string]*template.Template)
var commonTemplates = map[string]string{
	"domain": DomainTemplate,
	// "network": networkTemplate,
	// 可以添加其他常用模板
}

// RenderCachedTemplate 使用缓存的模板进行渲染
func RenderCachedTemplate(templateName string, data any) (string, error) {
	var tmpl *template.Template
	var ok bool

	// 检查模板是否在缓存中
	tmpl, ok = templateCache[templateName]
	if !ok {
		// 获取模板内容
		templateStr, exists := commonTemplates[templateName]
		if !exists {
			return "", fmt.Errorf("未找到名为 %s 的模板", templateName)
		}

		// 解析模板并添加到缓存
		var err error
		tmpl, err = template.New(templateName).Funcs(templateFuncs).Parse(templateStr)
		if err != nil {
			return "", fmt.Errorf("解析模板 %s 失败: %w", templateName, err)
		}
		templateCache[templateName] = tmpl
	}

	// 执行模板渲染
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("执行模板 %s 渲染失败: %w", templateName, err)
	}

	return buf.String(), nil
}

// SetDefaults 为结构体设置默认值
// 优先根据 config 标签从提供的 ConfigGetter 获取值，
// 如果 config 标签不存在或获取失败，则根据 default 标签设置默认值。
func SetDefaults(obj any) {
	// 获取传入对象的反射值
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	// 遍历所有字段
	for i := range t.NumField() {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 只处理零值字段
		if fieldValue.IsZero() {
			configApplied := false
			// 检查 config 标签
			configKey := field.Tag.Get("config")
			if configKey != "" {
				configValue := config.Get(configKey)
				// 假设 Get 返回非空字符串表示成功获取
				if configValue != "" {
					setDefaultByType(fieldValue, configValue)
					configApplied = true
				}
			}

			// 如果没有应用 config 值，则检查 default 标签
			if !configApplied {
				defaultValue := field.Tag.Get("default")
				if defaultValue != "" {
					setDefaultByType(fieldValue, defaultValue)
				}
			}
		}
	}
}

// setDefaultByType 根据字段类型设置默认值
func setDefaultByType(fieldValue reflect.Value, value string) {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, err := strconv.ParseInt(value, 10, 64); err == nil {
			fieldValue.SetInt(val)
		} // 可以在此处添加错误处理或日志记录
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := strconv.ParseUint(value, 10, 64); err == nil {
			fieldValue.SetUint(val)
		}
	case reflect.Float32, reflect.Float64:
		if val, err := strconv.ParseFloat(value, 64); err == nil {
			fieldValue.SetFloat(val)
		}
	case reflect.Bool:
		if val, err := strconv.ParseBool(value); err == nil {
			fieldValue.SetBool(val)
		}
	}
}
