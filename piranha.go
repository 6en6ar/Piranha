package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ledongthuc/pdf"
	"github.com/lu4p/cat"
)

var url = flag.String("u", "https://example.com", "url to be scraped")
var pdf_file = flag.String("p", "test.pdf", "pdf to be scraped")
var docx_file = flag.String("d", "test.docx", "docx to be scraped")
var banner = `

██████╗ ██╗██████╗  █████╗ ███╗   ██╗██╗  ██╗ █████╗ 
██╔══██╗██║██╔══██╗██╔══██╗████╗  ██║██║  ██║██╔══██╗
██████╔╝██║██████╔╝███████║██╔██╗ ██║███████║███████║
██╔═══╝ ██║██╔══██╗██╔══██║██║╚██╗██║██╔══██║██╔══██║
██║     ██║██║  ██║██║  ██║██║ ╚████║██║  ██║██║  ██║
╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝  ╚═╝
                                                        
			Coded by 6en6ar 3:-)

`

func main() {
	fmt.Println(banner)
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	pdf_set := false
	docx_set := false
	url_set := false
	isSet := func(f *flag.Flag) {
		if f.Name == "p" {
			pdf_set = true
		} else if f.Name == "d" {
			docx_set = true
		} else if f.Name == "u" {
			url_set = true
		}
	}
	flag.Visit(isSet)
	if pdf_set {
		pdfReader(*pdf_file)
	}
	if docx_set {
		docReader(*docx_file)
	}
	if url_set {
		GetWebContent(*url)
	}
	fmt.Println("[ + ] Output collected to passwords.txt")
	return

}
func GetWebContent(url string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("[ - ] Could not get website --> " + err.Error())

	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("[ - ] Response code not 200! --> " + err.Error())
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		fmt.Println("[ - ] Error reading body -->" + err.Error())
	}

	words := doc.Find("div").Map(func(i int, sel *goquery.Selection) string {
		return fmt.Sprintf("%d: %s", i+1, sel.Text())
	})
	text := strings.Join(words, " ")
	PrepareAndSave(text)
}
func RemoveDup(words []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range words {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func PrepareAndSave(w string) {
	file, err := os.OpenFile("passwords.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("[ - ] Error creating file --> " + err.Error())
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	words := strings.Fields(w)
	words = RemoveDup(words)
	for i := 0; i < len(words); i++ {
		//Remove words smaller than 5
		// Remove trailing , . or other signs
		chr := words[i][len(words[i])-1:]
		var res string
		if chr == "," || chr == "." || chr == "!" || chr == "?" {
			res = strings.TrimSuffix(words[i], chr)
		} else {
			res = words[i]
		}
		if len(res) < 5 {
			continue
		}
		_, _ = writer.WriteString(res + "\n")

	}
	writer.Flush()
}
func docReader(path string) {
	f, err := cat.File(path)
	if err != nil {
		fmt.Println("[ - ] Error opening file")
	}

	PrepareAndSave(f)
}
func pdfReader(path string) {
	f, r, err := pdf.Open(path)
	if err != nil {
		fmt.Println("[ - ] Error reading file")
	}
	defer f.Close()
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		fmt.Println("[ - ] Error getting text")
	}

	buf.ReadFrom(b)
	PrepareAndSave(buf.String())

}
