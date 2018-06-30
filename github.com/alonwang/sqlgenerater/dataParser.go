package main

import (
	"github.com/tealeg/xlsx"
)

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
