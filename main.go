// - Creating a data structure with load and sve methods
// - Using the net/http package to build web applications
// - Using the html/template package to process HTML templates
// - Using the regexp package to validate user input
// Reference: https://go.dev/doc/articles/wiki/
//
package main

import (
	"fmt"
	"log"
)

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.Save()
	p2, err := loadPage("TestPage")
	if err != nil {
		log.Println("loadPage error", err)
		return
	}
	fmt.Println("TestPage body:", string(p2.Body))
}
