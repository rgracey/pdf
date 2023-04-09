package main

import (
	"fmt"
	"os"

	"github.com/rgracey/pdf"
)

func main() {
	file, err := os.Open("sample.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	ast, _ := pdf.ParseStream(file)

	fmt.Printf("%+v", ast.Children())
}
