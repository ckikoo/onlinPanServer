package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/unidoc/unipdf/v3/model"
// )

// func main() {

// 	f, err := os.Open("1python.pdf")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	pdfReader, err := model.NewPdfReader(f)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(pdfReader.GetNumPages())

// 	fmt.Println(pdfReader.GetOutlines())

// 	lines, err := pdfReader.GetOutlines()
// 	if err != nil {
// 		panic(err)
// 	}

// 	for i := 0; i < len(lines.Items()); i++ {
// 		fmt.Printf("lines.Items()[i]: %+lv\n", lines.Items()[i])
// 	}
// 	// // Get outlines (bookmarks)
// 	// outlines, err := pdfReader.GetOutlines()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// // Print bookmarks
// 	// printOutlines(outlines, "")
// }

// // func printOutlines(outlines []*model.PdfOutlineItem, indent string) {
// // 	for _, outline := range outlines {
// // 		fmt.Printf("%sTitle: %s, Destination: %s\n", indent, outline.Title, outline.Dest)
// // 		if outline.Children != nil {
// // 			printOutlines(outline.Children, indent+"  ")
// // 		}
// // 	}
// // }
