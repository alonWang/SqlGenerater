package main

func main() {
	//file, err := xlsx.OpenFile("test.xlsm")
	//if err != nil {
	//	panic(err)
	//}
	//parseBody(file.Sheets[0])
}

//func ParseExcel(filePath string) (map[string]int, [][]string) {
//	file, err := xlsx.OpenFile("test.xlsm")
//	if err != nil {
//		panic(err)
//	}
//	row := file.Sheets[0].Rows[0]
//	headerMap := parseHeader(row)
//	body := parseBody(file.Sheets[0])
//	return headerMap, body
//}
//func parseBody(sheet *xlsx.Sheet) [][]string{
//	rows := len(sheet.Rows) - 1
//	body:=make([][]string,rows)
//	for row := range sheet.Rows[1:] {
//		for col := range sheet.Rows[row].Cells {
//			body[row]=append(body[row],sheet.Rows[row].Cells[col].Value)
//			body[row][col] = sheet.Rows[row].Cells[col].Value
//		}
//	}
//	return body
//
//}
//
//func parseHeader(row *xlsx.Row) map[string]int {
//	idx := 0
//	cells := row.Cells
//	headerMap := make(map[string]int)
//	for cIdx := range cells {
//		log.Println(cells[cIdx].Value + ":" + strconv.Itoa(idx))
//		headerMap[cells[cIdx].Value] = idx
//		idx++
//	}
//	return headerMap
//}
