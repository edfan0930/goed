package main

import (
	"github.com/edfan0930/goed/cmd"

	"github.com/edfan0930/goed/module/file"
)

type uu struct {
	Name string
	Male string
}

func main() {
	aa := make(chan string)
	file.ReadFile()
	<-aa
	return
	/* 	u := []uu{uu{"ebabab", "man"}, uu{"ac", "female"}}

	   	f := excel.NewExcel()
	   	if err := f.FieldName("cp", []string{"name", "male"}); err != nil {
	   		fmt.Println("17", err)
	   	}
	   	fmt.Println("uu", u)

	   	for i := range u {
	   		if err := f.WriteStruct(4, &u[i]); err != nil {
	   			fmt.Println("23", err)
	   			break
	   		}
	   	}

	   	if err := f.Save("/Users/red/go/src/abab/test.xlsx"); err != nil {
	   		fmt.Println("27", err)
	   	} */
	cmd.Execute()
}
