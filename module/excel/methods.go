package excel

import (
	"github.com/tealeg/xlsx"
)

type excelWrite struct {
	File  *xlsx.File
	Sheet *xlsx.Sheet
}

//NewExcel return excelWrite instance
func NewExcel() *excelWrite {
	return &excelWrite{
		File: xlsx.NewFile(),
	}
}

//AddSheet
//@params sheet string , sheet name
func (ex *excelWrite) AddSheet(sheet string) (*excelWrite, error) {
	var err error
	ex.Sheet, err = ex.File.AddSheet(sheet)
	return ex, err
}

//WriteSliceString
//@params cell []string , add to cell
func (ex *excelWrite) WriteSliceString(cell []string) {

	row := ex.Sheet.AddRow()
	for i := range cell {
		row.AddCell().Value = cell[i]
	}
	return
}

//WriteStruct
//@length length of row
//@v struct type
func (ex *excelWrite) WriteStruct(length int, v interface{}) {

	ex.Sheet.AddRow().WriteStruct(v, length)

	return
}

//Save
//Generate file
func (ex *excelWrite) Save(fileName string) error {
	return ex.File.Save(fileName)
}
