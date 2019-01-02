package excel

import (
	"testing"

	"github.com/tealeg/xlsx"
)

func TestNewExcel(t *testing.T) {
	a := NewExcel()

	b := new(excelWrite)
	b.File = xlsx.NewFile()

	if *a != *b {
		t.Error("NewExcel test error")
	}
}

func TestAddSheet(t *testing.T) {
	excel, err := NewExcel().AddSheet("pagination name")
	if err != nil {
		t.Error(err.Error())
	}
	if excel.Sheet.Name != "pagination name" {
		t.Error("Sheet not equal")
	}
}
