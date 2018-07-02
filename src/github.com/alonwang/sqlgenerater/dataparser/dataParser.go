package dataparser

import (
	"github.com/tealeg/xlsx"
	"os"
	"io/ioutil"
	"github.com/emirpasic/gods/maps/treemap"
)

func ParseExcel(jsonFile string) (map[string]int, [][]string) {
	jsFile, err := os.Open(jsonFile)
	if err != nil {
		panic(err)
	}
	jsData, err := ioutil.ReadAll(jsFile)
	//解析json,读取表名
	colsMap := treemap.NewWithStringComparator()
	colsMap.FromJSON(jsData)
	filePath, _ := colsMap.Get("excel_path")
	sheetName, _ := colsMap.Get("sheet_name")
	sheetIdLine, _ := colsMap.Get("sheet_id_line")
	originIgnoreLines, _ := colsMap.Get("ignore_line")
	var ignoreLines = make([]int, 0)
	for _, data := range originIgnoreLines.([]interface{}) {
		ignoreLines = append(ignoreLines, int(data.(float64))-1)
	}

	file, err := xlsx.OpenFile(filePath.(string))
	if err != nil {
		panic(err)
	}
	//匹配sheet
	var sheet xlsx.Sheet
	for _, st := range file.Sheets {
		if st.Name == sheetName.(string) {
			sheet = *st
			break
		}
	}
	//解析对应id列
	headerMap := parseHeader(sheet.Rows[int(sheetIdLine.(float64))-1])
	//解析数据
	body := parseBody(sheet, ignoreLines)
	return headerMap, body
}
func parseBody(sheet xlsx.Sheet, ignoreLines []int) [][]string {
	//去处无用行
	var rows = make([]xlsx.Row, 0)
	for row := range sheet.Rows {
		flag := true
		for i := range ignoreLines {
			if ignoreLines[i] == row {
				flag = false
				break
			}
		}
		if flag {
			rows = append(rows, *sheet.Rows[row])
		}
	}

	body := make([][]string, len(rows))
	line := make([]string, 0)
	idx := 0
	//去处空行
	for i := range rows {

		flag := true
		for j := range rows[i].Cells {
			line = append(line, rows[i].Cells[j].Value)
			//空或为0不取
			if len(line[j]) < 1 {
				flag = false
				line = line[:0]
				break
			}
		}
		if flag {
			for k := range line {
				body[idx] = append(body[idx], line[k])
			}
			idx++
			line = line[:0]
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
