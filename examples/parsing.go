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

	pdfDoc, _ := pdf.ParseStream(file)

	fmt.Println("version:", pdfDoc.Version())

	for _, obj := range pdfDoc.Objects() {
		fmt.Printf(
			"ID: %d, Generation: %d, Header: %v, Data: %v\n",
			obj.Ref.Id,
			obj.Ref.Generation,
			obj.Header,
			obj.Data,
		)
	}
}
