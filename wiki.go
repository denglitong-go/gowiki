package main

import (
	"log"
	"os"
)

const (
	fileExt = ".txt"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	filename := p.Title + fileExt
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + fileExt
	body, err := os.ReadFile(filename)
	if err != nil {
		log.Println("load page error", err)
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
