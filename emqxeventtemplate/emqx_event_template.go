package emqxeventtemplate

import (
	"strings"

	"github.com/yyle88/printgo"
	"github.com/yyle88/syntaxgo/syntaxgo_reflect"
)

// GetEmqxEventFieldSQL 根据消息体得到 event 的 select 语句
// https://docs.emqx.com/en/emqx/latest/data-integration/rule-sql-events-and-fields.html
func GetEmqxEventFieldSQL(object any, eventName string) string {
	var ptx = printgo.NewPTX()
	ptx.Println("SELECT")
	objectType := syntaxgo_reflect.GetTypeV3(object)
	for idx := 0; idx < objectType.NumField(); idx++ {
		field := objectType.Field(idx)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		var argName string
		if strings.Contains(jsonTag, ",") {
			parts := strings.Split(jsonTag, ",")
			argName = parts[0]
		} else {
			argName = jsonTag
		}
		ptx.Print("\t", argName)

		if idx+1 < objectType.NumField() {
			ptx.Println(",")
		} else {
			ptx.Println()
		}
	}
	ptx.Println(`FROM "` + eventName + `"`)
	res := ptx.String()
	return res
}

// GetEmqxEventTemplate 根据消息体得到 event 的 fields 模板
// https://docs.emqx.com/en/emqx/latest/data-integration/rule-sql-events-and-fields.html
func GetEmqxEventTemplate(object any) string {
	var ptx = printgo.NewPTX()
	ptx.Println("{")
	objectType := syntaxgo_reflect.GetTypeV3(object)
	for idx := 0; idx < objectType.NumField(); idx++ {
		field := objectType.Field(idx)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		var argName string
		if strings.Contains(jsonTag, ",") {
			parts := strings.Split(jsonTag, ",")
			argName = parts[0]
		} else {
			argName = jsonTag
		}

		if field.Type.Name() == "string" {
			ptx.Fprintf(`    "%s": "${%s}"`, argName, argName)
		} else {
			ptx.Fprintf(`    "%s": ${%s}`, argName, argName)
		}
		if idx+1 < objectType.NumField() {
			ptx.Println(",")
		} else {
			ptx.Println()
		}
	}
	ptx.Println("}")
	res := ptx.String()
	return res
}
