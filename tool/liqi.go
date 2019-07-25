package tool

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"bytes"
	"sort"
	"reflect"
)

func fetchLatestLiqiJson() (jsonContent []byte, err error) {
	apiGetVersionURL := appendRandv(ApiGetVersionZH)
	version, err := GetMajsoulVersion(apiGetVersionURL)
	if err != nil {
		return
	}

	apiGetResJsonURL := fmt.Sprintf(apiGetResVersionFormatZH, version.ResVersion)
	resource, err := getResource(apiGetResJsonURL)
	if err != nil {
		return
	}

	apiGetLiqiJsonURL := fmt.Sprintf(apiGetLiqiJsonFormatZH, resource.Res.LiqiJson.Prefix)
	return Fetch(apiGetLiqiJsonURL)
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

type rpcMethod struct {
	name         string
	requestType  string
	responseType string
}

type rpcService struct {
	name    string
	methods []*rpcMethod
}

type converter struct {
	protoBB bytes.Buffer
	indent  int

	rpcServiceList      []*rpcService
	messageContainError map[string]struct{}
}

func newConverter() *converter {
	return &converter{
		messageContainError: map[string]struct{}{},
	}
}

func (c *converter) addLine(line string) {
	c.protoBB.WriteString(strings.Repeat("\t", c.indent) + line + "\n")
}

func (c *converter) newLine() {
	c.protoBB.WriteString("\n")
}

func (c *converter) startDefine(defineType string, name string) {
	if c.indent == 0 {
		c.newLine()
	}
	c.addLine(fmt.Sprintf("%s %s {", defineType, name))
	c.indent++
}

func (c *converter) endDefine() {
	c.indent--
	c.addLine("}")
}

func (*converter) sortedKeys(mp interface{}) (keys []string) {
	rawKeys := reflect.ValueOf(mp).MapKeys()
	for _, k := range rawKeys {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)
	return
}

func (c *converter) parseFields(rawFields map[string]interface{}) error {
	type field struct {
		id    int
		name  string
		type_ interface{}
		rule  interface{}
	}
	fields := []field{}
	for k, v := range rawFields {
		_v, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("parseFields 解析 %s 失败", k)
		}
		fields = append(fields, field{
			id:    int(_v["id"].(float64)),
			name:  k,
			type_: _v["type"],
			rule:  _v["rule"],
		})
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].id < fields[j].id
	})
	for _, field := range fields {
		if field.rule != nil {
			c.addLine(fmt.Sprintf("%v %v %s = %d;", field.rule, field.type_, field.name, field.id))
		} else {
			c.addLine(fmt.Sprintf("%v %s = %d;", field.type_, field.name, field.id))
		}
	}
	return nil
}

func (c *converter) parseMethods(methods map[string]interface{}) (rpcMethods []*rpcMethod, err error) {
	methodNames := c.sortedKeys(methods)
	for _, methodName := range methodNames {
		method, ok := methods[methodName].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("parseMethods 解析 %s 失败", methodName)
		}
		m := rpcMethod{
			name:         methodName,
			requestType:  method["requestType"].(string),
			responseType: method["responseType"].(string),
		}
		rpcMethods = append(rpcMethods, &m)
		c.addLine(fmt.Sprintf("rpc %s (%s) returns (%s);", m.name, m.requestType, m.responseType))
	}
	return
}

func (c *converter) parseEnums(enums map[string]interface{}) error {
	type kv struct {
		k string
		v int
	}
	pairs := []kv{}
	for k, v := range enums {
		pairs = append(pairs, kv{k, int(v.(float64))})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].v < pairs[j].v
	})
	for _, pair := range pairs {
		c.addLine(fmt.Sprintf("%s = %d;", pair.k, pair.v))
	}
	return nil
}

func (c *converter) parseFieldsProtoItem(name string, item protoItem) (err error) {
	fields, ok := item["fields"]
	if !ok {
		return
	}

	if _, ok := fields["error"]; ok {
		c.messageContainError[name] = struct{}{}
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
		_names := c.sortedKeys(nested)
		for _, _name := range _names {
			if err = c.parseFieldsProtoItem(_name, nested[_name]); err != nil {
				return
			}
		}
	}
	c.endDefine()
	return nil
}

func (c *converter) LiqiJsonToProto3(liqiJsonContent []byte) (protoContent []byte, err error) {
	c.addLine("syntax = \"proto3\";\n\npackage lq;")
	lq := liqi{}
	if err = json.Unmarshal(liqiJsonContent, &lq); err != nil {
		return
	}
	items := lq.Nested.LQ.Nested
	names := c.sortedKeys(items)

	// 先处理 service 和 enum
	for _, name := range names {
		if methods, ok := items[name]["methods"]; ok {
			c.startDefine("service", name)
			rpcMethods, er := c.parseMethods(methods)
			if er != nil {
				return nil, er
			}
			c.rpcServiceList = append(c.rpcServiceList, &rpcService{
				name:    name,
				methods: rpcMethods,
			})
			c.endDefine()
		}
	}
	for _, name := range names {
		if values, ok := items[name]["values"]; ok {
			c.startDefine("enum", name)
			if err = c.parseEnums(values); err != nil {
				return
			}
			c.endDefine()
		}
	}
	for _, name := range names {
		if err = c.parseFieldsProtoItem(name, items[name]); err != nil {
			return
		}
	}
	return c.protoBB.Bytes(), nil
}

func liqiJsonToProto3(liqiJsonContent []byte) (protoContent []byte, err error) {
	c := newConverter()
	return c.LiqiJsonToProto3(liqiJsonContent)
}

func LiqiJsonToProto3(liqiJsonContent []byte, protoFilePath string) (err error) {
	content, err := liqiJsonToProto3(liqiJsonContent)
	if err != nil {
		return
	}
	return ioutil.WriteFile(protoFilePath, content, 0644)
}
