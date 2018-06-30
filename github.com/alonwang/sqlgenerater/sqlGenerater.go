package main

import (
	"os"
	"io/ioutil"
	"strings"
	"github.com/emirpasic/gods/maps/treemap"
	"fmt"
	"flag"
)

const INSERT_SQL_TEMPLATE = "insert into %s(%s)values\n%s;"

func main() {
	mappingJson := flag.String("mapJson", "", "映射json文件地址")
	excelFile := flag.String("excel", "", "excel文件")
	GenerateSql(*mappingJson, *excelFile)
}
func GenerateSql(mappingJson, excelFile string) {
	jsFile, err := os.Open(mappingJson)
	if err != nil {
		panic(err)
	}

	jsData, err := ioutil.ReadAll(jsFile)

	mmm := treemap.NewWithStringComparator()
	mmm.FromJSON(jsData)
	values := mmm.Keys()
	tableName := values[0].(string)
	originData, _ := mmm.Get(tableName)
	unOrderedMap := originData.(map[string]interface{})
	orderedColKV := treemap.NewWithStringComparator()
	for k, v := range unOrderedMap {
		orderedColKV.Put(k, v)
	}
	dataMap, data := ParseExcel(excelFile)
	decor := fillDecora(orderedColKV)
	values = orderedColKV.Values()
	content := fillContent(values, dataMap, data)
	output := fmt.Sprintf(INSERT_SQL_TEMPLATE, tableName, decor, content)
	newFile, err := os.Create(tableName + "_generate.sql")
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
