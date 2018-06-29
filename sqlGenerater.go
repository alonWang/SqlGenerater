package main

import (
	"os"
	"io/ioutil"
	"strings"
	"github.com/tealeg/xlsx"
	"github.com/emirpasic/gods/maps/treemap"
	"fmt"
)

const INSERT_SQL_TEMPLATE = "insert into %s(%s) values \n %s;"

func main() {
	jsFile, err := os.Open("example.json")
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
	dataMap, data := ParseExcel("test.xlsm")
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

	fmt.Println(output)
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

func ParseExcel(filePath string) (map[string]int, [][]string) {
	file, err := xlsx.OpenFile("test.xlsm")
	if err != nil {
		panic(err)
	}
	row := file.Sheets[0].Rows[0]
	headerMap := parseHeader(row)
	body := parseBody(file.Sheets[0])
	return headerMap, body
}
func parseBody(sheet *xlsx.Sheet) [][]string {
	rows := sheet.Rows[1:]
	body := make([][]string, len(rows))
	for i := range rows {
		for j := range rows[i].Cells {
			body[i] = append(body[i], rows[i].Cells[j].Value)
		}
	}
	return body

}
func parseHeader(row *xlsx.Row) map[string]int {
	idx := 0
	cells := row.Cells
	headerMap := make(map[string]int)
	for cIdx := range cells {
		headerMap[cells[cIdx].Value] = idx
		idx++
	}
	return headerMap
}
