package main

import (
	"os"
	"io/ioutil"
	"strings"
	"github.com/emirpasic/gods/maps/treemap"
	"fmt"
	"github.com/alonwang/sqlgenerater/dataParser"
	"flag"
)

const INSERT_SQL_TEMPLATE = "insert into %s(%s)values\n%s;"

func main() {
	mappingJson:= flag.String("f", "E:/sqlgenerate/example.json", "json file")
	flag.Parse()
	GenerateSql(*mappingJson)
}
func GenerateSql(mappingJson string) {
	jsFile, err := os.Open(mappingJson)
	if err != nil {
		panic(err)
	}
	jsData, err := ioutil.ReadAll(jsFile)
	//解析json,读取表名
	colsMap := treemap.NewWithStringComparator()
	colsMap.FromJSON(jsData)
	tableName,_:= colsMap.Get("table_name")

	//读取列对应关系并转化为有序map
	jsonData, _ := colsMap.Get("col_map")
	unOrderedMap := jsonData.(map[string]interface{})
	orderedColKV := treemap.NewWithStringComparator()
	for k, v := range unOrderedMap {
		orderedColKV.Put(k, v)
	}
	dataMap, data := dataparser.ParseExcel(mappingJson)
	decor := fillDecora(orderedColKV)
	values := orderedColKV.Values()
	content := fillContent(values, dataMap, data)
	output := fmt.Sprintf(INSERT_SQL_TEMPLATE, tableName, decor, content)
	newFile, err:=os.OpenFile(tableName.(string) + "_generate.sql", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	newFile.WriteString(output)
}
func fillContent(values []interface{}, dataMap map[string]int, data [][]string) string {
	contentTemplate := make([]string, 0)
	lines := make([]string, 0)
	for i := range data {
		for j := 0; j < len(values); j++ {
			contentTemplate = append(contentTemplate, data[i][dataMap[values[j].(string)]])
		}
		var line = "( "
		for _, value := range contentTemplate {
			line += value + ","
		}
		line = line[:len(line)-1]
		line += ")"
		contentTemplate = make([]string, 0)
		lines = append(lines, line)
	}
	return strings.Join(lines, ",\n")
}

func fillDecora(orderedMap *treemap.Map) string {
	var s []string
	orderedMap.Each(func(k interface{}, v interface{}) {
		s = append(s, k.(string))
	})
	return strings.Join(s, ",")
}
