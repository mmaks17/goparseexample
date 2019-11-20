package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const MainURL = "HTTPS://какойтосайтнабитриксе.ру"

func main() {
	// Make HTTP request
	response, err := http.Get("MainURL")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find and print image URLs
	document.Find("a").Each(func(index int, element *goquery.Selection) {
		imgSrc, exists := element.Attr("href")
		if exists && strings.Index(imgSrc, "katalog-produktsii") > 0 && len(imgSrc) > 20 && strings.Count(imgSrc, "/") == 4 {
			//	fmt.Println(imgSrc)
			parceTover(imgSrc)
		}
	})
}

func parceTover(url string) {
	// Make HTTP request
	response, err := http.Get("MainURL" + url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	document.Find("p.h5 a").Each(func(index int, element *goquery.Selection) {
		imgSrc, exists := element.Attr("href")
		if exists && strings.Count(imgSrc, "/") == 5 {
			getPDP(imgSrc)
		}
	})
}
func getPDP(url string) {
	response, err := http.Get("MainURL" + url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	f, err := os.OpenFile("lines.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	document.Find(".text-break-word").Each(func(index int, element *goquery.Selection) {
		fmt.Printf(element.Text() + " ")
		_, err = fmt.Fprintln(f, element.Text()+",")
	})
	document.Find(".lead").Each(func(index int, element *goquery.Selection) {
		fmt.Printf(" " + element.Text() + " \n ")
		_, err = fmt.Fprintln(f, element.Text()+" ,")
	})

	fmt.Println("MainURL" + url)
	_, err = fmt.Fprintln(f, "MainURL"+url+",")
	document.Find("a.fancybox-thumb").Each(func(index int, element *goquery.Selection) {
		imgSrc, exists := element.Attr("href")
		if exists {
			fmt.Printf("%s \n", "MainURL"+imgSrc)
			_, err = fmt.Fprintln(f, "MainURL"+imgSrc+",")
		}
	})
	_, err = fmt.Fprintln(f, "\n\r")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}
