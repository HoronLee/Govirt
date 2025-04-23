package xmlDefine

import (
	"encoding/xml"
	"errors"
	"strings"
)

// XMLParser 是一个用于处理XML的工具类
type XMLParser struct{}

// Marshal 将对象序列化为XML字符串
func (p *XMLParser) Marshal(v any) (string, error) {
	data, err := xml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MarshalIndent 将对象序列化为带缩进的XML字符串
func (p *XMLParser) MarshalIndent(v any, prefix, indent string) (string, error) {
	data, err := xml.MarshalIndent(v, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Unmarshal 将XML字符串反序列化为对象
func (p *XMLParser) Unmarshal(data string, v any) error {
	return xml.Unmarshal([]byte(data), v)
}

// GetValueByPath 通过点分层级路径获取XML中的字段值
// 例如: "person.address.city" 将获取 <person><address><city>值</city></address></person> 中的值
func (p *XMLParser) GetValueByPath(xmlStr string, path string) (string, error) {
	// 解析XML到通用的map结构
	var result map[string]any
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))
	err := decoder.Decode(&result)
	if err != nil {
		return "", err
	}

	// 分解路径
	parts := strings.Split(path, ".")

	// 递归查找值
	value, err := navigateMap(result, parts)
	if err != nil {
		return "", err
	}

	// 转换为字符串
	switch v := value.(type) {
	case string:
		return v, nil
	default:
		// 如果值不是字符串，尝试转换为字符串
		valueBytes, err := xml.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(valueBytes), nil
	}
}

// navigateMap 递归地在map中导航，按照路径查找值
func navigateMap(data any, path []string) (any, error) {
	if len(path) == 0 {
		return data, nil
	}

	current := path[0]
	remaining := path[1:]

	// 检查是否为map
	m, ok := data.(map[string]any)
	if !ok {
		return nil, errors.New("无法解析路径: 非map类型")
	}

	// 查找当前路径部分
	value, exists := m[current]
	if !exists {
		return nil, errors.New("路径不存在: " + current)
	}

	// 如果这是最后一个路径部分，返回找到的值
	if len(remaining) == 0 {
		return value, nil
	}

	// 继续递归查找
	return navigateMap(value, remaining)
}

// GetXMLElement 获取XML中特定元素的完整XML片段
func (p *XMLParser) GetXMLElement(xmlStr string, elementName string) (string, error) {
	startTag := "<" + elementName
	endTag := "</" + elementName + ">"

	startIdx := strings.Index(xmlStr, startTag)
	if startIdx == -1 {
		return "", errors.New("未找到元素: " + elementName)
	}

	// 找到开始标签的结束位置
	closeBracketIdx := strings.Index(xmlStr[startIdx:], ">")
	if closeBracketIdx == -1 {
		return "", errors.New("XML格式错误: 缺少闭合标签")
	}

	startIdx += closeBracketIdx + 1
	endIdx := strings.Index(xmlStr, endTag)
	if endIdx == -1 {
		return "", errors.New("XML格式错误: 缺少结束标签 " + endTag)
	}

	endIdx += len(endTag)
	return xmlStr[startIdx-closeBracketIdx : endIdx], nil
}

// NewXMLParser 创建一个新的XML解析器实例
func NewXMLParser() *XMLParser {
	return &XMLParser{}
}
