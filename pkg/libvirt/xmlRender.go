package libvirt

import (
	"bytes"
	"fmt"
	"text/template"
)

// RenderTemplate 通用模板渲染函数，接收任意模板字符串和数据结构
// templateStr: 包含模板变量的模板字符串
// data: 包含模板变量值的结构体实例
func RenderTemplate(templateStr string, data any) (string, error) {
	tmpl, err := template.New("template").Parse(templateStr)
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
	"domain": domainTemplate,
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
		tmpl, err = template.New(templateName).Parse(templateStr)
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

// RenderDomainXML 渲染虚拟机域定义的XML
// 这是一个针对DomainTemplateParams的便捷包装函数
func RenderDomainXML(params any) (string, error) {
	return RenderTemplate(domainTemplate, params)
}
