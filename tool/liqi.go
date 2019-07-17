package tool

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"bytes"
)

func fetchLatestLiqiJson() (jsonContent []byte, err error) {
	apiGetVersionURL := appendRandv(apiGetVersionZH)
	version, err := getVersion(apiGetVersionURL)
	if err != nil {
		return
	}

	apiGetResJsonURL := fmt.Sprintf(apiGetResVersionFormatZH, version.ResVersion)
	resource, err := getResource(apiGetResJsonURL)
	if err != nil {
		return
	}

	apiGetLiqiJsonURL := fmt.Sprintf(apiGetLiqiJsonFormatZH, resource.Res.LiqiJson.Prefix)
	return fetch(apiGetLiqiJsonURL)
}

func FetchLatestLiqiJson(filePath string) (err error) {
	jsonContent, err := fetchLatestLiqiJson()
	if err != nil {
		return
	}
	return ioutil.WriteFile(filePath, jsonContent, 0644)
}

//                "fields"   "game_url" "type"        string/int
//                "methods"  "login"    "requestType" string
//                "values"   "NULL"     int
//                "nested"   map[string]protoItem
type protoItem map[string]map[string]interface{}

type liqi struct {
	Nested struct {
		LQ struct {
			Nested map[string]protoItem `json:"nested"`
		} `json:"lq"`
	} `json:"nested"`
}

type converter struct {
	protoBB bytes.Buffer
	indent  int
}

func (c *converter) addLine(line string) {
	c.protoBB.WriteString(strings.Repeat("\t", c.indent) + line + "\n")
}

func (c *converter) newLine() {
	c.protoBB.WriteString("\n")
}

func (c *converter) startDefine(defineType string, name string) {
	c.addLine(fmt.Sprintf("%s %s {", defineType, name))
	c.indent++
}

func (c *converter) endDefine() {
	c.indent--
	c.addLine("}")
	if c.indent == 0 {
		c.newLine()
	}
}

func (c *converter) parseFields(fields map[string]interface{}) error {
	for fieldName, rawField := range fields {
		field, ok := rawField.(map[string]interface{})
		if !ok {
			return fmt.Errorf("parseFields 解析 %s 失败", fieldName)
		}
		if rule := field["rule"]; rule != nil {
			c.addLine(fmt.Sprintf("%v %v %s = %v;", rule, field["type"], fieldName, field["id"]))
		} else {
			c.addLine(fmt.Sprintf("%v %s = %v;", field["type"], fieldName, field["id"]))
		}
	}
	return nil
}

func (c *converter) parseMethods(methods map[string]interface{}) error {
	for methodName, rawMethod := range methods {
		method, ok := rawMethod.(map[string]interface{})
		if !ok {
			return fmt.Errorf("parseMethods 解析 %s 失败", methodName)
		}
		c.addLine(fmt.Sprintf("rpc %s (%v) returns (%v);", methodName, method["requestType"], method["responseType"]))
	}
	return nil
}

func (c *converter) parseValues(values map[string]interface{}) error {
	for name, value := range values {
		c.addLine(fmt.Sprintf("%s = %v;", name, value))
	}
	return nil
}

func (c *converter) parseFieldsProtoItem(name string, item protoItem) (err error) {
	fields, ok := item["fields"]
	if !ok {
		return
	}
	c.startDefine("message", name)
	if err = c.parseFields(fields); err != nil {
		return
	}
	if nestedItems, ok := item["nested"]; ok {
		var data []byte
		data, err = json.Marshal(nestedItems)
		if err != nil {
			return
		}
		nested := map[string]protoItem{}
		if err = json.Unmarshal(data, &nested); err != nil {
			return
		}
		for _name, _item := range nested {
			if err = c.parseFieldsProtoItem(_name, _item); err != nil {
				return
			}
		}
	}
	c.endDefine()
	return nil
}

func (c *converter) LiqiJsonToProto3(liqiJsonContent []byte) (protoContent []byte, err error) {
	c.addLine("syntax = \"proto3\";\n\npackage lq;\n")
	lq := liqi{}
	if err = json.Unmarshal(liqiJsonContent, &lq); err != nil {
		return
	}
	items := lq.Nested.LQ.Nested
	// 先处理 service 和 enum
	for name, item := range items {
		if methods, ok := item["methods"]; ok {
			c.startDefine("service", name)
			err = c.parseMethods(methods)
			c.endDefine()
		}
	}
	for name, item := range items {
		if values, ok := item["values"]; ok {
			c.startDefine("enum", name)
			err = c.parseValues(values)
			c.endDefine()
		}
	}
	// TODO: sort
	for name, item := range items {
		if err = c.parseFieldsProtoItem(name, item); err != nil {
			return
		}
	}
	return c.protoBB.Bytes(), nil
}

func liqiJsonToProto3(liqiJsonContent []byte) (protoContent []byte, err error) {
	return (&converter{}).LiqiJsonToProto3(liqiJsonContent)
}

func LiqiJsonToProto3(liqiJsonContent []byte, protoFilePath string) (err error) {
	content, err := liqiJsonToProto3(liqiJsonContent)
	if err != nil {
		return
	}
	return ioutil.WriteFile(protoFilePath, content, 0644)
}
